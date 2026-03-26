// bridge.go provides package-level helper functions that delegate to the
// core Engine. Generated code in cef/ calls these helpers for string
// conversion, refcount management, and pointer extraction.
//
// This file is handwritten — the rest of cef/ is generated.
package cef

import (
	"sync"
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/core"
	in "github.com/bnema/purego-cef/internal/ports/in"
	"github.com/ebitengine/purego"
)

// Type aliases for interfaces that are excluded from generation (skipPublicTypes)
// because their handler constructors need handwritten implementations.
type AudioHandler = in.AudioHandler
type LifeSpanHandler = in.LifeSpanHandler

// ---------------------------------------------------------------------------
// LifeSpanHandler constructor
// ---------------------------------------------------------------------------

type lifeSpanHandlerWrapper struct {
	LifeSpanHandler
	rawPtr *capi.CEFLifeSpanHandlerT
}

func (w *lifeSpanHandlerWrapper) RawPointer() unsafe.Pointer {
	return unsafe.Pointer(w.rawPtr)
}

func NewLifeSpanHandler(impl LifeSpanHandler) LifeSpanHandler {
	r := new(capi.CEFLifeSpanHandlerT)
	w := &lifeSpanHandlerWrapper{rawPtr: r}
	initRefCount(unsafe.Pointer(r), unsafe.Sizeof(*r), w)

	r.OverrideOnBeforePopup(purego.NewCallback(func(self uintptr, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11, arg12 uintptr) uintptr {
		browser := wrapBrowser(unsafe.Pointer(arg0))
		frame := wrapFrame(unsafe.Pointer(arg1))
		popupID := int32(arg2)
		targetURL := goString(unsafe.Pointer(arg3))
		targetFrameName := goString(unsafe.Pointer(arg4))
		targetDisposition := WindowOpenDisposition(arg5)
		userGesture := int32(arg6)
		popupFeatures := (*PopupFeatures)(unsafe.Pointer(arg7))
		windowInfo := (*WindowInfo)(unsafe.Pointer(arg8))
		settings := (*BrowserSettings)(unsafe.Pointer(arg10))

		blocked := impl.OnBeforePopup(browser, frame, popupID, targetURL, targetFrameName,
			targetDisposition, userGesture, popupFeatures, windowInfo,
			unsafe.Pointer(arg9), settings, unsafe.Pointer(arg11), unsafe.Pointer(arg12))

		if blocked {
			return 1
		}
		return 0
	}))

	r.OverrideOnBeforePopupAborted(purego.NewCallback(func(self uintptr, arg0, arg1 uintptr) {
		impl.OnBeforePopupAborted(wrapBrowser(unsafe.Pointer(arg0)), int32(arg1))
	}))

	r.OverrideOnBeforeDevToolsPopup(purego.NewCallback(func(self uintptr, arg0, arg1, arg2, arg3, arg4, arg5 uintptr) {
		impl.OnBeforeDevToolsPopup(wrapBrowser(unsafe.Pointer(arg0)),
			(*WindowInfo)(unsafe.Pointer(arg1)), unsafe.Pointer(arg2),
			(*BrowserSettings)(unsafe.Pointer(arg3)), unsafe.Pointer(arg4), unsafe.Pointer(arg5))
	}))

	r.OverrideOnAfterCreated(purego.NewCallback(func(self uintptr, arg0 uintptr) {
		impl.OnAfterCreated(wrapBrowser(unsafe.Pointer(arg0)))
	}))

	r.OverrideDoClose(purego.NewCallback(func(self uintptr, arg0 uintptr) uintptr {
		if impl.DoClose(wrapBrowser(unsafe.Pointer(arg0))) {
			return 1
		}
		return 0
	}))

	r.OverrideOnBeforeClose(purego.NewCallback(func(self uintptr, arg0 uintptr) {
		impl.OnBeforeClose(wrapBrowser(unsafe.Pointer(arg0)))
	}))

	w.LifeSpanHandler = impl
	return w
}

func wrapLifeSpanHandler(_ unsafe.Pointer) LifeSpanHandler { return nil }

// ---------------------------------------------------------------------------
// AudioHandler constructor
// ---------------------------------------------------------------------------

type audioHandlerWrapper struct {
	AudioHandler
	rawPtr   *capi.CEFAudioHandlerT
	mu       sync.Mutex
	channels int32
}

func (w *audioHandlerWrapper) RawPointer() unsafe.Pointer {
	return unsafe.Pointer(w.rawPtr)
}

func NewAudioHandler(impl AudioHandler) AudioHandler {
	r := new(capi.CEFAudioHandlerT)
	w := &audioHandlerWrapper{rawPtr: r}
	initRefCount(unsafe.Pointer(r), unsafe.Sizeof(*r), w)
	w.AudioHandler = impl

	r.OverrideGetAudioParameters(purego.NewCallback(func(self uintptr, arg0, arg1 uintptr) uintptr {
		browser := wrapBrowser(unsafe.Pointer(arg0))
		params := (*AudioParameters)(unsafe.Pointer(arg1))
		return uintptr(impl.GetAudioParameters(browser, params))
	}))

	r.OverrideOnAudioStreamStarted(purego.NewCallback(func(self uintptr, arg0, arg1, arg2 uintptr) {
		browser := wrapBrowser(unsafe.Pointer(arg0))
		params := (*AudioParameters)(unsafe.Pointer(arg1))
		channels := int32(arg2)
		w.mu.Lock()
		w.channels = channels
		w.mu.Unlock()
		impl.OnAudioStreamStarted(browser, params, channels)
	}))

	r.OverrideOnAudioStreamPacket(purego.NewCallback(func(self uintptr, arg0, arg1, arg2, arg3 uintptr) {
		browser := wrapBrowser(unsafe.Pointer(arg0))
		frames := int32(arg2)
		pts := int64(arg3)
		w.mu.Lock()
		ch := w.channels
		w.mu.Unlock()
		impl.OnAudioStreamPacket(browser, unsafe.Pointer(arg1), frames, pts)
		_ = ch // channel count available for future audio decoding helpers
	}))

	r.OverrideOnAudioStreamStopped(purego.NewCallback(func(self uintptr, arg0 uintptr) {
		impl.OnAudioStreamStopped(wrapBrowser(unsafe.Pointer(arg0)))
	}))

	r.OverrideOnAudioStreamError(purego.NewCallback(func(self uintptr, arg0, arg1 uintptr) {
		impl.OnAudioStreamError(wrapBrowser(unsafe.Pointer(arg0)), goString(unsafe.Pointer(arg1)))
	}))

	return w
}

func wrapAudioHandler(_ unsafe.Pointer) AudioHandler { return nil }

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
