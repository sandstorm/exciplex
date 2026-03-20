package src

// #include "exciplex_timeout.h"
import "C"
import (
	"unsafe"

	"github.com/dunglas/frankenphp"
)

type Callback interface {
	Call()
	Cleanup()
}

type CancellableCallback struct {
	inner     Callback
	cancelled bool // PHP-thread only — no sync needed
}

func (cc *CancellableCallback) Cancel() { cc.cancelled = true }
func (cc *CancellableCallback) Call() {
	if !cc.cancelled {
		cc.inner.Call()
	}
}
func (cc *CancellableCallback) Cleanup() { cc.inner.Cleanup() }

type GoCallback struct {
	fn func()
}

func (cb *GoCallback) Call()    { cb.fn() }
func (cb *GoCallback) Cleanup() {}

type PHPCallback struct {
	zval C.zval
}

func NewPHPCallback(zv *C.zval) *PHPCallback {
	cb := &PHPCallback{}
	C.exciplex_zval_copy(&cb.zval, zv)
	return cb
}

func (cb *PHPCallback) Call() {
	frankenphp.CallPHPCallable(unsafe.Pointer(&cb.zval), nil)
}

func (cb *PHPCallback) Cleanup() {
	C.exciplex_zval_dtor(&cb.zval)
}
