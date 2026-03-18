package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFLayoutT struct {
	_            structs.HostLayout
	Base         CEFBaseRefCountedT
	AsBoxLayout  uintptr
	AsFillLayout uintptr
	IsValid      uintptr
}

func (v *CEFLayoutT) OverrideAsBoxLayout(fn uintptr) { v.AsBoxLayout = fn }

func (v *CEFLayoutT) CallAsBoxLayout(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AsBoxLayout, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFLayoutT) OverrideAsFillLayout(fn uintptr) { v.AsFillLayout = fn }

func (v *CEFLayoutT) CallAsFillLayout(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AsFillLayout, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFLayoutT) OverrideIsValid(fn uintptr) { v.IsValid = fn }

func (v *CEFLayoutT) CallIsValid(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsValid, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterLayout(handle uintptr) {
}
