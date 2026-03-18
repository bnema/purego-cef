package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFJsdialogCallbackT struct {
	_    structs.HostLayout
	Base CEFBaseRefCountedT
	Cont uintptr
}

func (v *CEFJsdialogCallbackT) OverrideCont(fn uintptr) { v.Cont = fn }

func (v *CEFJsdialogCallbackT) CallCont(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cont, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFJsdialogHandlerT struct {
	_                    structs.HostLayout
	Base                 CEFBaseRefCountedT
	OnJsdialog           uintptr
	OnBeforeUnloadDialog uintptr
	OnResetDialogState   uintptr
	OnDialogClosed       uintptr
}

func (v *CEFJsdialogHandlerT) OverrideOnJsdialog(fn uintptr) { v.OnJsdialog = fn }

func (v *CEFJsdialogHandlerT) CallOnJsdialog(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnJsdialog, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFJsdialogHandlerT) OverrideOnBeforeUnloadDialog(fn uintptr) { v.OnBeforeUnloadDialog = fn }

func (v *CEFJsdialogHandlerT) CallOnBeforeUnloadDialog(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnBeforeUnloadDialog, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFJsdialogHandlerT) OverrideOnResetDialogState(fn uintptr) { v.OnResetDialogState = fn }

func (v *CEFJsdialogHandlerT) CallOnResetDialogState(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnResetDialogState, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFJsdialogHandlerT) OverrideOnDialogClosed(fn uintptr) { v.OnDialogClosed = fn }

func (v *CEFJsdialogHandlerT) CallOnDialogClosed(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnDialogClosed, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterJsdialogHandler(handle uintptr) {
}
