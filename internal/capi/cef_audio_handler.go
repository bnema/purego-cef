package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFAudioHandlerT struct {
	_                    structs.HostLayout
	Base                 CEFBaseRefCountedT
	GetAudioParameters   uintptr
	OnAudioStreamStarted uintptr
	OnAudioStreamPacket  uintptr
	OnAudioStreamStopped uintptr
	OnAudioStreamError   uintptr
}

func (v *CEFAudioHandlerT) OverrideGetAudioParameters(fn uintptr) { v.GetAudioParameters = fn }

func (v *CEFAudioHandlerT) CallGetAudioParameters(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetAudioParameters, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFAudioHandlerT) OverrideOnAudioStreamStarted(fn uintptr) { v.OnAudioStreamStarted = fn }

func (v *CEFAudioHandlerT) CallOnAudioStreamStarted(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnAudioStreamStarted, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFAudioHandlerT) OverrideOnAudioStreamPacket(fn uintptr) { v.OnAudioStreamPacket = fn }

func (v *CEFAudioHandlerT) CallOnAudioStreamPacket(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnAudioStreamPacket, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFAudioHandlerT) OverrideOnAudioStreamStopped(fn uintptr) { v.OnAudioStreamStopped = fn }

func (v *CEFAudioHandlerT) CallOnAudioStreamStopped(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnAudioStreamStopped, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFAudioHandlerT) OverrideOnAudioStreamError(fn uintptr) { v.OnAudioStreamError = fn }

func (v *CEFAudioHandlerT) CallOnAudioStreamError(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.OnAudioStreamError, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func RegisterAudioHandler(handle uintptr) {
}
