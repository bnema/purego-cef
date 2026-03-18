package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFTextfieldT struct {
	_                           structs.HostLayout
	Base                        CEFViewT
	SetPasswordInput            uintptr
	IsPasswordInput             uintptr
	SetReadOnly                 uintptr
	IsReadOnly                  uintptr
	GetText                     uintptr
	SetText                     uintptr
	AppendText                  uintptr
	InsertOrReplaceText         uintptr
	HasSelection                uintptr
	GetSelectedText             uintptr
	SelectAll                   uintptr
	ClearSelection              uintptr
	GetSelectedRange            uintptr
	SelectRange                 uintptr
	GetCursorPosition           uintptr
	SetTextColor                uintptr
	GetTextColor                uintptr
	SetSelectionTextColor       uintptr
	GetSelectionTextColor       uintptr
	SetSelectionBackgroundColor uintptr
	GetSelectionBackgroundColor uintptr
	SetFontList                 uintptr
	ApplyTextColor              uintptr
	ApplyTextStyle              uintptr
	IsCommandEnabled            uintptr
	ExecuteCommand              uintptr
	ClearEditHistory            uintptr
	SetPlaceholderText          uintptr
	GetPlaceholderText          uintptr
	SetPlaceholderTextColor     uintptr
	SetAccessibleName           uintptr
}

func (v *CEFTextfieldT) OverrideSetPasswordInput(fn uintptr) { v.SetPasswordInput = fn }

func (v *CEFTextfieldT) CallSetPasswordInput(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetPasswordInput, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideIsPasswordInput(fn uintptr) { v.IsPasswordInput = fn }

func (v *CEFTextfieldT) CallIsPasswordInput(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsPasswordInput, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSetReadOnly(fn uintptr) { v.SetReadOnly = fn }

func (v *CEFTextfieldT) CallSetReadOnly(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetReadOnly, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideIsReadOnly(fn uintptr) { v.IsReadOnly = fn }

func (v *CEFTextfieldT) CallIsReadOnly(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsReadOnly, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideGetText(fn uintptr) { v.GetText = fn }

func (v *CEFTextfieldT) CallGetText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSetText(fn uintptr) { v.SetText = fn }

func (v *CEFTextfieldT) CallSetText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideAppendText(fn uintptr) { v.AppendText = fn }

func (v *CEFTextfieldT) CallAppendText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AppendText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideInsertOrReplaceText(fn uintptr) { v.InsertOrReplaceText = fn }

func (v *CEFTextfieldT) CallInsertOrReplaceText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.InsertOrReplaceText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideHasSelection(fn uintptr) { v.HasSelection = fn }

func (v *CEFTextfieldT) CallHasSelection(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.HasSelection, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideGetSelectedText(fn uintptr) { v.GetSelectedText = fn }

func (v *CEFTextfieldT) CallGetSelectedText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSelectedText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSelectAll(fn uintptr) { v.SelectAll = fn }

func (v *CEFTextfieldT) CallSelectAll(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SelectAll, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideClearSelection(fn uintptr) { v.ClearSelection = fn }

func (v *CEFTextfieldT) CallClearSelection(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ClearSelection, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideGetSelectedRange(fn uintptr) { v.GetSelectedRange = fn }

func (v *CEFTextfieldT) CallGetSelectedRange(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSelectedRange, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSelectRange(fn uintptr) { v.SelectRange = fn }

func (v *CEFTextfieldT) CallSelectRange(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SelectRange, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideGetCursorPosition(fn uintptr) { v.GetCursorPosition = fn }

func (v *CEFTextfieldT) CallGetCursorPosition(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetCursorPosition, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSetTextColor(fn uintptr) { v.SetTextColor = fn }

func (v *CEFTextfieldT) CallSetTextColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetTextColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideGetTextColor(fn uintptr) { v.GetTextColor = fn }

func (v *CEFTextfieldT) CallGetTextColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetTextColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSetSelectionTextColor(fn uintptr) { v.SetSelectionTextColor = fn }

func (v *CEFTextfieldT) CallSetSelectionTextColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetSelectionTextColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideGetSelectionTextColor(fn uintptr) { v.GetSelectionTextColor = fn }

func (v *CEFTextfieldT) CallGetSelectionTextColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSelectionTextColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSetSelectionBackgroundColor(fn uintptr) {
	v.SetSelectionBackgroundColor = fn
}

func (v *CEFTextfieldT) CallSetSelectionBackgroundColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetSelectionBackgroundColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideGetSelectionBackgroundColor(fn uintptr) {
	v.GetSelectionBackgroundColor = fn
}

func (v *CEFTextfieldT) CallGetSelectionBackgroundColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSelectionBackgroundColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSetFontList(fn uintptr) { v.SetFontList = fn }

func (v *CEFTextfieldT) CallSetFontList(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetFontList, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideApplyTextColor(fn uintptr) { v.ApplyTextColor = fn }

func (v *CEFTextfieldT) CallApplyTextColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ApplyTextColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideApplyTextStyle(fn uintptr) { v.ApplyTextStyle = fn }

func (v *CEFTextfieldT) CallApplyTextStyle(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ApplyTextStyle, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideIsCommandEnabled(fn uintptr) { v.IsCommandEnabled = fn }

func (v *CEFTextfieldT) CallIsCommandEnabled(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsCommandEnabled, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideExecuteCommand(fn uintptr) { v.ExecuteCommand = fn }

func (v *CEFTextfieldT) CallExecuteCommand(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ExecuteCommand, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideClearEditHistory(fn uintptr) { v.ClearEditHistory = fn }

func (v *CEFTextfieldT) CallClearEditHistory(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.ClearEditHistory, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSetPlaceholderText(fn uintptr) { v.SetPlaceholderText = fn }

func (v *CEFTextfieldT) CallSetPlaceholderText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetPlaceholderText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideGetPlaceholderText(fn uintptr) { v.GetPlaceholderText = fn }

func (v *CEFTextfieldT) CallGetPlaceholderText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPlaceholderText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSetPlaceholderTextColor(fn uintptr) { v.SetPlaceholderTextColor = fn }

func (v *CEFTextfieldT) CallSetPlaceholderTextColor(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetPlaceholderTextColor, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFTextfieldT) OverrideSetAccessibleName(fn uintptr) { v.SetAccessibleName = fn }

func (v *CEFTextfieldT) CallSetAccessibleName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetAccessibleName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFTextfieldCreate func(Delegate unsafe.Pointer) unsafe.Pointer

func RegisterTextfield(handle uintptr) {
	purego.RegisterLibFunc(&CEFTextfieldCreate, handle, "cef_textfield_create")
}
