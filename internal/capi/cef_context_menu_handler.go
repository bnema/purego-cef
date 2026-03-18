package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFRunContextMenuCallbackT struct {
	_      structs.HostLayout
	Base   CEFBaseRefCountedT
	Cont   uintptr
	Cancel uintptr
}

func (v *CEFRunContextMenuCallbackT) OverrideCont(fn uintptr) { v.Cont = fn }

func (v *CEFRunContextMenuCallbackT) CallCont(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cont, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRunContextMenuCallbackT) OverrideCancel(fn uintptr) { v.Cancel = fn }

func (v *CEFRunContextMenuCallbackT) CallCancel(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cancel, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFRunQuickMenuCallbackT struct {
	_      structs.HostLayout
	Base   CEFBaseRefCountedT
	Cont   uintptr
	Cancel uintptr
}

func (v *CEFRunQuickMenuCallbackT) OverrideCont(fn uintptr) { v.Cont = fn }

func (v *CEFRunQuickMenuCallbackT) CallCont(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cont, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRunQuickMenuCallbackT) OverrideCancel(fn uintptr) { v.Cancel = fn }

func (v *CEFRunQuickMenuCallbackT) CallCancel(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cancel, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFContextMenuHandlerT struct {
	_                      structs.HostLayout
	Base                   CEFBaseRefCountedT
	OnBeforeContextMenu    uintptr
	RunContextMenu         uintptr
	OnContextMenuCommand   uintptr
	OnContextMenuDismissed uintptr
	RunQuickMenu           uintptr
	OnQuickMenuCommand     uintptr
	OnQuickMenuDismissed   uintptr
}

func (v *CEFContextMenuHandlerT) OverrideOnBeforeContextMenu(fn uintptr) { v.OnBeforeContextMenu = fn }

func (v *CEFContextMenuHandlerT) CallOnBeforeContextMenu(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnBeforeContextMenu, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuHandlerT) OverrideRunContextMenu(fn uintptr) { v.RunContextMenu = fn }

func (v *CEFContextMenuHandlerT) CallRunContextMenu(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.RunContextMenu, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuHandlerT) OverrideOnContextMenuCommand(fn uintptr) {
	v.OnContextMenuCommand = fn
}

func (v *CEFContextMenuHandlerT) CallOnContextMenuCommand(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnContextMenuCommand, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuHandlerT) OverrideOnContextMenuDismissed(fn uintptr) {
	v.OnContextMenuDismissed = fn
}

func (v *CEFContextMenuHandlerT) CallOnContextMenuDismissed(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnContextMenuDismissed, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuHandlerT) OverrideRunQuickMenu(fn uintptr) { v.RunQuickMenu = fn }

func (v *CEFContextMenuHandlerT) CallRunQuickMenu(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.RunQuickMenu, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuHandlerT) OverrideOnQuickMenuCommand(fn uintptr) { v.OnQuickMenuCommand = fn }

func (v *CEFContextMenuHandlerT) CallOnQuickMenuCommand(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnQuickMenuCommand, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuHandlerT) OverrideOnQuickMenuDismissed(fn uintptr) {
	v.OnQuickMenuDismissed = fn
}

func (v *CEFContextMenuHandlerT) CallOnQuickMenuDismissed(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnQuickMenuDismissed, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFContextMenuParamsT struct {
	_                        structs.HostLayout
	Base                     CEFBaseRefCountedT
	GetXcoord                uintptr
	GetYcoord                uintptr
	GetTypeFlags             uintptr
	GetLinkURL               uintptr
	GetUnfilteredLinkURL     uintptr
	GetSourceURL             uintptr
	HasImageContents         uintptr
	GetTitleText             uintptr
	GetPageURL               uintptr
	GetFrameURL              uintptr
	GetFrameCharset          uintptr
	GetMediaType             uintptr
	GetMediaStateFlags       uintptr
	GetSelectionText         uintptr
	GetMisspelledWord        uintptr
	GetDictionarySuggestions uintptr
	IsEditable               uintptr
	IsSpellCheckEnabled      uintptr
	GetEditStateFlags        uintptr
	IsCustomMenu             uintptr
}

func (v *CEFContextMenuParamsT) OverrideGetXcoord(fn uintptr) { v.GetXcoord = fn }

func (v *CEFContextMenuParamsT) CallGetXcoord(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetXcoord, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetYcoord(fn uintptr) { v.GetYcoord = fn }

func (v *CEFContextMenuParamsT) CallGetYcoord(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetYcoord, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetTypeFlags(fn uintptr) { v.GetTypeFlags = fn }

func (v *CEFContextMenuParamsT) CallGetTypeFlags(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetTypeFlags, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetLinkURL(fn uintptr) { v.GetLinkURL = fn }

func (v *CEFContextMenuParamsT) CallGetLinkURL(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetLinkURL, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetUnfilteredLinkURL(fn uintptr) { v.GetUnfilteredLinkURL = fn }

func (v *CEFContextMenuParamsT) CallGetUnfilteredLinkURL(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetUnfilteredLinkURL, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetSourceURL(fn uintptr) { v.GetSourceURL = fn }

func (v *CEFContextMenuParamsT) CallGetSourceURL(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSourceURL, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideHasImageContents(fn uintptr) { v.HasImageContents = fn }

func (v *CEFContextMenuParamsT) CallHasImageContents(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.HasImageContents, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetTitleText(fn uintptr) { v.GetTitleText = fn }

func (v *CEFContextMenuParamsT) CallGetTitleText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetTitleText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetPageURL(fn uintptr) { v.GetPageURL = fn }

func (v *CEFContextMenuParamsT) CallGetPageURL(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPageURL, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetFrameURL(fn uintptr) { v.GetFrameURL = fn }

func (v *CEFContextMenuParamsT) CallGetFrameURL(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFrameURL, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetFrameCharset(fn uintptr) { v.GetFrameCharset = fn }

func (v *CEFContextMenuParamsT) CallGetFrameCharset(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFrameCharset, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetMediaType(fn uintptr) { v.GetMediaType = fn }

func (v *CEFContextMenuParamsT) CallGetMediaType(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetMediaType, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetMediaStateFlags(fn uintptr) { v.GetMediaStateFlags = fn }

func (v *CEFContextMenuParamsT) CallGetMediaStateFlags(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetMediaStateFlags, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetSelectionText(fn uintptr) { v.GetSelectionText = fn }

func (v *CEFContextMenuParamsT) CallGetSelectionText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetSelectionText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetMisspelledWord(fn uintptr) { v.GetMisspelledWord = fn }

func (v *CEFContextMenuParamsT) CallGetMisspelledWord(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetMisspelledWord, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetDictionarySuggestions(fn uintptr) {
	v.GetDictionarySuggestions = fn
}

func (v *CEFContextMenuParamsT) CallGetDictionarySuggestions(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetDictionarySuggestions, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideIsEditable(fn uintptr) { v.IsEditable = fn }

func (v *CEFContextMenuParamsT) CallIsEditable(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsEditable, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideIsSpellCheckEnabled(fn uintptr) { v.IsSpellCheckEnabled = fn }

func (v *CEFContextMenuParamsT) CallIsSpellCheckEnabled(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsSpellCheckEnabled, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideGetEditStateFlags(fn uintptr) { v.GetEditStateFlags = fn }

func (v *CEFContextMenuParamsT) CallGetEditStateFlags(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetEditStateFlags, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFContextMenuParamsT) OverrideIsCustomMenu(fn uintptr) { v.IsCustomMenu = fn }

func (v *CEFContextMenuParamsT) CallIsCustomMenu(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsCustomMenu, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterContextMenuHandler(handle uintptr) {
}
