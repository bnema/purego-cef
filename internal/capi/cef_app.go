package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFAppT struct {
	_                             structs.HostLayout
	Base                          CEFBaseRefCountedT
	OnBeforeCommandLineProcessing uintptr
	OnRegisterCustomSchemes       uintptr
	GetResourceBundleHandler      uintptr
	GetBrowserProcessHandler      uintptr
	GetRenderProcessHandler       uintptr
}

func (v *CEFAppT) OverrideOnBeforeCommandLineProcessing(fn uintptr) {
	v.OnBeforeCommandLineProcessing = fn
}

func (v *CEFAppT) CallOnBeforeCommandLineProcessing(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnBeforeCommandLineProcessing, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFAppT) OverrideOnRegisterCustomSchemes(fn uintptr) { v.OnRegisterCustomSchemes = fn }

func (v *CEFAppT) CallOnRegisterCustomSchemes(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnRegisterCustomSchemes, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFAppT) OverrideGetResourceBundleHandler(fn uintptr) { v.GetResourceBundleHandler = fn }

func (v *CEFAppT) CallGetResourceBundleHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetResourceBundleHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFAppT) OverrideGetBrowserProcessHandler(fn uintptr) { v.GetBrowserProcessHandler = fn }

func (v *CEFAppT) CallGetBrowserProcessHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetBrowserProcessHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFAppT) OverrideGetRenderProcessHandler(fn uintptr) { v.GetRenderProcessHandler = fn }

func (v *CEFAppT) CallGetRenderProcessHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetRenderProcessHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFExecuteProcess func(Args unsafe.Pointer, Application unsafe.Pointer, WindowsSandboxInfo unsafe.Pointer) int32

var CEFInitialize func(Args unsafe.Pointer, Settings unsafe.Pointer, Application unsafe.Pointer, WindowsSandboxInfo unsafe.Pointer) int32

var CEFGetExitCode func() int32

var CEFShutdown func()

var CEFDoMessageLoopWork func()

var CEFRunMessageLoop func()

var CEFQuitMessageLoop func()

var CEFSetNestableTasksAllowed func(Allowed int32)

func RegisterApp(handle uintptr) {
	purego.RegisterLibFunc(&CEFExecuteProcess, handle, "cef_execute_process")
	purego.RegisterLibFunc(&CEFInitialize, handle, "cef_initialize")
	purego.RegisterLibFunc(&CEFGetExitCode, handle, "cef_get_exit_code")
	purego.RegisterLibFunc(&CEFShutdown, handle, "cef_shutdown")
	purego.RegisterLibFunc(&CEFDoMessageLoopWork, handle, "cef_do_message_loop_work")
	purego.RegisterLibFunc(&CEFRunMessageLoop, handle, "cef_run_message_loop")
	purego.RegisterLibFunc(&CEFQuitMessageLoop, handle, "cef_quit_message_loop")
	purego.RegisterLibFunc(&CEFSetNestableTasksAllowed, handle, "cef_set_nestable_tasks_allowed")
}
