package capi

import "structs"

// CEFStringT matches cef_string_utf16_t (aka cef_string_t on all platforms).
// Hand-written because cef_string_types.h is not a capi header.
type CEFStringT struct {
	_      structs.HostLayout
	Str    *uint16
	Length uintptr
	Dtor   uintptr
}
