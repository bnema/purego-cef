package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFMenuButtonT struct {
	_           structs.HostLayout
	Base        CEFLabelButtonT
	ShowMenu    uintptr
	TriggerMenu uintptr
}

func (v *CEFMenuButtonT) OverrideShowMenu(fn uintptr) { v.ShowMenu = fn }

func (v *CEFMenuButtonT) CallShowMenu(args ...uintptr) uintptr {
	if v.ShowMenu == 0 {
		return 0
	}
	r1, _, _ := purego.SyscallN(v.ShowMenu, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuButtonT) OverrideTriggerMenu(fn uintptr) { v.TriggerMenu = fn }

func (v *CEFMenuButtonT) CallTriggerMenu(args ...uintptr) uintptr {
	if v.TriggerMenu == 0 {
		return 0
	}
	r1, _, _ := purego.SyscallN(v.TriggerMenu, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFMenuButtonCreate func(Delegate unsafe.Pointer, Text unsafe.Pointer) unsafe.Pointer

func RegisterMenuButton(handle uintptr) {
	purego.RegisterLibFunc(&CEFMenuButtonCreate, handle, "cef_menu_button_create")
}
