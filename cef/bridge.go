// bridge.go provides package-level helper functions that delegate to the
// core Engine. Generated code in cef/ calls these helpers for string
// conversion, refcount management, and pointer extraction.
//
// This file is handwritten — the rest of cef/ is generated.
package cef

import (
	"sync"
	"unsafe"

	"github.com/bnema/purego"
	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/core"
	portin "github.com/bnema/purego-cef/internal/ports/in"
)

// RawLifeSpanHandler is the low-level/raw lifespan handler interface.
// End users should usually implement LifeSpanHandler and pass it to
// NewLifeSpanHandler instead of implementing this directly.
type RawLifeSpanHandler = portin.RawLifeSpanHandler

// RawAudioHandler is the low-level/raw audio handler interface.
// End users should usually implement AudioHandler and pass it to
// NewAudioHandler instead of using this directly.
type RawAudioHandler = portin.RawAudioHandler

// LifeSpanHandler is the user-facing lifespan handler interface with typed
// out-params.
//
// Note: client out-params may arrive as nil when the underlying raw client
// handle originated inside CEF and cannot be meaningfully wrapped back into Go.
type LifeSpanHandler interface {
	OnBeforePopup(browser Browser, frame Frame, popupID int32, targetURL string, targetFrameName string,
		targetDisposition WindowOpenDisposition, userGesture int32, popupFeatures *PopupFeatures,
		windowInfo *WindowInfo, client *RawClient, settings *BrowserSettings,
		extraInfo *DictionaryValue, noJavascriptAccess *bool) bool
	OnBeforePopupAborted(browser Browser, popupID int32)
	OnBeforeDevToolsPopup(browser Browser, windowInfo *WindowInfo, client *RawClient,
		settings *BrowserSettings, extraInfo *DictionaryValue, useDefaultWindow *bool)
	OnAfterCreated(browser Browser)
	DoClose(browser Browser) bool
	OnBeforeClose(browser Browser)
}

// AudioHandler handles audio events with a safe [][]float32 data signature.
// The constructor decodes the raw float** callback input for you.
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

type rawLifeSpanHandlerWrapper struct {
	RawLifeSpanHandler
	rawPtr *capi.CEFLifeSpanHandlerT
}

func (w *rawLifeSpanHandlerWrapper) RawPointer() unsafe.Pointer {
	return unsafe.Pointer(w.rawPtr)
}

func NewRawLifeSpanHandler(impl RawLifeSpanHandler) RawLifeSpanHandler {
	r := new(capi.CEFLifeSpanHandlerT)
	w := &rawLifeSpanHandlerWrapper{rawPtr: r}
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
			unsafe.Pointer(arg9), settings, unsafe.Pointer(arg11), (*int32)(unsafe.Pointer(arg12)))

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
			(*BrowserSettings)(unsafe.Pointer(arg3)), unsafe.Pointer(arg4), (*int32)(unsafe.Pointer(arg5)))
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

	w.RawLifeSpanHandler = impl
	return w
}

func wrapLifeSpanHandler(_ unsafe.Pointer) RawLifeSpanHandler { return nil }

// ---------------------------------------------------------------------------
// Safe lifespan handler adapter
// ---------------------------------------------------------------------------

type safeLifeSpanHandlerWrapper struct {
	impl   LifeSpanHandler
	rawPtr *capi.CEFLifeSpanHandlerT
}

func (w *safeLifeSpanHandlerWrapper) RawPointer() unsafe.Pointer {
	return unsafe.Pointer(w.rawPtr)
}

// portin.LifeSpanHandler interface compliance — callbacks go through purego,
// so these methods should never be called directly. Panic to catch misuse.
func (w *safeLifeSpanHandlerWrapper) OnBeforePopup(Browser, Frame, int32, string, string, WindowOpenDisposition, int32, *PopupFeatures, *WindowInfo, unsafe.Pointer, *BrowserSettings, unsafe.Pointer, *int32) bool {
	panic("safeLifeSpanHandlerWrapper: raw OnBeforePopup called directly; callbacks go through purego")
}
func (w *safeLifeSpanHandlerWrapper) OnBeforePopupAborted(Browser, int32) {
	panic("safeLifeSpanHandlerWrapper: raw OnBeforePopupAborted called directly; callbacks go through purego")
}
func (w *safeLifeSpanHandlerWrapper) OnBeforeDevToolsPopup(Browser, *WindowInfo, unsafe.Pointer, *BrowserSettings, unsafe.Pointer, *int32) {
	panic("safeLifeSpanHandlerWrapper: raw OnBeforeDevToolsPopup called directly; callbacks go through purego")
}
func (w *safeLifeSpanHandlerWrapper) OnAfterCreated(Browser) {
	panic("safeLifeSpanHandlerWrapper: raw OnAfterCreated called directly; callbacks go through purego")
}
func (w *safeLifeSpanHandlerWrapper) DoClose(Browser) bool {
	panic("safeLifeSpanHandlerWrapper: raw DoClose called directly; callbacks go through purego")
}
func (w *safeLifeSpanHandlerWrapper) OnBeforeClose(Browser) {
	panic("safeLifeSpanHandlerWrapper: raw OnBeforeClose called directly; callbacks go through purego")
}

// NewLifeSpanHandler creates a raw lifespan handler from the user-facing
// typed lifespan handler interface. It converts raw callback out-params to
// typed Go values and writes back any changes the consumer makes.
func NewLifeSpanHandler(impl LifeSpanHandler) RawLifeSpanHandler {
	r := new(capi.CEFLifeSpanHandlerT)
	w := &safeLifeSpanHandlerWrapper{rawPtr: r, impl: impl}
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

		// Decode out-param: client (cef_client_t**)
		var clientVal RawClient
		if arg9 != 0 {
			if cp := *(*unsafe.Pointer)(unsafe.Pointer(arg9)); cp != nil {
				clientVal = wrapRawClient(cp)
			}
		}

		// Decode out-param: extraInfo (cef_dictionary_value_t**)
		var extraInfoVal DictionaryValue
		if arg11 != 0 {
			if ep := *(*unsafe.Pointer)(unsafe.Pointer(arg11)); ep != nil {
				extraInfoVal = wrapDictionaryValue(ep)
			}
		}

		// Decode out-param: noJavascriptAccess (int*)
		var noJS bool
		if arg12 != 0 {
			noJS = *(*int32)(unsafe.Pointer(arg12)) != 0
		}

		blocked := impl.OnBeforePopup(browser, frame, popupID, targetURL, targetFrameName,
			targetDisposition, userGesture, popupFeatures, windowInfo,
			&clientVal, settings, &extraInfoVal, &noJS)

		// Write back out-params
		if arg9 != 0 && clientVal != nil {
			if rp := extractRawPointer(clientVal); rp != nil {
				*(*unsafe.Pointer)(unsafe.Pointer(arg9)) = rp
			}
		}
		if arg11 != 0 && extraInfoVal != nil {
			if rp := extractRawPointer(extraInfoVal); rp != nil {
				*(*unsafe.Pointer)(unsafe.Pointer(arg11)) = rp
			}
		}
		if arg12 != 0 {
			if noJS {
				*(*int32)(unsafe.Pointer(arg12)) = 1
			} else {
				*(*int32)(unsafe.Pointer(arg12)) = 0
			}
		}

		if blocked {
			return 1
		}
		return 0
	}))

	r.OverrideOnBeforePopupAborted(purego.NewCallback(func(self uintptr, arg0, arg1 uintptr) {
		impl.OnBeforePopupAborted(wrapBrowser(unsafe.Pointer(arg0)), int32(arg1))
	}))

	r.OverrideOnBeforeDevToolsPopup(purego.NewCallback(func(self uintptr, arg0, arg1, arg2, arg3, arg4, arg5 uintptr) {
		browser := wrapBrowser(unsafe.Pointer(arg0))
		windowInfo := (*WindowInfo)(unsafe.Pointer(arg1))
		settings := (*BrowserSettings)(unsafe.Pointer(arg3))

		// Decode out-param: client (cef_client_t**)
		var clientVal RawClient
		if arg2 != 0 {
			if cp := *(*unsafe.Pointer)(unsafe.Pointer(arg2)); cp != nil {
				clientVal = wrapRawClient(cp)
			}
		}

		// Decode out-param: extraInfo (cef_dictionary_value_t**)
		var extraInfoVal DictionaryValue
		if arg4 != 0 {
			if ep := *(*unsafe.Pointer)(unsafe.Pointer(arg4)); ep != nil {
				extraInfoVal = wrapDictionaryValue(ep)
			}
		}

		// Decode out-param: useDefaultWindow (int*)
		var useDefault bool
		if arg5 != 0 {
			useDefault = *(*int32)(unsafe.Pointer(arg5)) != 0
		}

		impl.OnBeforeDevToolsPopup(browser, windowInfo, &clientVal, settings, &extraInfoVal, &useDefault)

		// Write back out-params
		if arg2 != 0 && clientVal != nil {
			if rp := extractRawPointer(clientVal); rp != nil {
				*(*unsafe.Pointer)(unsafe.Pointer(arg2)) = rp
			}
		}
		if arg4 != 0 && extraInfoVal != nil {
			if rp := extractRawPointer(extraInfoVal); rp != nil {
				*(*unsafe.Pointer)(unsafe.Pointer(arg4)) = rp
			}
		}
		if arg5 != 0 {
			if useDefault {
				*(*int32)(unsafe.Pointer(arg5)) = 1
			} else {
				*(*int32)(unsafe.Pointer(arg5)) = 0
			}
		}
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

	return w
}

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

// NewAudioHandler creates a CEF handler with decoded [][]float32 audio packets.
func NewAudioHandler(impl AudioHandler) RawAudioHandler {
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

// newRawAudioHandler creates a CEF handler from the low-level raw audio handler
// interface. Most users should prefer NewAudioHandler.
func newRawAudioHandler(impl RawAudioHandler) RawAudioHandler {
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
	impl   RawAudioHandler
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

func wrapAudioHandler(_ unsafe.Pointer) RawAudioHandler { return nil }

// NewRawAudioHandler exposes the low-level raw audio handler constructor for
// advanced callers. Most users should prefer NewAudioHandler.
func NewRawAudioHandler(impl RawAudioHandler) RawAudioHandler {
	return newRawAudioHandler(impl)
}

// ---------------------------------------------------------------------------
// Engine access — initialised once via sync.Once in Init()
// ---------------------------------------------------------------------------

var (
	eng      *core.Engine
	initOnce sync.Once
	initErr  error
)

func mustEng() *core.Engine {
	if eng == nil {
		panic("cef: engine not initialized; call cef.Init() first")
	}
	return eng
}

// cefString converts a Go string to a CEF UTF-16 string.
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

// ---------------------------------------------------------------------------
// User-facing Client with safe handler types
// ---------------------------------------------------------------------------

// Client is the user-facing client interface. It differs from RawClient in two ways:
//   - GetAudioHandler() returns cef.AudioHandler (decoded [][]float32) instead of RawAudioHandler
//   - GetLifeSpanHandler() returns cef.LifeSpanHandler (typed out-params) instead of RawLifeSpanHandler
//
// Use NewClient to create a raw CEF client from a Client implementation.
type Client interface {
	GetAudioHandler() AudioHandler
	GetCommandHandler() CommandHandler
	GetContextMenuHandler() ContextMenuHandler
	GetDialogHandler() DialogHandler
	GetDisplayHandler() DisplayHandler
	GetDownloadHandler() DownloadHandler
	GetDragHandler() DragHandler
	GetFindHandler() FindHandler
	GetFocusHandler() FocusHandler
	GetFrameHandler() FrameHandler
	GetPermissionHandler() PermissionHandler
	GetJsdialogHandler() JsdialogHandler
	GetKeyboardHandler() KeyboardHandler
	GetLifeSpanHandler() LifeSpanHandler
	GetLoadHandler() LoadHandler
	GetPrintHandler() PrintHandler
	GetRenderHandler() RenderHandler
	GetRequestHandler() RequestHandler
	OnProcessMessageReceived(browser Browser, frame Frame, sourceProcess ProcessID, message ProcessMessage) int32
}

// clientAdapter wraps a Client to satisfy RawClient by converting safe handler
// types to their raw equivalents.
type clientAdapter struct {
	impl Client
}

func (a *clientAdapter) GetAudioHandler() RawAudioHandler {
	h := a.impl.GetAudioHandler()
	if h == nil {
		return nil
	}
	return NewAudioHandler(h)
}

func (a *clientAdapter) GetLifeSpanHandler() RawLifeSpanHandler {
	h := a.impl.GetLifeSpanHandler()
	if h == nil {
		return nil
	}
	return NewLifeSpanHandler(h)
}

func (a *clientAdapter) GetCommandHandler() CommandHandler {
	return a.impl.GetCommandHandler()
}
func (a *clientAdapter) GetContextMenuHandler() ContextMenuHandler {
	return a.impl.GetContextMenuHandler()
}
func (a *clientAdapter) GetDialogHandler() DialogHandler {
	return a.impl.GetDialogHandler()
}
func (a *clientAdapter) GetDisplayHandler() DisplayHandler {
	return a.impl.GetDisplayHandler()
}
func (a *clientAdapter) GetDownloadHandler() DownloadHandler {
	return a.impl.GetDownloadHandler()
}
func (a *clientAdapter) GetDragHandler() DragHandler {
	return a.impl.GetDragHandler()
}
func (a *clientAdapter) GetFindHandler() FindHandler {
	return a.impl.GetFindHandler()
}
func (a *clientAdapter) GetFocusHandler() FocusHandler {
	return a.impl.GetFocusHandler()
}
func (a *clientAdapter) GetFrameHandler() FrameHandler {
	return a.impl.GetFrameHandler()
}
func (a *clientAdapter) GetPermissionHandler() PermissionHandler {
	return a.impl.GetPermissionHandler()
}
func (a *clientAdapter) GetJsdialogHandler() JsdialogHandler {
	return a.impl.GetJsdialogHandler()
}
func (a *clientAdapter) GetKeyboardHandler() KeyboardHandler {
	return a.impl.GetKeyboardHandler()
}
func (a *clientAdapter) GetLoadHandler() LoadHandler {
	return a.impl.GetLoadHandler()
}
func (a *clientAdapter) GetPrintHandler() PrintHandler {
	return a.impl.GetPrintHandler()
}
func (a *clientAdapter) GetRenderHandler() RenderHandler {
	return a.impl.GetRenderHandler()
}
func (a *clientAdapter) GetRequestHandler() RequestHandler {
	return a.impl.GetRequestHandler()
}
func (a *clientAdapter) OnProcessMessageReceived(browser Browser, frame Frame, sourceProcess ProcessID, message ProcessMessage) int32 {
	return a.impl.OnProcessMessageReceived(browser, frame, sourceProcess, message)
}

// NewClient creates a raw CEF client from a user-facing Client implementation.
// It wraps AudioHandler via NewAudioHandler ([][]float32 decoding) and
// LifeSpanHandler via NewLifeSpanHandler (typed out-params), then delegates
// to the generated raw client constructor.
func NewClient(impl Client) RawClient {
	return NewRawClient(&clientAdapter{impl: impl})
}
