package capi

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFMainArgsT struct {
	_    structs.HostLayout
	Argc int32
	Argv uintptr
}

type CEFWindowInfoT struct {
	_                          structs.HostLayout
	Size                       uintptr
	WindowName                 CEFStringT
	Bounds                     CEFRectT
	ParentWindow               CEFWindowHandleT
	WindowlessRenderingEnabled int32
	SharedTextureEnabled       int32
	ExternalBeginFrameEnabled  int32
	Window                     CEFWindowHandleT
	RuntimeStyle               CEFRuntimeStyleT
}

type CEFAcceleratedPaintNativePixmapPlaneT struct {
	_      structs.HostLayout
	Stride uint32
	Offset uint64
	Size   uint64
	Fd     int32
}

type CEFAcceleratedPaintInfoT struct {
	_          structs.HostLayout
	Size       uintptr
	Planes     [kacceleratedpaintmaxplanes]CEFAcceleratedPaintNativePixmapPlaneT
	PlaneCount int32
	Modifier   uint64
	Format     CEFColorTypeT
	Extra      CEFAcceleratedPaintInfoCommonT
}

var CEFGetXdisplay func() unsafe.Pointer

func RegisterTypesLinux(handle uintptr) {
	purego.RegisterLibFunc(&CEFGetXdisplay, handle, "cef_get_xdisplay")
}
