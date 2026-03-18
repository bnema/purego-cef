package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFMenuButtonPressedLockT struct {
	_    structs.HostLayout
	Base CEFBaseRefCountedT
}

type CEFMenuButtonDelegateT struct {
	_                   structs.HostLayout
	Base                CEFButtonDelegateT
	OnMenuButtonPressed uintptr
}

func (v *CEFMenuButtonDelegateT) OverrideOnMenuButtonPressed(fn uintptr) { v.OnMenuButtonPressed = fn }

func (v *CEFMenuButtonDelegateT) CallOnMenuButtonPressed(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnMenuButtonPressed, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterMenuButtonDelegate(handle uintptr) {
}
