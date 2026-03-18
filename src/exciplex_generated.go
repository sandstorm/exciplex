package exciplex

// #include <stdlib.h>
// #include "exciplex.h"
import "C"
import (
	"runtime/cgo"
	"unsafe"

	"github.com/dunglas/frankenphp"
)

func init() {
	frankenphp.RegisterExtension(unsafe.Pointer(&C.exciplex_module_entry))
}

func registerGoObject(obj interface{}) C.uintptr_t {
	handle := cgo.NewHandle(obj)
	return C.uintptr_t(handle)
}

func getGoObject(handle C.uintptr_t) interface{} {
	h := cgo.Handle(handle)
	return h.Value()
}

//export removeGoObject
func removeGoObject(handle C.uintptr_t) {
	h := cgo.Handle(handle)
	h.Delete()
}

//export go_exciplex_set_timeout
func go_exciplex_set_timeout(callback *C.zval, interval float64) C.uintptr_t {
	timer := SetTimeout(callback, interval)
	if timer == nil {
		return 0
	}
	return registerGoObject(timer)
}

//export go_exciplex_set_interval
func go_exciplex_set_interval(callback *C.zval, initialDelay float64, interval float64) C.uintptr_t {
	timer := SetInterval(callback, initialDelay, interval)
	if timer == nil {
		return 0
	}
	return registerGoObject(timer)
}

//export go_exciplex_timer_stop
func go_exciplex_timer_stop(handle C.uintptr_t) {
	obj := getGoObject(handle)
	if obj == nil {
		return
	}
	timer := obj.(*Timer)
	timer.Stop()
}

//export go_exciplex_start_profiler
func go_exciplex_start_profiler(initialDelay float64, interval float64) C.uintptr_t {
	profiler := StartProfiler(initialDelay, interval)
	if profiler == nil {
		return 0
	}
	return registerGoObject(profiler)
}

//export go_exciplex_profiler_stop
func go_exciplex_profiler_stop(handle C.uintptr_t) {
	obj := getGoObject(handle)
	if obj == nil {
		return
	}
	profiler := obj.(*Profiler)
	profiler.Stop()
}

//export go_exciplex_profiler_get_log
func go_exciplex_profiler_get_log(handle C.uintptr_t) unsafe.Pointer {
	obj := getGoObject(handle)
	if obj == nil {
		return nil
	}
	profiler := obj.(*Profiler)
	return frankenphp.PHPString(profiler.GetLog(), false)
}
