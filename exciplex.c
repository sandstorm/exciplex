#include <php.h>
#include <Zend/zend_API.h>
#include <Zend/zend_hash.h>
#include <Zend/zend_types.h>
#include <stddef.h>

#include "exciplex.h"
#include "exciplex_timeout.h"
#include "exciplex_arginfo.h"
#include "_cgo_export.h"


PHP_MINIT_FUNCTION(exciplex) {
    exciplex_timeout_minit();
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
PHP_FUNCTION(exciplex_set_timeout)
{
    zval *callable_callback;
    double interval = 0.0;
    ZEND_PARSE_PARAMETERS_START(2, 2)
        Z_PARAM_ZVAL(callable_callback)
        Z_PARAM_DOUBLE(interval)
    ZEND_PARSE_PARAMETERS_END();
    go_exciplex_set_timeout(callable_callback, (double) interval);
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
    go_exciplex_set_interval(callable_callback, (double) initialDelay, (double) interval);
}

