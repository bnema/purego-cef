package cefstr

import (
	"fmt"
	"runtime"
	"structs"
	"sync/atomic"
	"unicode/utf16"
	"unsafe"

	"github.com/ebitengine/purego"
)

type CEFStringUTF16 struct {
	_      structs.HostLayout
	Str    *uint16
	Length uintptr
	Dtor   uintptr
}

var (
	stringSet   func(*uint16, uintptr, *CEFStringUTF16, int32) int32
	stringClear func(*CEFStringUTF16)
	bound       atomic.Bool
)

func Bind(handle uintptr) {
	purego.RegisterLibFunc(&stringSet, handle, "cef_string_utf16_set")
	purego.RegisterLibFunc(&stringClear, handle, "cef_string_utf16_clear")
	bound.Store(true)
}

func FromGo(s string) (CEFStringUTF16, func(), error) {
	if !bound.Load() {
		return CEFStringUTF16{}, nil, fmt.Errorf("cef string helpers not bound")
	}
	encoded := utf16.Encode([]rune(s))
	var src *uint16
	if len(encoded) > 0 {
		src = &encoded[0]
	}
	var out CEFStringUTF16
	if ok := stringSet(src, uintptr(len(encoded)), &out, 1); ok != 1 {
		return CEFStringUTF16{}, nil, fmt.Errorf("cef_string_utf16_set failed")
	}
	cleanup := func() {
		stringClear(&out)
		runtime.KeepAlive(encoded)
	}
	return out, cleanup, nil
}

func ToGo(value *CEFStringUTF16) string {
	if value == nil || value.Str == nil || value.Length == 0 {
		return ""
	}
	slice := unsafe.Slice(value.Str, value.Length)
	return string(utf16.Decode(slice))
}

func clearBindingsForTest() {
	stringSet = nil
	stringClear = nil
	bound.Store(false)
}
