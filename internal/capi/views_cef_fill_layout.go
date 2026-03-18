package capi

import (
	"structs"
)

type CEFFillLayoutT struct {
	_    structs.HostLayout
	Base CEFLayoutT
}

func RegisterFillLayout(handle uintptr) {
}
