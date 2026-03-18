package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFRequestContextHandlerT struct {
	_                           structs.HostLayout
	Base                        CEFBaseRefCountedT
	OnRequestContextInitialized uintptr
	GetResourceRequestHandler   uintptr
}

func (v *CEFRequestContextHandlerT) OverrideOnRequestContextInitialized(fn uintptr) {
	v.OnRequestContextInitialized = fn
}

func (v *CEFRequestContextHandlerT) CallOnRequestContextInitialized(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnRequestContextInitialized, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestContextHandlerT) OverrideGetResourceRequestHandler(fn uintptr) {
	v.GetResourceRequestHandler = fn
}

func (v *CEFRequestContextHandlerT) CallGetResourceRequestHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetResourceRequestHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterRequestContextHandler(handle uintptr) {
}
