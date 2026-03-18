package exciplex

// #include <stdlib.h>
// #include "exciplex_timeout.h"
import "C"
import (
	"hash/fnv"
	"runtime/cgo"
	"slices"
	"strconv"
	"strings"
	"time"
	"unsafe"
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

func startTimer(callback *C.zval, interval, delay float64, repeated bool, onProcessed func()) *Timer {
	var onProcessedHandle C.uintptr_t
	if onProcessed != nil {
		onProcessedHandle = C.uintptr_t(cgo.NewHandle(onProcessed))
	}
	state := C.exciplex_setup_timeout(callback, C.bool(repeated), onProcessedHandle)
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
	return startTimer(callback, interval, interval, false, nil)
}

// export_php:function exciplex_set_interval(callable $callable, float $initialDelay, float $interval): ExciplexTimer
func SetInterval(callback *C.zval, initialDelay, interval float64) *Timer {
	return startTimer(callback, interval, initialDelay, true, nil)
}

// export_php:method ExciplexTimer::stop(): void
func (t *Timer) Stop() {
	// Cancel on PHP thread: sets CANCELLED, removes from pending list, cleans callback
	C.exciplex_cancel_timeout(t.state)
	t.state = nil
	// Signal goroutine to exit and free state
	t.stopChan <- struct{}{}
}

type LogEntry struct {
	stacktrace  []string
	occurrences int
}

// export_php:class ExciplexProfiler
type Profiler struct {
	timer *Timer
	logs  map[uint64]*LogEntry
}

// export_php:function start_profiler(float $initalDelay, float $interval): ExciplexProfiler
func StartProfiler(initialDelay float64, interval float64) *Profiler {
	p := &Profiler{
		logs: make(map[uint64]*LogEntry),
	}
	p.timer = startTimer(nil, interval, initialDelay, true, func() {
		cstr := C.exciplex_capture_stack_trace()
		defer C.free(unsafe.Pointer(cstr))
		goStr := C.GoString(cstr)
		if goStr == "" {
			return
		}
		h := fnv.New64a()
		h.Write([]byte(goStr))
		key := h.Sum64()

		if entry, ok := p.logs[key]; ok {
			entry.occurrences++
		} else {
			frames := strings.Split(goStr, "\n")
			slices.Reverse(frames)
			p.logs[key] = &LogEntry{stacktrace: frames, occurrences: 1}
		}
	})
	return p
}

// export_php:method ExciplexProfiler::stop(): void
func (p *Profiler) Stop() {
	p.timer.Stop()
}

// export_php:method ExciplexProfiler::getLog(): string
func (p *Profiler) GetLog() string {
	var sb strings.Builder
	for _, entry := range p.logs {
		sb.WriteString(strings.Join(entry.stacktrace, ";"))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(entry.occurrences))
		sb.WriteByte('\n')
	}
	return sb.String()
}
