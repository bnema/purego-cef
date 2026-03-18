package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFProcessMessageT struct {
	_                     structs.HostLayout
	Base                  CEFBaseRefCountedT
	IsValid               uintptr
	IsReadOnly            uintptr
	Copy                  uintptr
	GetName               uintptr
	GetArgumentList       uintptr
	GetSharedMemoryRegion uintptr
}

func (v *CEFProcessMessageT) OverrideIsValid(fn uintptr) { v.IsValid = fn }

func (v *CEFProcessMessageT) CallIsValid(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsValid, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFProcessMessageT) OverrideIsReadOnly(fn uintptr) { v.IsReadOnly = fn }

func (v *CEFProcessMessageT) CallIsReadOnly(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsReadOnly, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFProcessMessageT) OverrideCopy(fn uintptr) { v.Copy = fn }

func (v *CEFProcessMessageT) CallCopy(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Copy, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFProcessMessageT) OverrideGetName(fn uintptr) { v.GetName = fn }

func (v *CEFProcessMessageT) CallGetName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFProcessMessageT) OverrideGetArgumentList(fn uintptr) { v.GetArgumentList = fn }

func (v *CEFProcessMessageT) CallGetArgumentList(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetArgumentList, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFProcessMessageT) OverrideGetSharedMemoryRegion(fn uintptr) { v.GetSharedMemoryRegion = fn }

func (v *CEFProcessMessageT) CallGetSharedMemoryRegion(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSharedMemoryRegion, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFProcessMessageCreate func(Name unsafe.Pointer) unsafe.Pointer

func RegisterProcessMessage(handle uintptr) {
	purego.RegisterLibFunc(&CEFProcessMessageCreate, handle, "cef_process_message_create")
}
