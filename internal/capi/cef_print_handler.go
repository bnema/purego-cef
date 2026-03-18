package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFPrintDialogCallbackT struct {
	_      structs.HostLayout
	Base   CEFBaseRefCountedT
	Cont   uintptr
	Cancel uintptr
}

func (v *CEFPrintDialogCallbackT) OverrideCont(fn uintptr) { v.Cont = fn }

func (v *CEFPrintDialogCallbackT) CallCont(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cont, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPrintDialogCallbackT) OverrideCancel(fn uintptr) { v.Cancel = fn }

func (v *CEFPrintDialogCallbackT) CallCancel(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cancel, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFPrintJobCallbackT struct {
	_    structs.HostLayout
	Base CEFBaseRefCountedT
	Cont uintptr
}

func (v *CEFPrintJobCallbackT) OverrideCont(fn uintptr) { v.Cont = fn }

func (v *CEFPrintJobCallbackT) CallCont(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cont, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFPrintHandlerT struct {
	_               structs.HostLayout
	Base            CEFBaseRefCountedT
	OnPrintStart    uintptr
	OnPrintSettings uintptr
	OnPrintDialog   uintptr
	OnPrintJob      uintptr
	OnPrintReset    uintptr
	GetPdfPaperSize uintptr
}

func (v *CEFPrintHandlerT) OverrideOnPrintStart(fn uintptr) { v.OnPrintStart = fn }

func (v *CEFPrintHandlerT) CallOnPrintStart(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnPrintStart, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPrintHandlerT) OverrideOnPrintSettings(fn uintptr) { v.OnPrintSettings = fn }

func (v *CEFPrintHandlerT) CallOnPrintSettings(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnPrintSettings, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPrintHandlerT) OverrideOnPrintDialog(fn uintptr) { v.OnPrintDialog = fn }

func (v *CEFPrintHandlerT) CallOnPrintDialog(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnPrintDialog, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPrintHandlerT) OverrideOnPrintJob(fn uintptr) { v.OnPrintJob = fn }

func (v *CEFPrintHandlerT) CallOnPrintJob(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnPrintJob, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPrintHandlerT) OverrideOnPrintReset(fn uintptr) { v.OnPrintReset = fn }

func (v *CEFPrintHandlerT) CallOnPrintReset(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnPrintReset, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPrintHandlerT) OverrideGetPdfPaperSize(fn uintptr) { v.GetPdfPaperSize = fn }

func (v *CEFPrintHandlerT) CallGetPdfPaperSize(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPdfPaperSize, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterPrintHandler(handle uintptr) {
}
