package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFDevToolsMessageObserverT struct {
	_                       structs.HostLayout
	Base                    CEFBaseRefCountedT
	OnDevToolsMessage       uintptr
	OnDevToolsMethodResult  uintptr
	OnDevToolsEvent         uintptr
	OnDevToolsAgentAttached uintptr
	OnDevToolsAgentDetached uintptr
}

func (v *CEFDevToolsMessageObserverT) OverrideOnDevToolsMessage(fn uintptr) { v.OnDevToolsMessage = fn }

func (v *CEFDevToolsMessageObserverT) CallOnDevToolsMessage(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnDevToolsMessage, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFDevToolsMessageObserverT) OverrideOnDevToolsMethodResult(fn uintptr) {
	v.OnDevToolsMethodResult = fn
}

func (v *CEFDevToolsMessageObserverT) CallOnDevToolsMethodResult(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnDevToolsMethodResult, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFDevToolsMessageObserverT) OverrideOnDevToolsEvent(fn uintptr) { v.OnDevToolsEvent = fn }

func (v *CEFDevToolsMessageObserverT) CallOnDevToolsEvent(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnDevToolsEvent, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFDevToolsMessageObserverT) OverrideOnDevToolsAgentAttached(fn uintptr) {
	v.OnDevToolsAgentAttached = fn
}

func (v *CEFDevToolsMessageObserverT) CallOnDevToolsAgentAttached(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnDevToolsAgentAttached, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFDevToolsMessageObserverT) OverrideOnDevToolsAgentDetached(fn uintptr) {
	v.OnDevToolsAgentDetached = fn
}

func (v *CEFDevToolsMessageObserverT) CallOnDevToolsAgentDetached(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnDevToolsAgentDetached, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterDevtoolsMessageObserver(handle uintptr) {
}
