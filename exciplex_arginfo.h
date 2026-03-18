/* This is a generated file, edit exciplex.stub.php instead. */

/* Functions */

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_exciplex_set_timeout, 0, 2, ExciplexTimer, 0)
	ZEND_ARG_TYPE_INFO(0, callable, IS_CALLABLE, 0)
	ZEND_ARG_TYPE_INFO(0, interval, IS_DOUBLE, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_exciplex_set_interval, 0, 3, ExciplexTimer, 0)
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

/* ExciplexTimer class */

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_ExciplexTimer_stop, 0, 0, IS_VOID, 0)
ZEND_END_ARG_INFO()

ZEND_METHOD(ExciplexTimer, stop);

static const zend_function_entry class_ExciplexTimer_methods[] = {
	ZEND_ME(ExciplexTimer, stop, arginfo_class_ExciplexTimer_stop, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static zend_class_entry *register_class_ExciplexTimer(void) {
	zend_class_entry ce, *class_entry;

	INIT_CLASS_ENTRY(ce, "ExciplexTimer", class_ExciplexTimer_methods);
	class_entry = zend_register_internal_class_ex(&ce, NULL);
	class_entry->ce_flags |= ZEND_ACC_FINAL | ZEND_ACC_NO_DYNAMIC_PROPERTIES;

	return class_entry;
}
