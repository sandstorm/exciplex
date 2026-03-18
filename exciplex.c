#include <php.h>
#include <Zend/zend_API.h>
#include <Zend/zend_hash.h>
#include <Zend/zend_types.h>
#include <stddef.h>

#include "exciplex.h"
#include "exciplex_timeout.h"
#include "exciplex_arginfo.h"
#include "_cgo_export.h"

// --- ExciplexTimer PHP object infrastructure ---

zend_class_entry *exciplex_timer_ce = NULL;
static zend_object_handlers exciplex_object_handlers;

typedef struct {
    uintptr_t go_handle;
    zend_object std; /* must be last */
} exciplex_object;

static inline exciplex_object *exciplex_object_from_obj(zend_object *obj) {
    return (exciplex_object*)((char*)(obj) - offsetof(exciplex_object, std));
}

static zend_object *exciplex_create_object(zend_class_entry *ce) {
    exciplex_object *intern = ecalloc(1, sizeof(exciplex_object) + zend_object_properties_size(ce));
    zend_object_std_init(&intern->std, ce);
    object_properties_init(&intern->std, ce);
    intern->std.handlers = &exciplex_object_handlers;
    intern->go_handle = 0;
    return &intern->std;
}

static void exciplex_free_object(zend_object *object) {
    // Do NOT call removeGoObject here — this runs during PHP request
    // shutdown, after FrankenPHP may have torn down the Go context
    // for this thread. The Go Timer's lifecycle is managed by its
    // goroutine; the cgo.Handle is cleaned up when the goroutine exits.
    zend_object_std_dtor(object);
}

// --- Module lifecycle ---

PHP_MINIT_FUNCTION(exciplex) {
    exciplex_timeout_minit();

    // Init object handlers
    memcpy(&exciplex_object_handlers, &std_object_handlers, sizeof(zend_object_handlers));
    exciplex_object_handlers.free_obj = exciplex_free_object;
    exciplex_object_handlers.clone_obj = NULL;
    exciplex_object_handlers.offset = offsetof(exciplex_object, std);

    // Register ExciplexTimer class
    exciplex_timer_ce = register_class_ExciplexTimer();
    exciplex_timer_ce->create_object = exciplex_create_object;

    return SUCCESS;
}

PHP_RSHUTDOWN_FUNCTION(exciplex) {
    exciplex_timeout_rshutdown();
    return SUCCESS;
}

zend_module_entry exciplex_module_entry = {STANDARD_MODULE_HEADER,
                                         "exciplex",
                                         ext_functions,             /* Functions */
                                         PHP_MINIT(exciplex),  /* MINIT */
                                         NULL,                      /* MSHUTDOWN */
                                         NULL,                      /* RINIT */
                                         PHP_RSHUTDOWN(exciplex),   /* RSHUTDOWN */
                                         NULL,                      /* MINFO */
                                         "1.0.0",                   /* Version */
                                         STANDARD_MODULE_PROPERTIES};

// --- PHP functions ---

PHP_FUNCTION(exciplex_set_timeout)
{
    zval *callable_callback;
    double interval = 0.0;
    ZEND_PARSE_PARAMETERS_START(2, 2)
        Z_PARAM_ZVAL(callable_callback)
        Z_PARAM_DOUBLE(interval)
    ZEND_PARSE_PARAMETERS_END();

    uintptr_t handle = go_exciplex_set_timeout(callable_callback, (double) interval);
    if (handle == 0) {
        RETURN_NULL();
    }

    object_init_ex(return_value, exciplex_timer_ce);
    exciplex_object *intern = exciplex_object_from_obj(Z_OBJ_P(return_value));
    intern->go_handle = handle;
}

PHP_FUNCTION(exciplex_set_interval)
{
    zval *callable_callback;
    double initialDelay = 0.0;
    double interval = 0.0;
    ZEND_PARSE_PARAMETERS_START(3, 3)
        Z_PARAM_ZVAL(callable_callback)
        Z_PARAM_DOUBLE(initialDelay)
        Z_PARAM_DOUBLE(interval)
    ZEND_PARSE_PARAMETERS_END();

    uintptr_t handle = go_exciplex_set_interval(callable_callback, (double) initialDelay, (double) interval);
    if (handle == 0) {
        RETURN_NULL();
    }

    object_init_ex(return_value, exciplex_timer_ce);
    exciplex_object *intern = exciplex_object_from_obj(Z_OBJ_P(return_value));
    intern->go_handle = handle;
}

// --- ExciplexTimer methods ---

PHP_METHOD(ExciplexTimer, stop)
{
    ZEND_PARSE_PARAMETERS_NONE();

    exciplex_object *intern = exciplex_object_from_obj(Z_OBJ_P(ZEND_THIS));
    if (intern->go_handle == 0) {
        zend_throw_error(NULL, "Timer not initialized");
        RETURN_THROWS();
    }
    go_exciplex_timer_stop(intern->go_handle);
}
