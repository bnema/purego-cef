package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFResponseT struct {
	_               structs.HostLayout
	Base            CEFBaseRefCountedT
	IsReadOnly      uintptr
	GetError        uintptr
	SetError        uintptr
	GetStatus       uintptr
	SetStatus       uintptr
	GetStatusText   uintptr
	SetStatusText   uintptr
	GetMimeType     uintptr
	SetMimeType     uintptr
	GetCharset      uintptr
	SetCharset      uintptr
	GetHeaderByName uintptr
	SetHeaderByName uintptr
	GetHeaderMap    uintptr
	SetHeaderMap    uintptr
	GetURL          uintptr
	SetURL          uintptr
}

func (v *CEFResponseT) OverrideIsReadOnly(fn uintptr) { v.IsReadOnly = fn }

func (v *CEFResponseT) CallIsReadOnly(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsReadOnly, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideGetError(fn uintptr) { v.GetError = fn }

func (v *CEFResponseT) CallGetError(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetError, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideSetError(fn uintptr) { v.SetError = fn }

func (v *CEFResponseT) CallSetError(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetError, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideGetStatus(fn uintptr) { v.GetStatus = fn }

func (v *CEFResponseT) CallGetStatus(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetStatus, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideSetStatus(fn uintptr) { v.SetStatus = fn }

func (v *CEFResponseT) CallSetStatus(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetStatus, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideGetStatusText(fn uintptr) { v.GetStatusText = fn }

func (v *CEFResponseT) CallGetStatusText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetStatusText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideSetStatusText(fn uintptr) { v.SetStatusText = fn }

func (v *CEFResponseT) CallSetStatusText(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetStatusText, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideGetMimeType(fn uintptr) { v.GetMimeType = fn }

func (v *CEFResponseT) CallGetMimeType(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetMimeType, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideSetMimeType(fn uintptr) { v.SetMimeType = fn }

func (v *CEFResponseT) CallSetMimeType(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetMimeType, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideGetCharset(fn uintptr) { v.GetCharset = fn }

func (v *CEFResponseT) CallGetCharset(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetCharset, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideSetCharset(fn uintptr) { v.SetCharset = fn }

func (v *CEFResponseT) CallSetCharset(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetCharset, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideGetHeaderByName(fn uintptr) { v.GetHeaderByName = fn }

func (v *CEFResponseT) CallGetHeaderByName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetHeaderByName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideSetHeaderByName(fn uintptr) { v.SetHeaderByName = fn }

func (v *CEFResponseT) CallSetHeaderByName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetHeaderByName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideGetHeaderMap(fn uintptr) { v.GetHeaderMap = fn }

func (v *CEFResponseT) CallGetHeaderMap(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetHeaderMap, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideSetHeaderMap(fn uintptr) { v.SetHeaderMap = fn }

func (v *CEFResponseT) CallSetHeaderMap(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetHeaderMap, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideGetURL(fn uintptr) { v.GetURL = fn }

func (v *CEFResponseT) CallGetURL(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetURL, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFResponseT) OverrideSetURL(fn uintptr) { v.SetURL = fn }

func (v *CEFResponseT) CallSetURL(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetURL, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFResponseCreate func() unsafe.Pointer

func RegisterResponse(handle uintptr) {
	purego.RegisterLibFunc(&CEFResponseCreate, handle, "cef_response_create")
}
