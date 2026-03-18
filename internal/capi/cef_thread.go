package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFThreadT struct {
	_                   structs.HostLayout
	Base                CEFBaseRefCountedT
	GetTaskRunner       uintptr
	GetPlatformThreadID uintptr
	Stop                uintptr
	IsRunning           uintptr
}

func (v *CEFThreadT) OverrideGetTaskRunner(fn uintptr) { v.GetTaskRunner = fn }

func (v *CEFThreadT) CallGetTaskRunner(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetTaskRunner, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFThreadT) OverrideGetPlatformThreadID(fn uintptr) { v.GetPlatformThreadID = fn }

func (v *CEFThreadT) CallGetPlatformThreadID(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPlatformThreadID, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFThreadT) OverrideStop(fn uintptr) { v.Stop = fn }

func (v *CEFThreadT) CallStop(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Stop, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFThreadT) OverrideIsRunning(fn uintptr) { v.IsRunning = fn }

func (v *CEFThreadT) CallIsRunning(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsRunning, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFThreadCreate func(DisplayName unsafe.Pointer, Priority CEFThreadPriorityT, MessageLoopType CEFMessageLoopTypeT, Stoppable int32, ComInitMode CEFComInitModeT) unsafe.Pointer

func RegisterThread(handle uintptr) {
	purego.RegisterLibFunc(&CEFThreadCreate, handle, "cef_thread_create")
}
