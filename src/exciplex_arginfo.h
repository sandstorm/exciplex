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

ZEND_BEGIN_ARG_WITH_RETURN_OBJ_INFO_EX(arginfo_start_profiler, 0, 2, ExciplexProfiler, 0)
	ZEND_ARG_TYPE_INFO(0, initialDelay, IS_DOUBLE, 0)
	ZEND_ARG_TYPE_INFO(0, interval, IS_DOUBLE, 0)
ZEND_END_ARG_INFO()

ZEND_FUNCTION(exciplex_set_timeout);
ZEND_FUNCTION(exciplex_set_interval);
ZEND_FUNCTION(start_profiler);

static const zend_function_entry ext_functions[] = {
	ZEND_FE(exciplex_set_timeout, arginfo_exciplex_set_timeout)
	ZEND_FE(exciplex_set_interval, arginfo_exciplex_set_interval)
	ZEND_FE(start_profiler, arginfo_start_profiler)
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

/* ExciplexProfiler class */

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_ExciplexProfiler_stop, 0, 0, IS_VOID, 0)
ZEND_END_ARG_INFO()

ZEND_BEGIN_ARG_WITH_RETURN_TYPE_INFO_EX(arginfo_class_ExciplexProfiler_getLog, 0, 0, IS_STRING, 0)
ZEND_END_ARG_INFO()

ZEND_METHOD(ExciplexProfiler, stop);
ZEND_METHOD(ExciplexProfiler, getLog);

static const zend_function_entry class_ExciplexProfiler_methods[] = {
	ZEND_ME(ExciplexProfiler, stop, arginfo_class_ExciplexProfiler_stop, ZEND_ACC_PUBLIC)
	ZEND_ME(ExciplexProfiler, getLog, arginfo_class_ExciplexProfiler_getLog, ZEND_ACC_PUBLIC)
	ZEND_FE_END
};

static zend_class_entry *register_class_ExciplexProfiler(void) {
	zend_class_entry ce, *class_entry;

	INIT_CLASS_ENTRY(ce, "ExciplexProfiler", class_ExciplexProfiler_methods);
	class_entry = zend_register_internal_class_ex(&ce, NULL);
	class_entry->ce_flags |= ZEND_ACC_FINAL | ZEND_ACC_NO_DYNAMIC_PROPERTIES;

	return class_entry;
}
