package capi

import (
	"structs"
)

type CEFAcceleratedPaintInfoCommonT struct {
	_                    structs.HostLayout
	Size                 uintptr
	Timestamp            uint64
	CodedSize            CEFSizeT
	VisibleRect          CEFRectT
	ContentRect          CEFRectT
	SourceSize           CEFSizeT
	CaptureUpdateRect    CEFRectT
	RegionCaptureRect    CEFRectT
	CaptureCounter       uint64
	HasCaptureUpdateRect uintptr
	HasRegionCaptureRect uintptr
	HasSourceSize        uintptr
	HasCaptureCounter    uintptr
}

func RegisterTypesOsr(handle uintptr) {
}
