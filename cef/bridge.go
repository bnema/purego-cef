// bridge.go provides package-level helper functions that delegate to the
// core Engine. Generated code in cef/ calls these helpers for string
// conversion, refcount management, and pointer extraction.
//
// This file is handwritten — the rest of cef/ is generated.
package cef

import (
	"unsafe"

	"github.com/bnema/purego-cef/internal/core"
	in "github.com/bnema/purego-cef/internal/ports/in"
)

// Type aliases for interfaces that are excluded from generation (skipPublicTypes)
// because their handler constructors need handwritten implementations.
type AudioHandler = in.AudioHandler
type LifeSpanHandler = in.LifeSpanHandler

// Stub constructors for skipped handler types.
// TODO: implement proper handler constructors with out-param handling.
func NewAudioHandler(_ AudioHandler) AudioHandler          { return nil }
func NewLifeSpanHandler(_ LifeSpanHandler) LifeSpanHandler { return nil }

// eng is the core Engine instance, wired at Init() time.
var eng *core.Engine

// cefString converts a Go string to a CEF UTF-16 string.
func cefString(s string) core.CEFStringT {
	return eng.CefString(s)
}

// freeCefString releases a CEF string's backing memory.
func freeCefString(cs *core.CEFStringT) {
	eng.FreeCefString(cs)
}

// goString converts a pointer to a CEF string to a Go string.
func goString(cs unsafe.Pointer) string {
	return core.GoString(cs)
}

// goStringUserfree converts a cef_string_userfree_t to a Go string and frees it.
func goStringUserfree(ptr unsafe.Pointer) string {
	return eng.GoStringUserfree(ptr)
}

// initRefCount wires refcount callbacks into a CEF base struct header.
func initRefCount(base unsafe.Pointer, size uintptr, owner any) {
	eng.Refs().InitRefCount(base, size, owner)
}

// addRef increments the refcount for the object at base.
func addRef(base unsafe.Pointer) {
	eng.Refs().AddRef(base)
}

// extractRawPointer returns the underlying raw CEF pointer from an interface.
func extractRawPointer(v any) unsafe.Pointer {
	return core.ExtractRawPointer(v)
}

// extractOrWrapRawPointer returns the raw pointer for v, calling wrap if needed.
func extractOrWrapRawPointer(v any, wrap func() any) unsafe.Pointer {
	return core.ExtractOrWrapRawPointer(v, wrap)
}

// decodeSlice converts a raw pointer and count into a Go slice of T.
func decodeSlice[T any](ptr uintptr, count int) []T {
	return core.DecodeSlice[T](ptr, count)
}
