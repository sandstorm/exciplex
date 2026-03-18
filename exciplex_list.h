#ifndef _EXCIPLEX_LIST_H
#define _EXCIPLEX_LIST_H

#include <stdatomic.h>
#include <stdlib.h>

// Common node for both list types
typedef struct _exciplex_list_node {
    void *data;
    struct _exciplex_list_node *next;
} exciplex_list_node;

// ---------- Simple list (single-thread, no locking) ----------

typedef struct {
    exciplex_list_node *head;
} exciplex_list;

static inline void exciplex_list_prepend(exciplex_list *list, void *data) {
    exciplex_list_node *node = malloc(sizeof(exciplex_list_node));
    node->data = data;
    node->next = list->head;
    list->head = node;
}

// Remove first node matching data. Returns data if found, NULL otherwise.
static inline void *exciplex_list_remove(exciplex_list *list, void *data) {
    exciplex_list_node **pp = &list->head;
    while (*pp != NULL) {
        if ((*pp)->data == data) {
            exciplex_list_node *node = *pp;
            *pp = node->next;
            free(node);
            return data;
        }
        pp = &(*pp)->next;
    }
    return NULL;
}

// Detach and return the entire chain. Caller walks via ->next and frees nodes.
static inline exciplex_list_node *exciplex_list_drain(exciplex_list *list) {
    exciplex_list_node *head = list->head;
    list->head = NULL;
    return head;
}

// ---------- MPSC stack (multi-producer, single-consumer, lock-free) ----------

typedef struct {
    _Atomic(exciplex_list_node *) head;
} exciplex_mpsc_stack;

// Push from any thread (lock-free CAS).
static inline void exciplex_mpsc_stack_push(exciplex_mpsc_stack *stack, void *data) {
    exciplex_list_node *node = malloc(sizeof(exciplex_list_node));
    node->data = data;
    node->next = atomic_load(&stack->head);
    while (!atomic_compare_exchange_weak(&stack->head, &node->next, node)) {
        // node->next is updated by CAS on failure
    }
}

// Pop all nodes atomically (single consumer). Returns chain, caller frees.
static inline exciplex_list_node *exciplex_mpsc_stack_pop_all(exciplex_mpsc_stack *stack) {
    return atomic_exchange(&stack->head, (exciplex_list_node *)NULL);
}

#endif
