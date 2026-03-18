package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFRenderProcessHandlerT struct {
	_                        structs.HostLayout
	Base                     CEFBaseRefCountedT
	OnWebKitInitialized      uintptr
	OnBrowserCreated         uintptr
	OnBrowserDestroyed       uintptr
	GetLoadHandler           uintptr
	OnContextCreated         uintptr
	OnContextReleased        uintptr
	OnUncaughtException      uintptr
	OnFocusedNodeChanged     uintptr
	OnProcessMessageReceived uintptr
}

func (v *CEFRenderProcessHandlerT) OverrideOnWebKitInitialized(fn uintptr) {
	v.OnWebKitInitialized = fn
}

func (v *CEFRenderProcessHandlerT) CallOnWebKitInitialized(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnWebKitInitialized, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRenderProcessHandlerT) OverrideOnBrowserCreated(fn uintptr) { v.OnBrowserCreated = fn }

func (v *CEFRenderProcessHandlerT) CallOnBrowserCreated(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnBrowserCreated, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRenderProcessHandlerT) OverrideOnBrowserDestroyed(fn uintptr) { v.OnBrowserDestroyed = fn }

func (v *CEFRenderProcessHandlerT) CallOnBrowserDestroyed(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnBrowserDestroyed, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRenderProcessHandlerT) OverrideGetLoadHandler(fn uintptr) { v.GetLoadHandler = fn }

func (v *CEFRenderProcessHandlerT) CallGetLoadHandler(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetLoadHandler, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRenderProcessHandlerT) OverrideOnContextCreated(fn uintptr) { v.OnContextCreated = fn }

func (v *CEFRenderProcessHandlerT) CallOnContextCreated(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnContextCreated, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRenderProcessHandlerT) OverrideOnContextReleased(fn uintptr) { v.OnContextReleased = fn }

func (v *CEFRenderProcessHandlerT) CallOnContextReleased(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnContextReleased, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRenderProcessHandlerT) OverrideOnUncaughtException(fn uintptr) {
	v.OnUncaughtException = fn
}

func (v *CEFRenderProcessHandlerT) CallOnUncaughtException(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnUncaughtException, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRenderProcessHandlerT) OverrideOnFocusedNodeChanged(fn uintptr) {
	v.OnFocusedNodeChanged = fn
}

func (v *CEFRenderProcessHandlerT) CallOnFocusedNodeChanged(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnFocusedNodeChanged, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRenderProcessHandlerT) OverrideOnProcessMessageReceived(fn uintptr) {
	v.OnProcessMessageReceived = fn
}

func (v *CEFRenderProcessHandlerT) CallOnProcessMessageReceived(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnProcessMessageReceived, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterRenderProcessHandler(handle uintptr) {
}
