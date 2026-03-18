package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFRequestT struct {
	_                       structs.HostLayout
	Base                    CEFBaseRefCountedT
	IsReadOnly              uintptr
	GetURL                  uintptr
	SetURL                  uintptr
	GetMethod               uintptr
	SetMethod               uintptr
	SetReferrer             uintptr
	GetReferrerURL          uintptr
	GetReferrerPolicy       uintptr
	GetPostData             uintptr
	SetPostData             uintptr
	GetHeaderMap            uintptr
	SetHeaderMap            uintptr
	GetHeaderByName         uintptr
	SetHeaderByName         uintptr
	Set                     uintptr
	GetFlags                uintptr
	SetFlags                uintptr
	GetFirstPartyForCookies uintptr
	SetFirstPartyForCookies uintptr
	GetResourceType         uintptr
	GetTransitionType       uintptr
	GetIdentifier           uintptr
}

func (v *CEFRequestT) OverrideIsReadOnly(fn uintptr) { v.IsReadOnly = fn }

func (v *CEFRequestT) CallIsReadOnly(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsReadOnly, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetURL(fn uintptr) { v.GetURL = fn }

func (v *CEFRequestT) CallGetURL(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetURL, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideSetURL(fn uintptr) { v.SetURL = fn }

func (v *CEFRequestT) CallSetURL(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetURL, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetMethod(fn uintptr) { v.GetMethod = fn }

func (v *CEFRequestT) CallGetMethod(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetMethod, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideSetMethod(fn uintptr) { v.SetMethod = fn }

func (v *CEFRequestT) CallSetMethod(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetMethod, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideSetReferrer(fn uintptr) { v.SetReferrer = fn }

func (v *CEFRequestT) CallSetReferrer(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetReferrer, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetReferrerURL(fn uintptr) { v.GetReferrerURL = fn }

func (v *CEFRequestT) CallGetReferrerURL(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetReferrerURL, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetReferrerPolicy(fn uintptr) { v.GetReferrerPolicy = fn }

func (v *CEFRequestT) CallGetReferrerPolicy(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetReferrerPolicy, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetPostData(fn uintptr) { v.GetPostData = fn }

func (v *CEFRequestT) CallGetPostData(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetPostData, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideSetPostData(fn uintptr) { v.SetPostData = fn }

func (v *CEFRequestT) CallSetPostData(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetPostData, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetHeaderMap(fn uintptr) { v.GetHeaderMap = fn }

func (v *CEFRequestT) CallGetHeaderMap(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetHeaderMap, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideSetHeaderMap(fn uintptr) { v.SetHeaderMap = fn }

func (v *CEFRequestT) CallSetHeaderMap(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetHeaderMap, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetHeaderByName(fn uintptr) { v.GetHeaderByName = fn }

func (v *CEFRequestT) CallGetHeaderByName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetHeaderByName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideSetHeaderByName(fn uintptr) { v.SetHeaderByName = fn }

func (v *CEFRequestT) CallSetHeaderByName(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetHeaderByName, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideSet(fn uintptr) { v.Set = fn }

func (v *CEFRequestT) CallSet(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.Set, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetFlags(fn uintptr) { v.GetFlags = fn }

func (v *CEFRequestT) CallGetFlags(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFlags, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideSetFlags(fn uintptr) { v.SetFlags = fn }

func (v *CEFRequestT) CallSetFlags(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetFlags, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetFirstPartyForCookies(fn uintptr) { v.GetFirstPartyForCookies = fn }

func (v *CEFRequestT) CallGetFirstPartyForCookies(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFirstPartyForCookies, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideSetFirstPartyForCookies(fn uintptr) { v.SetFirstPartyForCookies = fn }

func (v *CEFRequestT) CallSetFirstPartyForCookies(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetFirstPartyForCookies, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetResourceType(fn uintptr) { v.GetResourceType = fn }

func (v *CEFRequestT) CallGetResourceType(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetResourceType, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetTransitionType(fn uintptr) { v.GetTransitionType = fn }

func (v *CEFRequestT) CallGetTransitionType(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetTransitionType, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFRequestT) OverrideGetIdentifier(fn uintptr) { v.GetIdentifier = fn }

func (v *CEFRequestT) CallGetIdentifier(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetIdentifier, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFPostDataT struct {
	_                   structs.HostLayout
	Base                CEFBaseRefCountedT
	IsReadOnly          uintptr
	HasExcludedElements uintptr
	GetElementCount     uintptr
	GetElements         uintptr
	RemoveElement       uintptr
	AddElement          uintptr
	RemoveElements      uintptr
}

func (v *CEFPostDataT) OverrideIsReadOnly(fn uintptr) { v.IsReadOnly = fn }

func (v *CEFPostDataT) CallIsReadOnly(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsReadOnly, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataT) OverrideHasExcludedElements(fn uintptr) { v.HasExcludedElements = fn }

func (v *CEFPostDataT) CallHasExcludedElements(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.HasExcludedElements, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataT) OverrideGetElementCount(fn uintptr) { v.GetElementCount = fn }

func (v *CEFPostDataT) CallGetElementCount(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetElementCount, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataT) OverrideGetElements(fn uintptr) { v.GetElements = fn }

func (v *CEFPostDataT) CallGetElements(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetElements, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataT) OverrideRemoveElement(fn uintptr) { v.RemoveElement = fn }

func (v *CEFPostDataT) CallRemoveElement(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.RemoveElement, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataT) OverrideAddElement(fn uintptr) { v.AddElement = fn }

func (v *CEFPostDataT) CallAddElement(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AddElement, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataT) OverrideRemoveElements(fn uintptr) { v.RemoveElements = fn }

func (v *CEFPostDataT) CallRemoveElements(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.RemoveElements, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

type CEFPostDataElementT struct {
	_             structs.HostLayout
	Base          CEFBaseRefCountedT
	IsReadOnly    uintptr
	SetToEmpty    uintptr
	SetToFile     uintptr
	SetToBytes    uintptr
	GetType       uintptr
	GetFile       uintptr
	GetBytesCount uintptr
	GetBytes      uintptr
}

func (v *CEFPostDataElementT) OverrideIsReadOnly(fn uintptr) { v.IsReadOnly = fn }

func (v *CEFPostDataElementT) CallIsReadOnly(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsReadOnly, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataElementT) OverrideSetToEmpty(fn uintptr) { v.SetToEmpty = fn }

func (v *CEFPostDataElementT) CallSetToEmpty(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetToEmpty, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataElementT) OverrideSetToFile(fn uintptr) { v.SetToFile = fn }

func (v *CEFPostDataElementT) CallSetToFile(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetToFile, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataElementT) OverrideSetToBytes(fn uintptr) { v.SetToBytes = fn }

func (v *CEFPostDataElementT) CallSetToBytes(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.SetToBytes, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataElementT) OverrideGetType(fn uintptr) { v.GetType = fn }

func (v *CEFPostDataElementT) CallGetType(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetType, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataElementT) OverrideGetFile(fn uintptr) { v.GetFile = fn }

func (v *CEFPostDataElementT) CallGetFile(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetFile, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataElementT) OverrideGetBytesCount(fn uintptr) { v.GetBytesCount = fn }

func (v *CEFPostDataElementT) CallGetBytesCount(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetBytesCount, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFPostDataElementT) OverrideGetBytes(fn uintptr) { v.GetBytes = fn }

func (v *CEFPostDataElementT) CallGetBytes(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetBytes, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFRequestCreate func() unsafe.Pointer

var CEFPostDataCreate func() unsafe.Pointer

var CEFPostDataElementCreate func() unsafe.Pointer

func RegisterRequest(handle uintptr) {
	purego.RegisterLibFunc(&CEFRequestCreate, handle, "cef_request_create")
	purego.RegisterLibFunc(&CEFPostDataCreate, handle, "cef_post_data_create")
	purego.RegisterLibFunc(&CEFPostDataElementCreate, handle, "cef_post_data_element_create")
}
