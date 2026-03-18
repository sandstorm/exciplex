package exciplex

// #include <stdlib.h>
// #include "exciplex_timeout.h"
import "C"
import (
	"runtime/cgo"
	"time"
	"unsafe"
)

//export go_exciplex_on_processed
func go_exciplex_on_processed(handle C.uintptr_t) {
	h := cgo.Handle(handle)
	fn := h.Value().(func())
	fn()
}

// export_php:function exciplex_set_timeout(callable $callable, float $interval): void
func SetTimeout(callback *C.zval, interval float64) {
	state := C.exciplex_setup_timeout(callback)
	if state == nil {
		return
	}
	go func() {
		<-time.After(time.Duration(interval * float64(time.Second)))
		C.exciplex_trigger_timeout(state)
	}()
}

// export_php:function exciplex_set_interval(callable $callable, float $initialDelay, float $interval): void
func SetInterval(callback *C.zval, initialDelay, interval float64) {
	//ch := make(chan struct{}, 1)
	//h := cgo.NewHandle(func() { ch <- struct{}{} })

	state := C.exciplex_setup_repeating_timeout(callback, 0) //C.uintptr_t(h))
	if state == nil {
		//h.Delete()
		return
	}
	go func() {
		//defer h.Delete()
		time.Sleep(time.Duration(initialDelay * float64(time.Second)))
		for {
			if C.exciplex_trigger_timeout(state) < 0 {
				C.free(unsafe.Pointer(state))
				return
			}
			//<-ch // block until handler processes this tick
			time.Sleep(time.Duration(interval * float64(time.Second)))
		}
	}()
}
