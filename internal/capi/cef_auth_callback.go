package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFAuthCallbackT struct {
	_      structs.HostLayout
	Base   CEFBaseRefCountedT
	Cont   uintptr
	Cancel uintptr
}

func (v *CEFAuthCallbackT) OverrideCont(fn uintptr) { v.Cont = fn }

func (v *CEFAuthCallbackT) CallCont(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cont, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFAuthCallbackT) OverrideCancel(fn uintptr) { v.Cancel = fn }

func (v *CEFAuthCallbackT) CallCancel(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cancel, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterAuthCallback(handle uintptr) {
}
