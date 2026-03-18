package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFSharedMemoryRegionT struct {
	_       structs.HostLayout
	Base    CEFBaseRefCountedT
	IsValid uintptr
	Size    uintptr
	Memory  uintptr
}

func (v *CEFSharedMemoryRegionT) OverrideIsValid(fn uintptr) { v.IsValid = fn }

func (v *CEFSharedMemoryRegionT) CallIsValid(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsValid, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFSharedMemoryRegionT) OverrideSize(fn uintptr) { v.Size = fn }

func (v *CEFSharedMemoryRegionT) CallSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Size, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFSharedMemoryRegionT) OverrideMemory(fn uintptr) { v.Memory = fn }

func (v *CEFSharedMemoryRegionT) CallMemory(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Memory, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterSharedMemoryRegion(handle uintptr) {
}
