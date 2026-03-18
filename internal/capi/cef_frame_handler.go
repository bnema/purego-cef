package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFFrameHandlerT struct {
	_                  structs.HostLayout
	Base               CEFBaseRefCountedT
	OnFrameCreated     uintptr
	OnFrameDestroyed   uintptr
	OnFrameAttached    uintptr
	OnFrameDetached    uintptr
	OnMainFrameChanged uintptr
}

func (v *CEFFrameHandlerT) OverrideOnFrameCreated(fn uintptr) { v.OnFrameCreated = fn }

func (v *CEFFrameHandlerT) CallOnFrameCreated(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnFrameCreated, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFFrameHandlerT) OverrideOnFrameDestroyed(fn uintptr) { v.OnFrameDestroyed = fn }

func (v *CEFFrameHandlerT) CallOnFrameDestroyed(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnFrameDestroyed, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFFrameHandlerT) OverrideOnFrameAttached(fn uintptr) { v.OnFrameAttached = fn }

func (v *CEFFrameHandlerT) CallOnFrameAttached(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnFrameAttached, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFFrameHandlerT) OverrideOnFrameDetached(fn uintptr) { v.OnFrameDetached = fn }

func (v *CEFFrameHandlerT) CallOnFrameDetached(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnFrameDetached, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFFrameHandlerT) OverrideOnMainFrameChanged(fn uintptr) { v.OnMainFrameChanged = fn }

func (v *CEFFrameHandlerT) CallOnMainFrameChanged(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnMainFrameChanged, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterFrameHandler(handle uintptr) {
}
