package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFLoadHandlerT struct {
	_                    structs.HostLayout
	Base                 CEFBaseRefCountedT
	OnLoadingStateChange uintptr
	OnLoadStart          uintptr
	OnLoadEnd            uintptr
	OnLoadError          uintptr
}

func (v *CEFLoadHandlerT) OverrideOnLoadingStateChange(fn uintptr) { v.OnLoadingStateChange = fn }

func (v *CEFLoadHandlerT) CallOnLoadingStateChange(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnLoadingStateChange, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFLoadHandlerT) OverrideOnLoadStart(fn uintptr) { v.OnLoadStart = fn }

func (v *CEFLoadHandlerT) CallOnLoadStart(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnLoadStart, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFLoadHandlerT) OverrideOnLoadEnd(fn uintptr) { v.OnLoadEnd = fn }

func (v *CEFLoadHandlerT) CallOnLoadEnd(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnLoadEnd, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFLoadHandlerT) OverrideOnLoadError(fn uintptr) { v.OnLoadError = fn }

func (v *CEFLoadHandlerT) CallOnLoadError(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnLoadError, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterLoadHandler(handle uintptr) {
}
