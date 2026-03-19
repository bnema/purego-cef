//go:build amd64 && linux

package cef

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

// defaultAPIVersion is the CEF API version we target. This must match
// CEF_API_VERSION_LAST from the headers used for generation.
const defaultAPIVersion = 14500

// setAPIVersion patches the internal g_version global inside libcef.so so
// that cef_api_version() returns our target version instead of -1. Without
// this, CEF's versioned CToCpp wrappers reject all client structs.
//
// The cef_api_version function is a trivial global read on x86-64:
//
//	push %rbp; mov %rsp,%rbp; mov <offset>(%rip),%eax; pop %rbp; ret
//
// We validate the opcode and decode the RIP-relative offset to locate
// g_version, then write our target version.
func setAPIVersion(handle uintptr) error {
	sym, err := purego.Dlsym(handle, "cef_api_version")
	if err != nil {
		return fmt.Errorf("resolve cef_api_version: %w", err)
	}

	// Validate the opcode at sym+4 is "mov offset(%rip),%eax" (0x8b 0x05).
	op1 := *(*byte)(unsafe.Pointer(sym + 4))
	op2 := *(*byte)(unsafe.Pointer(sym + 5))
	if op1 != 0x8b || op2 != 0x05 {
		return fmt.Errorf("unexpected cef_api_version encoding: got %02x %02x, want 8b 05", op1, op2)
	}

	// Decode the 4-byte signed RIP-relative offset at sym+6.
	// RIP points to the next instruction at sym+10.
	offset := *(*int32)(unsafe.Pointer(sym + 6))
	gVersionAddr := sym + 10 + uintptr(offset)

	*(*int32)(unsafe.Pointer(gVersionAddr)) = defaultAPIVersion
	return nil
}
