package capi

import (
	"structs"
)

type CEFRegistrationT struct {
	_    structs.HostLayout
	Base CEFBaseRefCountedT
}

func RegisterRegistration(_ uintptr) {}
