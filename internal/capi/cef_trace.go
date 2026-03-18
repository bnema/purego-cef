package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFEndTracingCallbackT struct {
	_                    structs.HostLayout
	Base                 CEFBaseRefCountedT
	OnEndTracingComplete uintptr
}

func (v *CEFEndTracingCallbackT) OverrideOnEndTracingComplete(fn uintptr) {
	v.OnEndTracingComplete = fn
}

func (v *CEFEndTracingCallbackT) CallOnEndTracingComplete(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnEndTracingComplete, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFBeginTracing func(Categories unsafe.Pointer, Callback unsafe.Pointer) int32

var CEFEndTracing func(TracingFile unsafe.Pointer, Callback unsafe.Pointer) int32

var CEFNowFromSystemTraceTime func() int64

func RegisterTrace(handle uintptr) {
	purego.RegisterLibFunc(&CEFBeginTracing, handle, "cef_begin_tracing")
	purego.RegisterLibFunc(&CEFEndTracing, handle, "cef_end_tracing")
	purego.RegisterLibFunc(&CEFNowFromSystemTraceTime, handle, "cef_now_from_system_trace_time")
}
