package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFImageT struct {
	_                     structs.HostLayout
	Base                  CEFBaseRefCountedT
	IsEmpty               uintptr
	IsSame                uintptr
	AddBitmap             uintptr
	AddPng                uintptr
	AddJpeg               uintptr
	GetWidth              uintptr
	GetHeight             uintptr
	HasRepresentation     uintptr
	RemoveRepresentation  uintptr
	GetRepresentationInfo uintptr
	GetAsBitmap           uintptr
	GetAsPng              uintptr
	GetAsJpeg             uintptr
}

func (v *CEFImageT) OverrideIsEmpty(fn uintptr) { v.IsEmpty = fn }

func (v *CEFImageT) CallIsEmpty(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsEmpty, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideIsSame(fn uintptr) { v.IsSame = fn }

func (v *CEFImageT) CallIsSame(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.IsSame, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideAddBitmap(fn uintptr) { v.AddBitmap = fn }

func (v *CEFImageT) CallAddBitmap(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AddBitmap, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideAddPng(fn uintptr) { v.AddPng = fn }

func (v *CEFImageT) CallAddPng(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AddPng, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideAddJpeg(fn uintptr) { v.AddJpeg = fn }

func (v *CEFImageT) CallAddJpeg(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.AddJpeg, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideGetWidth(fn uintptr) { v.GetWidth = fn }

func (v *CEFImageT) CallGetWidth(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetWidth, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideGetHeight(fn uintptr) { v.GetHeight = fn }

func (v *CEFImageT) CallGetHeight(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetHeight, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideHasRepresentation(fn uintptr) { v.HasRepresentation = fn }

func (v *CEFImageT) CallHasRepresentation(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.HasRepresentation, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideRemoveRepresentation(fn uintptr) { v.RemoveRepresentation = fn }

func (v *CEFImageT) CallRemoveRepresentation(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.RemoveRepresentation, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideGetRepresentationInfo(fn uintptr) { v.GetRepresentationInfo = fn }

func (v *CEFImageT) CallGetRepresentationInfo(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetRepresentationInfo, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideGetAsBitmap(fn uintptr) { v.GetAsBitmap = fn }

func (v *CEFImageT) CallGetAsBitmap(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetAsBitmap, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideGetAsPng(fn uintptr) { v.GetAsPng = fn }

func (v *CEFImageT) CallGetAsPng(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetAsPng, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

func (v *CEFImageT) OverrideGetAsJpeg(fn uintptr) { v.GetAsJpeg = fn }

func (v *CEFImageT) CallGetAsJpeg(args ...uintptr) uintptr {
	r1, _, _ := purego.SyscallN(v.GetAsJpeg, append([]uintptr{uintptr(unsafe.Pointer(v))}, args...)...)
	return r1
}

var CEFImageCreate func() unsafe.Pointer

func RegisterImage(handle uintptr) {
	purego.RegisterLibFunc(&CEFImageCreate, handle, "cef_image_create")
}
