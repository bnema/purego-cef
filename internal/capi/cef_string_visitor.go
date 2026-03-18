package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFStringVisitorT struct {
	_     structs.HostLayout
	Base  CEFBaseRefCountedT
	Visit uintptr
}

func (v *CEFStringVisitorT) OverrideVisit(fn uintptr) { v.Visit = fn }

func (v *CEFStringVisitorT) CallVisit(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Visit, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterStringVisitor(handle uintptr) {
}
