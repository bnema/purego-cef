package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFFindHandlerT struct {
	_            structs.HostLayout
	Base         CEFBaseRefCountedT
	OnFindResult uintptr
}

func (v *CEFFindHandlerT) OverrideOnFindResult(fn uintptr) { v.OnFindResult = fn }

func (v *CEFFindHandlerT) CallOnFindResult(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnFindResult, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterFindHandler(handle uintptr) {
}
