package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFSharedProcessMessageBuilderT struct {
	_       structs.HostLayout
	Base    CEFBaseRefCountedT
	IsValid uintptr
	Size    uintptr
	Memory  uintptr
	Build   uintptr
}

func (v *CEFSharedProcessMessageBuilderT) OverrideIsValid(fn uintptr) { v.IsValid = fn }

func (v *CEFSharedProcessMessageBuilderT) CallIsValid(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsValid, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFSharedProcessMessageBuilderT) OverrideSize(fn uintptr) { v.Size = fn }

func (v *CEFSharedProcessMessageBuilderT) CallSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Size, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFSharedProcessMessageBuilderT) OverrideMemory(fn uintptr) { v.Memory = fn }

func (v *CEFSharedProcessMessageBuilderT) CallMemory(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Memory, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFSharedProcessMessageBuilderT) OverrideBuild(fn uintptr) { v.Build = fn }

func (v *CEFSharedProcessMessageBuilderT) CallBuild(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Build, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFSharedProcessMessageBuilderCreate func(Name unsafe.Pointer, ByteSize uintptr) unsafe.Pointer

func RegisterSharedProcessMessageBuilder(handle uintptr) {
	purego.RegisterLibFunc(&CEFSharedProcessMessageBuilderCreate, handle, "cef_shared_process_message_builder_create")
}
