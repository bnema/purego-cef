package capi

import (
	"structs"
)

type CEFPanelDelegateT struct {
	_    structs.HostLayout
	Base CEFViewDelegateT
}

func RegisterPanelDelegate(handle uintptr) {
}
