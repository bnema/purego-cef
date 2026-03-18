package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFMenuModelT struct {
	_                   structs.HostLayout
	Base                CEFBaseRefCountedT
	IsSubMenu           uintptr
	Clear               uintptr
	GetCount            uintptr
	AddSeparator        uintptr
	AddItem             uintptr
	AddCheckItem        uintptr
	AddRadioItem        uintptr
	AddSubMenu          uintptr
	InsertSeparatorAt   uintptr
	InsertItemAt        uintptr
	InsertCheckItemAt   uintptr
	InsertRadioItemAt   uintptr
	InsertSubMenuAt     uintptr
	Remove              uintptr
	RemoveAt            uintptr
	GetIndexOf          uintptr
	GetCommandIDAt      uintptr
	SetCommandIDAt      uintptr
	GetLabel            uintptr
	GetLabelAt          uintptr
	SetLabel            uintptr
	SetLabelAt          uintptr
	GetType             uintptr
	GetTypeAt           uintptr
	GetGroupID          uintptr
	GetGroupIDAt        uintptr
	SetGroupID          uintptr
	SetGroupIDAt        uintptr
	GetSubMenu          uintptr
	GetSubMenuAt        uintptr
	IsVisible           uintptr
	IsVisibleAt         uintptr
	SetVisible          uintptr
	SetVisibleAt        uintptr
	IsEnabled           uintptr
	IsEnabledAt         uintptr
	SetEnabled          uintptr
	SetEnabledAt        uintptr
	IsChecked           uintptr
	IsCheckedAt         uintptr
	SetChecked          uintptr
	SetCheckedAt        uintptr
	HasAccelerator      uintptr
	HasAcceleratorAt    uintptr
	SetAccelerator      uintptr
	SetAcceleratorAt    uintptr
	RemoveAccelerator   uintptr
	RemoveAcceleratorAt uintptr
	GetAccelerator      uintptr
	GetAcceleratorAt    uintptr
	SetColor            uintptr
	SetColorAt          uintptr
	GetColor            uintptr
	GetColorAt          uintptr
	SetFontList         uintptr
	SetFontListAt       uintptr
}

func (v *CEFMenuModelT) OverrideIsSubMenu(fn uintptr) { v.IsSubMenu = fn }

func (v *CEFMenuModelT) CallIsSubMenu(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsSubMenu, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideClear(fn uintptr) { v.Clear = fn }

func (v *CEFMenuModelT) CallClear(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Clear, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetCount(fn uintptr) { v.GetCount = fn }

func (v *CEFMenuModelT) CallGetCount(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetCount, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideAddSeparator(fn uintptr) { v.AddSeparator = fn }

func (v *CEFMenuModelT) CallAddSeparator(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AddSeparator, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideAddItem(fn uintptr) { v.AddItem = fn }

func (v *CEFMenuModelT) CallAddItem(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AddItem, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideAddCheckItem(fn uintptr) { v.AddCheckItem = fn }

func (v *CEFMenuModelT) CallAddCheckItem(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AddCheckItem, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideAddRadioItem(fn uintptr) { v.AddRadioItem = fn }

func (v *CEFMenuModelT) CallAddRadioItem(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AddRadioItem, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideAddSubMenu(fn uintptr) { v.AddSubMenu = fn }

func (v *CEFMenuModelT) CallAddSubMenu(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AddSubMenu, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideInsertSeparatorAt(fn uintptr) { v.InsertSeparatorAt = fn }

func (v *CEFMenuModelT) CallInsertSeparatorAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.InsertSeparatorAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideInsertItemAt(fn uintptr) { v.InsertItemAt = fn }

func (v *CEFMenuModelT) CallInsertItemAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.InsertItemAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideInsertCheckItemAt(fn uintptr) { v.InsertCheckItemAt = fn }

func (v *CEFMenuModelT) CallInsertCheckItemAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.InsertCheckItemAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideInsertRadioItemAt(fn uintptr) { v.InsertRadioItemAt = fn }

func (v *CEFMenuModelT) CallInsertRadioItemAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.InsertRadioItemAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideInsertSubMenuAt(fn uintptr) { v.InsertSubMenuAt = fn }

func (v *CEFMenuModelT) CallInsertSubMenuAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.InsertSubMenuAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideRemove(fn uintptr) { v.Remove = fn }

func (v *CEFMenuModelT) CallRemove(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Remove, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideRemoveAt(fn uintptr) { v.RemoveAt = fn }

func (v *CEFMenuModelT) CallRemoveAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.RemoveAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetIndexOf(fn uintptr) { v.GetIndexOf = fn }

func (v *CEFMenuModelT) CallGetIndexOf(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetIndexOf, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetCommandIDAt(fn uintptr) { v.GetCommandIDAt = fn }

func (v *CEFMenuModelT) CallGetCommandIDAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetCommandIDAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetCommandIDAt(fn uintptr) { v.SetCommandIDAt = fn }

func (v *CEFMenuModelT) CallSetCommandIDAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetCommandIDAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetLabel(fn uintptr) { v.GetLabel = fn }

func (v *CEFMenuModelT) CallGetLabel(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetLabel, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetLabelAt(fn uintptr) { v.GetLabelAt = fn }

func (v *CEFMenuModelT) CallGetLabelAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetLabelAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetLabel(fn uintptr) { v.SetLabel = fn }

func (v *CEFMenuModelT) CallSetLabel(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetLabel, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetLabelAt(fn uintptr) { v.SetLabelAt = fn }

func (v *CEFMenuModelT) CallSetLabelAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetLabelAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetType(fn uintptr) { v.GetType = fn }

func (v *CEFMenuModelT) CallGetType(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetType, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetTypeAt(fn uintptr) { v.GetTypeAt = fn }

func (v *CEFMenuModelT) CallGetTypeAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetTypeAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetGroupID(fn uintptr) { v.GetGroupID = fn }

func (v *CEFMenuModelT) CallGetGroupID(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetGroupID, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetGroupIDAt(fn uintptr) { v.GetGroupIDAt = fn }

func (v *CEFMenuModelT) CallGetGroupIDAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetGroupIDAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetGroupID(fn uintptr) { v.SetGroupID = fn }

func (v *CEFMenuModelT) CallSetGroupID(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetGroupID, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetGroupIDAt(fn uintptr) { v.SetGroupIDAt = fn }

func (v *CEFMenuModelT) CallSetGroupIDAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetGroupIDAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetSubMenu(fn uintptr) { v.GetSubMenu = fn }

func (v *CEFMenuModelT) CallGetSubMenu(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSubMenu, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetSubMenuAt(fn uintptr) { v.GetSubMenuAt = fn }

func (v *CEFMenuModelT) CallGetSubMenuAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSubMenuAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideIsVisible(fn uintptr) { v.IsVisible = fn }

func (v *CEFMenuModelT) CallIsVisible(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsVisible, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideIsVisibleAt(fn uintptr) { v.IsVisibleAt = fn }

func (v *CEFMenuModelT) CallIsVisibleAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsVisibleAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetVisible(fn uintptr) { v.SetVisible = fn }

func (v *CEFMenuModelT) CallSetVisible(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetVisible, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetVisibleAt(fn uintptr) { v.SetVisibleAt = fn }

func (v *CEFMenuModelT) CallSetVisibleAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetVisibleAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideIsEnabled(fn uintptr) { v.IsEnabled = fn }

func (v *CEFMenuModelT) CallIsEnabled(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsEnabled, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideIsEnabledAt(fn uintptr) { v.IsEnabledAt = fn }

func (v *CEFMenuModelT) CallIsEnabledAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsEnabledAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetEnabled(fn uintptr) { v.SetEnabled = fn }

func (v *CEFMenuModelT) CallSetEnabled(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetEnabled, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetEnabledAt(fn uintptr) { v.SetEnabledAt = fn }

func (v *CEFMenuModelT) CallSetEnabledAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetEnabledAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideIsChecked(fn uintptr) { v.IsChecked = fn }

func (v *CEFMenuModelT) CallIsChecked(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsChecked, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideIsCheckedAt(fn uintptr) { v.IsCheckedAt = fn }

func (v *CEFMenuModelT) CallIsCheckedAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsCheckedAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetChecked(fn uintptr) { v.SetChecked = fn }

func (v *CEFMenuModelT) CallSetChecked(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetChecked, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetCheckedAt(fn uintptr) { v.SetCheckedAt = fn }

func (v *CEFMenuModelT) CallSetCheckedAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetCheckedAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideHasAccelerator(fn uintptr) { v.HasAccelerator = fn }

func (v *CEFMenuModelT) CallHasAccelerator(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.HasAccelerator, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideHasAcceleratorAt(fn uintptr) { v.HasAcceleratorAt = fn }

func (v *CEFMenuModelT) CallHasAcceleratorAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.HasAcceleratorAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetAccelerator(fn uintptr) { v.SetAccelerator = fn }

func (v *CEFMenuModelT) CallSetAccelerator(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetAccelerator, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetAcceleratorAt(fn uintptr) { v.SetAcceleratorAt = fn }

func (v *CEFMenuModelT) CallSetAcceleratorAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetAcceleratorAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideRemoveAccelerator(fn uintptr) { v.RemoveAccelerator = fn }

func (v *CEFMenuModelT) CallRemoveAccelerator(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.RemoveAccelerator, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideRemoveAcceleratorAt(fn uintptr) { v.RemoveAcceleratorAt = fn }

func (v *CEFMenuModelT) CallRemoveAcceleratorAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.RemoveAcceleratorAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetAccelerator(fn uintptr) { v.GetAccelerator = fn }

func (v *CEFMenuModelT) CallGetAccelerator(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetAccelerator, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetAcceleratorAt(fn uintptr) { v.GetAcceleratorAt = fn }

func (v *CEFMenuModelT) CallGetAcceleratorAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetAcceleratorAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetColor(fn uintptr) { v.SetColor = fn }

func (v *CEFMenuModelT) CallSetColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetColorAt(fn uintptr) { v.SetColorAt = fn }

func (v *CEFMenuModelT) CallSetColorAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetColorAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetColor(fn uintptr) { v.GetColor = fn }

func (v *CEFMenuModelT) CallGetColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideGetColorAt(fn uintptr) { v.GetColorAt = fn }

func (v *CEFMenuModelT) CallGetColorAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetColorAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetFontList(fn uintptr) { v.SetFontList = fn }

func (v *CEFMenuModelT) CallSetFontList(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetFontList, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFMenuModelT) OverrideSetFontListAt(fn uintptr) { v.SetFontListAt = fn }

func (v *CEFMenuModelT) CallSetFontListAt(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetFontListAt, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFMenuModelCreate func(Delegate unsafe.Pointer) unsafe.Pointer

func RegisterMenuModel(handle uintptr) {
	purego.RegisterLibFunc(&CEFMenuModelCreate, handle, "cef_menu_model_create")
}
