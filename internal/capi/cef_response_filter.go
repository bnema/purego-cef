package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFResponseFilterT struct {
	_          structs.HostLayout
	Base       CEFBaseRefCountedT
	InitFilter uintptr
	Filter     uintptr
}

func (v *CEFResponseFilterT) OverrideInitFilter(fn uintptr) { v.InitFilter = fn }

func (v *CEFResponseFilterT) CallInitFilter(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.InitFilter, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseFilterT) OverrideFilter(fn uintptr) { v.Filter = fn }

func (v *CEFResponseFilterT) CallFilter(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Filter, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterResponseFilter(handle uintptr) {
}
