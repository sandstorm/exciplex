#include <php.h>
#include <Zend/zend_execute.h>
#include <stdatomic.h>
#include <stdlib.h>

#include "exciplex_timeout.h"

// Exported Go function — calls the Go closure stored in the cgo.Handle.
extern void go_exciplex_on_processed(uintptr_t handle);

// Per-PHP-thread list of all pending timeout states (for RSHUTDOWN cleanup).
// Only accessed on the PHP thread — no synchronization needed.
static __thread exciplex_list pending = {NULL};

// Per-PHP-thread MPSC stack. Goroutines push triggered states here
// (via the triggered_stack pointer captured at setup time), then set vm_interrupt.
// The interrupt handler pops all and processes them.
static __thread exciplex_mpsc_stack triggered = {(exciplex_list_node *)NULL};

// Previous interrupt handler to chain to
static void (*prev_zend_interrupt_function)(zend_execute_data *) = NULL;

// Interrupt handler — runs on the PHP thread when EG(vm_interrupt) is set.
// Pops all triggered states and processes each callback.
static void exciplex_interrupt_handler(zend_execute_data *execute_data) {
    exciplex_list_node *node = exciplex_mpsc_stack_pop_all(&triggered);

    while (node != NULL) {
        exciplex_list_node *next = node->next;
        exciplex_timeout_state *state = (exciplex_timeout_state *)node->data;
        free(node);

        if (state->repeating) {
            // Reset status BEFORE calling — goroutine can retrigger after interval.
            // If callback calls die(), RSHUTDOWN will cancel the still-PENDING state.
            atomic_store(&state->status, EXCIPLEX_TIMEOUT_PENDING);

            zval retval;
            ZVAL_UNDEF(&retval);
            call_user_function(NULL, NULL, &state->callback, &retval, 0, NULL);
            zval_ptr_dtor(&retval);

            // Notify the Go goroutine that this tick was processed
            if (state->on_processed_handle != 0) {
                go_exciplex_on_processed(state->on_processed_handle);
            }
        } else {
            // One-shot: set CANCELLED so goroutine exits on next trigger attempt.
            // Don't free state — goroutine will free via defer.
            // Don't remove from pending — RSHUTDOWN will clean up.
            atomic_store(&state->status, EXCIPLEX_TIMEOUT_CANCELLED);

            zval callback;
            ZVAL_COPY_VALUE(&callback, &state->callback);
            ZVAL_UNDEF(&state->callback); // prevent RSHUTDOWN double-dtor

            zval retval;
            ZVAL_UNDEF(&retval);
            call_user_function(NULL, NULL, &callback, &retval, 0, NULL);
            zval_ptr_dtor(&retval);
            zval_ptr_dtor(&callback);
        }

        node = next;
    }

    if (prev_zend_interrupt_function) {
        prev_zend_interrupt_function(execute_data);
    }
}

exciplex_timeout_state *exciplex_setup_timeout(zval *callback, bool repeating, uintptr_t on_processed_handle) {
    exciplex_timeout_state *state = malloc(sizeof(exciplex_timeout_state));
    if (state == NULL) {
        return NULL;
    }

    ZVAL_COPY(&state->callback, callback);
    state->vm_interrupt = &EG(vm_interrupt);
    state->triggered_stack = &triggered;
    atomic_store(&state->status, EXCIPLEX_TIMEOUT_PENDING);
    state->repeating = repeating;
    state->on_processed_handle = on_processed_handle;

    exciplex_list_prepend(&pending, state);

    return state;
}

// Called from PHP thread to cancel a timeout (e.g. $timer->stop()).
// State stays in pending list — RSHUTDOWN will free it.
void exciplex_cancel_timeout(exciplex_timeout_state *state) {
    atomic_store(&state->status, EXCIPLEX_TIMEOUT_CANCELLED);
    zval_ptr_dtor(&state->callback);
    ZVAL_UNDEF(&state->callback);
}

// Called from a Go goroutine (arbitrary OS thread) after the timer fires.
// Returns 0 on success, -1 if cancelled.
int exciplex_trigger_timeout(exciplex_timeout_state *state) {
    // Copy vm_interrupt locally — after status transitions,
    // RSHUTDOWN on the PHP thread may free the state.
    zend_atomic_bool *vm_interrupt = state->vm_interrupt;

    // Whitelist: allow push from PENDING or TRIGGERED states.
    // Try PENDING → TRIGGERED first, then TRIGGERED → TRIGGERED.
    // If neither succeeds, status must be CANCELLED.
    int expected = EXCIPLEX_TIMEOUT_PENDING;
    if (!atomic_compare_exchange_strong(&state->status, &expected, EXCIPLEX_TIMEOUT_TRIGGERED)) {
        expected = EXCIPLEX_TIMEOUT_TRIGGERED;
        if (!atomic_compare_exchange_strong(&state->status, &expected, EXCIPLEX_TIMEOUT_TRIGGERED)) {
            // Status is CANCELLED — RSHUTDOWN will free state.
            return -1;
        }
    }

    // Push onto the MPSC stack and signal the PHP VM.
    // Each push = one callback invocation by the handler.
    exciplex_mpsc_stack_push(state->triggered_stack, state);
    zend_atomic_bool_store(vm_interrupt, true);
    return 0;
}

void exciplex_timeout_minit(void) {
    prev_zend_interrupt_function = zend_interrupt_function;
    zend_interrupt_function = exciplex_interrupt_handler;
}

// Called during request shutdown (on the PHP thread).
void exciplex_timeout_rshutdown(void) {
    // First: cancel all states so goroutines stop pushing new triggers.
    exciplex_list_node *node = exciplex_list_drain(&pending);
    exciplex_list_node *first = node;
    while (node != NULL) {
        exciplex_timeout_state *state = (exciplex_timeout_state *)node->data;

        // Set CANCELLED — goroutines will see this and stop.
        atomic_store(&state->status, EXCIPLEX_TIMEOUT_CANCELLED);

        // Clean callback if not already cleaned (cancel/one-shot handler set UNDEF).
        zval_ptr_dtor(&state->callback);

        node = node->next;
    }

    // Second: drain any triggered nodes pushed before/during cancellation.
    exciplex_list_node *tnode = exciplex_mpsc_stack_pop_all(&triggered);
    while (tnode != NULL) {
        exciplex_list_node *tnext = tnode->next;
        free(tnode);
        tnode = tnext;
    }

    // Third: free all states and list nodes. At this point all goroutines
    // have seen CANCELLED (or will see it and just return -1 without
    // accessing state further since the MPSC stack was drained).
    node = first;
    while (node != NULL) {
        exciplex_list_node *next = node->next;
        free(node->data); // free state
        free(node);       // free list node
        node = next;
    }
}
