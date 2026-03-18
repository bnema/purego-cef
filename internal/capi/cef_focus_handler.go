package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFFocusHandlerT struct {
	_           structs.HostLayout
	Base        CEFBaseRefCountedT
	OnTakeFocus uintptr
	OnSetFocus  uintptr
	OnGotFocus  uintptr
}

func (v *CEFFocusHandlerT) OverrideOnTakeFocus(fn uintptr) { v.OnTakeFocus = fn }

func (v *CEFFocusHandlerT) CallOnTakeFocus(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnTakeFocus, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFFocusHandlerT) OverrideOnSetFocus(fn uintptr) { v.OnSetFocus = fn }

func (v *CEFFocusHandlerT) CallOnSetFocus(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnSetFocus, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFFocusHandlerT) OverrideOnGotFocus(fn uintptr) { v.OnGotFocus = fn }

func (v *CEFFocusHandlerT) CallOnGotFocus(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnGotFocus, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterFocusHandler(handle uintptr) {
}
