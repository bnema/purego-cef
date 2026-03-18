package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFBeforeDownloadCallbackT struct {
	_    structs.HostLayout
	Base CEFBaseRefCountedT
	Cont uintptr
}

func (v *CEFBeforeDownloadCallbackT) OverrideCont(fn uintptr) { v.Cont = fn }

func (v *CEFBeforeDownloadCallbackT) CallCont(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cont, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFDownloadItemCallbackT struct {
	_      structs.HostLayout
	Base   CEFBaseRefCountedT
	Cancel uintptr
	Pause  uintptr
	Resume uintptr
}

func (v *CEFDownloadItemCallbackT) OverrideCancel(fn uintptr) { v.Cancel = fn }

func (v *CEFDownloadItemCallbackT) CallCancel(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Cancel, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFDownloadItemCallbackT) OverridePause(fn uintptr) { v.Pause = fn }

func (v *CEFDownloadItemCallbackT) CallPause(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Pause, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFDownloadItemCallbackT) OverrideResume(fn uintptr) { v.Resume = fn }

func (v *CEFDownloadItemCallbackT) CallResume(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Resume, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFDownloadHandlerT struct {
	_                 structs.HostLayout
	Base              CEFBaseRefCountedT
	CanDownload       uintptr
	OnBeforeDownload  uintptr
	OnDownloadUpdated uintptr
}

func (v *CEFDownloadHandlerT) OverrideCanDownload(fn uintptr) { v.CanDownload = fn }

func (v *CEFDownloadHandlerT) CallCanDownload(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.CanDownload, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFDownloadHandlerT) OverrideOnBeforeDownload(fn uintptr) { v.OnBeforeDownload = fn }

func (v *CEFDownloadHandlerT) CallOnBeforeDownload(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnBeforeDownload, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFDownloadHandlerT) OverrideOnDownloadUpdated(fn uintptr) { v.OnDownloadUpdated = fn }

func (v *CEFDownloadHandlerT) CallOnDownloadUpdated(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnDownloadUpdated, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterDownloadHandler(handle uintptr) {
}
