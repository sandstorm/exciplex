package exciplex

// #include <stdlib.h>
// #include "exciplex_timeout.h"
import "C"
import (
	"runtime/cgo"
	"time"
)

//export go_exciplex_on_processed
func go_exciplex_on_processed(handle C.uintptr_t) {
	h := cgo.Handle(handle)
	fn := h.Value().(func())
	fn()
}

// export_php:class ExciplexTimer
type Timer struct {
	stopChan chan interface{}
	state    *C.exciplex_timeout_state
}

func startTimer(callback *C.zval, interval, delay float64, repeated bool) *Timer {
	state := C.exciplex_setup_timeout(callback, C.bool(repeated), 0)
	if state == nil {
		return nil
	}

	timeout := &Timer{
		stopChan: make(chan interface{}, 1),
		state:    state,
	}

	go func() {
		repeat := repeated
		for ok, waitTime := true, delay; ok; ok, waitTime = repeat, interval {
			select {
			case <-timeout.stopChan:
				repeat = false
			case <-time.After(time.Duration(waitTime * float64(time.Second))):
				repeat = C.exciplex_trigger_timeout(state) == 0
			}
		}
	}()

	return timeout
}

// export_php:function exciplex_set_timeout(callable $callable, float $interval): ExciplexTimer
func SetTimeout(callback *C.zval, interval float64) *Timer {
	return startTimer(callback, interval, interval, false)
}

// export_php:function exciplex_set_interval(callable $callable, float $initialDelay, float $interval): ExciplexTimer
func SetInterval(callback *C.zval, initialDelay, interval float64) *Timer {
	return startTimer(callback, interval, initialDelay, true)
}

// export_php:method ExciplexTimer::stop(): void
func (t *Timer) Stop() {
	// Cancel on PHP thread: sets CANCELLED, removes from pending list, cleans callback
	C.exciplex_cancel_timeout(t.state)
	t.state = nil
	// Signal goroutine to exit and free state
	t.stopChan <- struct{}{}
}
