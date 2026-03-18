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

// CEFColorT is typedef uint32_t cef_color_t (from cef_types.h).
// Hand-written because the parser doesn't handle simple typedef aliases.
type CEFColorT = uint32

// CEFWindowHandleT is typedef unsigned long cef_window_handle_t (linux).
// Hand-written because it comes from cef_types_linux.h as a simple typedef.
type CEFWindowHandleT = uint64

// CEFBasetimeT matches cef_basetime_t (from cef_time.h).
// Hand-written because cef_time.h is not a capi header.
type CEFBasetimeT struct {
	_   structs.HostLayout
	Val int64
}

// kAcceleratedPaintMaxPlanes is #define kAcceleratedPaintMaxPlanes 4 (from cef_types_linux.h).
const kacceleratedpaintmaxplanes = 4
