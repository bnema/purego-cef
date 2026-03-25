// adapter.go implements portout.CAPI using purego bindings.
// This is the only handwritten file in the capi package.
package capi

import (
	"unsafe"

	portout "github.com/bnema/purego-cef/internal/ports/out"
	"github.com/ebitengine/purego"
)

// Adapter implements portout.CAPI using purego bindings.
type Adapter struct {
	handle uintptr

	// String functions bound from CEF shared library.
	stringSet     func(*uint16, uintptr, unsafe.Pointer, int32) int32
	stringClear   func(unsafe.Pointer)
	stringFreeFn  func(unsafe.Pointer)
	stringListSz  func(uintptr) uintptr
	stringListVal func(uintptr, uintptr, unsafe.Pointer) int32

	// Core functions bound from CEF shared library (matches out.AppFunctions).
	initialize              func(unsafe.Pointer, unsafe.Pointer, unsafe.Pointer, unsafe.Pointer) int32
	shutdown                func()
	doMsgLoopWork           func()
	executeProcess          func(unsafe.Pointer, unsafe.Pointer, unsafe.Pointer) int32
	getExitCode             func() int32
	runMessageLoop          func()
	quitMessageLoop         func()
	setNestableTasksAllowed func(int32)
}

// New creates a CAPI Adapter, registering all symbols from the handle.
func New(handle uintptr) *Adapter {
	a := &Adapter{handle: handle}
	Register(handle)
	a.bindStringFuncs(handle)
	a.bindCoreFuncs(handle)
	return a
}

func (a *Adapter) bindStringFuncs(handle uintptr) {
	purego.RegisterLibFunc(&a.stringSet, handle, "cef_string_utf16_set")
	purego.RegisterLibFunc(&a.stringClear, handle, "cef_string_utf16_clear")
	purego.RegisterLibFunc(&a.stringFreeFn, handle, "cef_string_userfree_utf16_free")
	purego.RegisterLibFunc(&a.stringListSz, handle, "cef_string_list_size")
	purego.RegisterLibFunc(&a.stringListVal, handle, "cef_string_list_value")
}

func (a *Adapter) bindCoreFuncs(handle uintptr) {
	purego.RegisterLibFunc(&a.initialize, handle, "cef_initialize")
	purego.RegisterLibFunc(&a.shutdown, handle, "cef_shutdown")
	purego.RegisterLibFunc(&a.doMsgLoopWork, handle, "cef_do_message_loop_work")
	purego.RegisterLibFunc(&a.executeProcess, handle, "cef_execute_process")
	purego.RegisterLibFunc(&a.getExitCode, handle, "cef_get_exit_code")
	purego.RegisterLibFunc(&a.runMessageLoop, handle, "cef_run_message_loop")
	purego.RegisterLibFunc(&a.quitMessageLoop, handle, "cef_quit_message_loop")
	purego.RegisterLibFunc(&a.setNestableTasksAllowed, handle, "cef_set_nestable_tasks_allowed")
}

func (a *Adapter) Initialize(args, settings, app, sandboxInfo unsafe.Pointer) int32 {
	return a.initialize(args, settings, app, sandboxInfo)
}

func (a *Adapter) Shutdown() {
	a.shutdown()
}

func (a *Adapter) DoMessageLoopWork() {
	a.doMsgLoopWork()
}

func (a *Adapter) ExecuteProcess(args, app, sandboxInfo unsafe.Pointer) int32 {
	return a.executeProcess(args, app, sandboxInfo)
}

func (a *Adapter) GetExitCode() int32 {
	return a.getExitCode()
}

func (a *Adapter) RunMessageLoop() {
	a.runMessageLoop()
}

func (a *Adapter) QuitMessageLoop() {
	a.quitMessageLoop()
}

// SetNestableTasksAllowed adapts the generated port signature (unsafe.Pointer)
// to the actual C signature (int). The port-out template uses unsafe.Pointer
// for all free function params — this is a known limitation.
func (a *Adapter) SetNestableTasksAllowed(allowed unsafe.Pointer) {
	a.setNestableTasksAllowed(int32(uintptr(allowed)))
}

func (a *Adapter) NewCallback(fn any) uintptr {
	return purego.NewCallback(fn)
}

func (a *Adapter) StringSet(src *uint16, srcLen uintptr, output unsafe.Pointer, cp int32) int32 {
	return a.stringSet(src, srcLen, output, cp)
}

func (a *Adapter) StringClear(s unsafe.Pointer) {
	a.stringClear(s)
}

func (a *Adapter) StringUserfreeFree(s unsafe.Pointer) {
	a.stringFreeFn(s)
}

func (a *Adapter) StringListSize(list uintptr) uintptr {
	return a.stringListSz(list)
}

func (a *Adapter) StringListValue(list uintptr, index uintptr, value unsafe.Pointer) int32 {
	return a.stringListVal(list, index, value)
}

// Ensure Adapter implements portout.CAPI at compile time.
var _ portout.CAPI = (*Adapter)(nil)
