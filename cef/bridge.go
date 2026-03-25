// bridge.go provides package-level helper functions that delegate to the
// core Engine. Generated code in cef/ calls these helpers for string
// conversion, refcount management, and pointer extraction.
//
// This file is handwritten — the rest of cef/ is generated.
package cef

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/core"
	portin "github.com/bnema/purego-cef/internal/ports/in"
	"github.com/ebitengine/purego"
)

// Type alias for LifeSpanHandler — uses generated interface (out-params are unsafe.Pointer).
type LifeSpanHandler = portin.LifeSpanHandler

// AudioHandler handles audio events with a safe [][]float32 data signature.
// This is NOT an alias of the generated portin.AudioHandler (which uses unsafe.Pointer
// for the data param). The constructor decodes the raw float** to [][]float32.
type AudioHandler interface {
	GetAudioParameters(browser Browser, params *AudioParameters) int32
	OnAudioStreamStarted(browser Browser, params *AudioParameters, channels int32)
	OnAudioStreamPacket(browser Browser, data [][]float32, frames int32, pts int64)
	OnAudioStreamStopped(browser Browser)
	OnAudioStreamError(browser Browser, message string)
}

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
// AudioHandler constructor — decodes float** to [][]float32
// ---------------------------------------------------------------------------

type audioHandlerWrapper struct {
	impl     AudioHandler
	rawPtr   *capi.CEFAudioHandlerT
	mu       sync.Mutex
	channels int32
}

func (w *audioHandlerWrapper) RawPointer() unsafe.Pointer {
	return unsafe.Pointer(w.rawPtr)
}

// Satisfy the generated portin.AudioHandler interface for extractRawPointer.
func (w *audioHandlerWrapper) GetAudioParameters(browser Browser, params *AudioParameters) int32 {
	return w.impl.GetAudioParameters(browser, params)
}
func (w *audioHandlerWrapper) OnAudioStreamStarted(browser Browser, params *AudioParameters, channels int32) {
	w.impl.OnAudioStreamStarted(browser, params, channels)
}
func (w *audioHandlerWrapper) OnAudioStreamPacket(browser Browser, _ unsafe.Pointer, frames int32, pts int64) {
	// The actual decoded data is passed through the safe interface by the callback below.
	// This method exists only for interface compliance; it is never called directly.
}
func (w *audioHandlerWrapper) OnAudioStreamStopped(browser Browser) {
	w.impl.OnAudioStreamStopped(browser)
}
func (w *audioHandlerWrapper) OnAudioStreamError(browser Browser, message string) {
	w.impl.OnAudioStreamError(browser, message)
}

// NewSafeAudioHandler creates a CEF handler with safe [][]float32 audio data decoding.
// Use this instead of NewAudioHandler when you want decoded audio packets.
func NewSafeAudioHandler(impl AudioHandler) portin.AudioHandler {
	r := new(capi.CEFAudioHandlerT)
	w := &audioHandlerWrapper{rawPtr: r, impl: impl}
	initRefCount(unsafe.Pointer(r), unsafe.Sizeof(*r), w)

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
		decoded := core.DecodeAudioPacket(unsafe.Pointer(arg1), ch, frames)
		impl.OnAudioStreamPacket(browser, decoded, frames, pts)
	}))

	r.OverrideOnAudioStreamStopped(purego.NewCallback(func(self uintptr, arg0 uintptr) {
		impl.OnAudioStreamStopped(wrapBrowser(unsafe.Pointer(arg0)))
	}))

	r.OverrideOnAudioStreamError(purego.NewCallback(func(self uintptr, arg0, arg1 uintptr) {
		impl.OnAudioStreamError(wrapBrowser(unsafe.Pointer(arg0)), goString(unsafe.Pointer(arg1)))
	}))

	return w
}

// NewAudioHandler creates a CEF handler from the generated portin.AudioHandler interface.
// For decoded [][]float32 audio data, use NewSafeAudioHandler instead.
func NewAudioHandler(impl portin.AudioHandler) portin.AudioHandler {
	r := new(capi.CEFAudioHandlerT)
	w := &rawAudioHandlerWrapper{impl: impl, rawPtr: r}
	initRefCount(unsafe.Pointer(r), unsafe.Sizeof(*r), w)

	r.OverrideGetAudioParameters(purego.NewCallback(func(self uintptr, arg0, arg1 uintptr) uintptr {
		return uintptr(impl.GetAudioParameters(wrapBrowser(unsafe.Pointer(arg0)), (*AudioParameters)(unsafe.Pointer(arg1))))
	}))
	r.OverrideOnAudioStreamStarted(purego.NewCallback(func(self uintptr, arg0, arg1, arg2 uintptr) {
		impl.OnAudioStreamStarted(wrapBrowser(unsafe.Pointer(arg0)), (*AudioParameters)(unsafe.Pointer(arg1)), int32(arg2))
	}))
	r.OverrideOnAudioStreamPacket(purego.NewCallback(func(self uintptr, arg0, arg1, arg2, arg3 uintptr) {
		impl.OnAudioStreamPacket(wrapBrowser(unsafe.Pointer(arg0)), unsafe.Pointer(arg1), int32(arg2), int64(arg3))
	}))
	r.OverrideOnAudioStreamStopped(purego.NewCallback(func(self uintptr, arg0 uintptr) {
		impl.OnAudioStreamStopped(wrapBrowser(unsafe.Pointer(arg0)))
	}))
	r.OverrideOnAudioStreamError(purego.NewCallback(func(self uintptr, arg0, arg1 uintptr) {
		impl.OnAudioStreamError(wrapBrowser(unsafe.Pointer(arg0)), goString(unsafe.Pointer(arg1)))
	}))

	return w
}

type rawAudioHandlerWrapper struct {
	impl   portin.AudioHandler
	rawPtr *capi.CEFAudioHandlerT
}

func (w *rawAudioHandlerWrapper) RawPointer() unsafe.Pointer { return unsafe.Pointer(w.rawPtr) }
func (w *rawAudioHandlerWrapper) GetAudioParameters(b Browser, p *AudioParameters) int32 {
	return w.impl.GetAudioParameters(b, p)
}
func (w *rawAudioHandlerWrapper) OnAudioStreamStarted(b Browser, p *AudioParameters, c int32) {
	w.impl.OnAudioStreamStarted(b, p, c)
}
func (w *rawAudioHandlerWrapper) OnAudioStreamPacket(b Browser, d unsafe.Pointer, f int32, p int64) {
	w.impl.OnAudioStreamPacket(b, d, f, p)
}
func (w *rawAudioHandlerWrapper) OnAudioStreamStopped(b Browser) { w.impl.OnAudioStreamStopped(b) }
func (w *rawAudioHandlerWrapper) OnAudioStreamError(b Browser, m string) {
	w.impl.OnAudioStreamError(b, m)
}

func wrapAudioHandler(_ unsafe.Pointer) portin.AudioHandler { return nil }

// ---------------------------------------------------------------------------
// Engine access — atomic for thread safety (#8)
// ---------------------------------------------------------------------------

// engPtr is the core Engine instance, wired at Init() time.
// Accessed atomically so bridge helpers are safe from any goroutine.
var engPtr atomic.Pointer[core.Engine]

func mustEng() *core.Engine {
	e := engPtr.Load()
	if e == nil {
		panic("cef: engine not initialized; call cef.Init() first")
	}
	return e
}

// cefString converts a Go string to a CEF UTF-16 string.
// CEFStringT is layout-identical in core and capi packages (both mirror cef_string_t).
func cefString(s string) core.CEFStringT {
	return mustEng().CefString(s)
}

// freeCefString releases a CEF string's backing memory.
func freeCefString(cs *core.CEFStringT) {
	mustEng().FreeCefString(cs)
}

// goString converts a pointer to a CEF string to a Go string.
func goString(cs unsafe.Pointer) string {
	return core.GoString(cs)
}

// goStringUserfree converts a cef_string_userfree_t to a Go string and frees it.
func goStringUserfree(ptr unsafe.Pointer) string {
	return mustEng().GoStringUserfree(ptr)
}

// initRefCount wires refcount callbacks into a CEF base struct header.
func initRefCount(base unsafe.Pointer, size uintptr, owner any) {
	mustEng().Refs().InitRefCount(base, size, owner)
}

// addRef increments the refcount for the object at base.
func addRef(base unsafe.Pointer) {
	mustEng().Refs().AddRef(base)
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
