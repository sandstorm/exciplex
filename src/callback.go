package exciplex

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
