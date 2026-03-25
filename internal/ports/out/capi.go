package out

import "unsafe"

// CAPI composes all outbound port interfaces.
// It embeds the generated AppFunctions port and adds infrastructure
// methods not present in the parsed CEF headers (string ops from
// cef_string_types.h and purego callback creation).
type CAPI interface {
	// Generated outbound port — CEF free functions from cef_app_capi.h.
	AppFunctions

	// Infrastructure — not in parsed CEF headers, handwritten.
	NewCallback(fn any) uintptr
	StringSet(src *uint16, srcLen uintptr, output unsafe.Pointer, copy int32) int32
	StringClear(s unsafe.Pointer)
	StringUserfreeFree(s unsafe.Pointer)
	StringListSize(list uintptr) uintptr
	StringListValue(list uintptr, index uintptr, value unsafe.Pointer) int32
}
