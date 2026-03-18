package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFAccessibilityHandlerT struct {
	_                             structs.HostLayout
	Base                          CEFBaseRefCountedT
	OnAccessibilityTreeChange     uintptr
	OnAccessibilityLocationChange uintptr
}

func (v *CEFAccessibilityHandlerT) OverrideOnAccessibilityTreeChange(fn uintptr) {
	v.OnAccessibilityTreeChange = fn
}

func (v *CEFAccessibilityHandlerT) CallOnAccessibilityTreeChange(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnAccessibilityTreeChange, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFAccessibilityHandlerT) OverrideOnAccessibilityLocationChange(fn uintptr) {
	v.OnAccessibilityLocationChange = fn
}

func (v *CEFAccessibilityHandlerT) CallOnAccessibilityLocationChange(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnAccessibilityLocationChange, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterAccessibilityHandler(handle uintptr) {
}
