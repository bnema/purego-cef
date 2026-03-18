package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFLifeSpanHandlerT struct {
	_                     structs.HostLayout
	Base                  uintptr
	OnBeforePopup         uintptr
	OnBeforePopupAborted  uintptr
	OnBeforeDevToolsPopup uintptr
	OnAfterCreated        uintptr
	DoClose               uintptr
	OnBeforeClose         uintptr
}

func (v *CEFLifeSpanHandlerT) OverrideOnBeforePopup(fn uintptr) { v.OnBeforePopup = fn }

func (v *CEFLifeSpanHandlerT) CallOnBeforePopup(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnBeforePopup, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFLifeSpanHandlerT) OverrideOnBeforePopupAborted(fn uintptr) { v.OnBeforePopupAborted = fn }

func (v *CEFLifeSpanHandlerT) CallOnBeforePopupAborted(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnBeforePopupAborted, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFLifeSpanHandlerT) OverrideOnBeforeDevToolsPopup(fn uintptr) { v.OnBeforeDevToolsPopup = fn }

func (v *CEFLifeSpanHandlerT) CallOnBeforeDevToolsPopup(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnBeforeDevToolsPopup, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFLifeSpanHandlerT) OverrideOnAfterCreated(fn uintptr) { v.OnAfterCreated = fn }

func (v *CEFLifeSpanHandlerT) CallOnAfterCreated(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnAfterCreated, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFLifeSpanHandlerT) OverrideDoClose(fn uintptr) { v.DoClose = fn }

func (v *CEFLifeSpanHandlerT) CallDoClose(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.DoClose, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFLifeSpanHandlerT) OverrideOnBeforeClose(fn uintptr) { v.OnBeforeClose = fn }

func (v *CEFLifeSpanHandlerT) CallOnBeforeClose(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnBeforeClose, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterLifeSpanHandler(handle uintptr) {
}
