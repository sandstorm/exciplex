#ifndef _EXCIPLEX_TIMEOUT_H
#define _EXCIPLEX_TIMEOUT_H

#include <php.h>
#include <stdatomic.h>
#include <stdint.h>

#include "exciplex_list.h"

// Status values for the timeout state machine
#define EXCIPLEX_TIMEOUT_PENDING   0
#define EXCIPLEX_TIMEOUT_TRIGGERED 1
#define EXCIPLEX_TIMEOUT_CANCELLED 2

typedef struct _exciplex_timeout_state {
    zend_atomic_bool *vm_interrupt;
    exciplex_mpsc_stack *triggered_stack;
    atomic_int status;
    bool repeating;
    uintptr_t callback_handle;
} exciplex_timeout_state;

// Called from PHP thread — one-shot or repeating setup
exciplex_timeout_state *exciplex_setup_timeout(bool repeated, uintptr_t callback_handle);

// Called from PHP thread to cancel a timeout (e.g. $timer->stop()).
// Only sets CANCELLED — RSHUTDOWN handles all resource cleanup.
void exciplex_cancel_timeout(exciplex_timeout_state *state);

// Called from Go goroutine after timer fires.
// Returns 0 on success, -1 if cancelled (request ended).
int exciplex_trigger_timeout(exciplex_timeout_state *state);

// Capture the current PHP stack trace as a malloc'd newline-separated string.
// Caller must free() the result.
char *exciplex_capture_stack_trace(void);

// Thin wrappers around PHP macros for use from Go
void exciplex_zval_copy(zval *dest, zval *src);
void exciplex_zval_dtor(zval *zv);

// Lifecycle hooks
void exciplex_timeout_minit(void);
void exciplex_timeout_rshutdown(void);

#endif
