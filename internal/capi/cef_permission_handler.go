package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFMediaAccessCallbackT struct {
	_      structs.HostLayout
	Base   CEFBaseRefCountedT
	Cont   uintptr
	Cancel uintptr
}

func (v *CEFMediaAccessCallbackT) OverrideCont(fn uintptr) { v.Cont = fn }

func (v *CEFMediaAccessCallbackT) CallCont(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cont, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMediaAccessCallbackT) OverrideCancel(fn uintptr) { v.Cancel = fn }

func (v *CEFMediaAccessCallbackT) CallCancel(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cancel, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFPermissionPromptCallbackT struct {
	_    structs.HostLayout
	Base CEFBaseRefCountedT
	Cont uintptr
}

func (v *CEFPermissionPromptCallbackT) OverrideCont(fn uintptr) { v.Cont = fn }

func (v *CEFPermissionPromptCallbackT) CallCont(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cont, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFPermissionHandlerT struct {
	_                              structs.HostLayout
	Base                           CEFBaseRefCountedT
	OnRequestMediaAccessPermission uintptr
	OnShowPermissionPrompt         uintptr
	OnDismissPermissionPrompt      uintptr
}

func (v *CEFPermissionHandlerT) OverrideOnRequestMediaAccessPermission(fn uintptr) {
	v.OnRequestMediaAccessPermission = fn
}

func (v *CEFPermissionHandlerT) CallOnRequestMediaAccessPermission(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnRequestMediaAccessPermission, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPermissionHandlerT) OverrideOnShowPermissionPrompt(fn uintptr) {
	v.OnShowPermissionPrompt = fn
}

func (v *CEFPermissionHandlerT) CallOnShowPermissionPrompt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnShowPermissionPrompt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPermissionHandlerT) OverrideOnDismissPermissionPrompt(fn uintptr) {
	v.OnDismissPermissionPrompt = fn
}

func (v *CEFPermissionHandlerT) CallOnDismissPermissionPrompt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnDismissPermissionPrompt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterPermissionHandler(handle uintptr) {
}
