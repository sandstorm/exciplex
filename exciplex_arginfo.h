/* This is a generated file, edit exciplex.stub.php instead. */

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_exciplex_set_timeout, 0, 2, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, callable, IS_CALLABLE, 0)
	ZEND_ARG_TYPE_INFO(0, interval, IS_DOUBLE, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_exciplex_set_interval, 0, 3, IS_VOID, 0)
	ZEND_ARG_TYPE_INFO(0, callable, IS_CALLABLE, 0)
	ZEND_ARG_TYPE_INFO(0, initialDelay, IS_DOUBLE, 0)
	ZEND_ARG_TYPE_INFO(0, interval, IS_DOUBLE, 0)
ZEND_END_ARG_INFO()

ZEND_FUNCTION(exciplex_set_timeout);
ZEND_FUNCTION(exciplex_set_interval);

static const zend_function_entry ext_functions[] = {
	ZEND_FE(exciplex_set_timeout, arginfo_exciplex_set_timeout)
	ZEND_FE(exciplex_set_interval, arginfo_exciplex_set_interval)
	ZEND_FE_END
};
