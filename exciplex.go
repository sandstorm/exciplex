package exciplex

// #include "exciplex_timeout.h"
import "C"
import "time"

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
