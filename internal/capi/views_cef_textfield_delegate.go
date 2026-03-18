package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFTextfieldDelegateT struct {
	_                 structs.HostLayout
	Base              CEFViewDelegateT
	OnKeyEvent        uintptr
	OnAfterUserAction uintptr
}

func (v *CEFTextfieldDelegateT) OverrideOnKeyEvent(fn uintptr) { v.OnKeyEvent = fn }

func (v *CEFTextfieldDelegateT) CallOnKeyEvent(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnKeyEvent, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldDelegateT) OverrideOnAfterUserAction(fn uintptr) { v.OnAfterUserAction = fn }

func (v *CEFTextfieldDelegateT) CallOnAfterUserAction(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnAfterUserAction, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterTextfieldDelegate(handle uintptr) {
}
