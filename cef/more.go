package cef

import (
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"

	"github.com/bnema/purego-cef/cef/internal/raw"
)

// NewKeyEvent creates a KeyEvent with Size pre-filled and the given parameters set.
func NewKeyEvent(eventType KeyEventType, windowsKeyCode, nativeKeyCode int32, modifiers uint32) KeyEvent {
	var evt KeyEvent
	evt.Size = unsafe.Sizeof(evt)
	*(*int32)(unsafe.Pointer(&evt.Type)) = int32(eventType)
	evt.WindowsKeyCode = windowsKeyCode
	evt.NativeKeyCode = nativeKeyCode
	evt.Modifiers = modifiers
	return evt
}

// KeyEventSetType sets the key event type using an unsafe cast from the public
// KeyEventType to the internal CEFKeyEventTypeT.
func KeyEventSetType(e *KeyEvent, t KeyEventType) {
	*(*int32)(unsafe.Pointer(&e.Type)) = int32(t)
}

// KeyEventSetModifiers sets the modifier flags on the key event.
func KeyEventSetModifiers(e *KeyEvent, m uint32) {
	e.Modifiers = m
}

// KeyEventSetCharacter sets the character value on the key event.
func KeyEventSetCharacter(e *KeyEvent, c uint16) {
	e.Character = c
}

// KeyEventSetUnmodifiedCharacter sets the unmodified character value on the key event.
func KeyEventSetUnmodifiedCharacter(e *KeyEvent, c uint16) {
	e.UnmodifiedCharacter = c
}

// KeyEventSetIsSystemKey sets whether this is a system key event.
func KeyEventSetIsSystemKey(e *KeyEvent, v bool) {
	if v {
		e.IsSystemKey = 1
	} else {
		e.IsSystemKey = 0
	}
}

// KeyEventSetFocusOnEditableField sets whether focus is on an editable field.
func KeyEventSetFocusOnEditableField(e *KeyEvent, v bool) {
	if v {
		e.FocusOnEditableField = 1
	} else {
		e.FocusOnEditableField = 0
	}
}

// NewScreenInfo creates a ScreenInfo with Size pre-filled and the given parameters set.
func NewScreenInfo(deviceScaleFactor float32, depth, depthPerComponent int32,
	isMonochrome bool, rect, availableRect Rect) ScreenInfo {
	var si ScreenInfo
	si.Size = unsafe.Sizeof(si)
	si.DeviceScaleFactor = deviceScaleFactor
	si.Depth = depth
	si.DepthPerComponent = depthPerComponent
	if isMonochrome {
		si.IsMonochrome = 1
	}
	si.Rect = rect
	si.AvailableRect = availableRect
	return si
}

// DecodeAudioPacket converts a CEF audio data pointer (float** layout) into a
// Go [][]float32 slice indexed by [channel][frame].
func DecodeAudioPacket(data unsafe.Pointer, channels, frames int32) [][]float32 {
	if data == nil || channels <= 0 || frames <= 0 {
		return nil
	}
	channelPtrs := unsafe.Slice((**float32)(data), int(channels))
	result := make([][]float32, int(channels))
	for i, cp := range channelPtrs {
		if cp != nil {
			result[i] = unsafe.Slice(cp, int(frames))
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// LifeSpanHandler — hand-written with safe signatures (skipped by cefgen)
// ---------------------------------------------------------------------------

// LifeSpanHandler handles events related to browser life span.
type LifeSpanHandler interface {
	OnBeforePopup(browser Browser, frame Frame, popupID int32, targetURL string,
		targetFrameName string, targetDisposition WindowOpenDisposition,
		userGesture int32, popupFeatures *PopupFeatures, windowInfo *WindowInfo,
		client *Client, settings *BrowserSettings, extraInfo *DictionaryValue,
		noJavascriptAccess *bool) bool
	OnBeforePopupAborted(browser Browser, popupID int32)
	OnBeforeDevToolsPopup(browser Browser, windowInfo *WindowInfo,
		client *Client, settings *BrowserSettings, extraInfo *DictionaryValue,
		useDefaultWindow *bool)
	OnAfterCreated(browser Browser)
	DoClose(browser Browser) bool
	OnBeforeClose(browser Browser)
}

type lifeSpanHandlerWrapper struct {
	LifeSpanHandler
	rawPtr *raw.CEFLifeSpanHandlerT
}

func (w *lifeSpanHandlerWrapper) rawPointer() unsafe.Pointer {
	return unsafe.Pointer(w.rawPtr)
}

// NewLifeSpanHandler creates a CEF handler backed by the given implementation.
func NewLifeSpanHandler(impl LifeSpanHandler) LifeSpanHandler {
	r := new(raw.CEFLifeSpanHandlerT)
	initRefCount(unsafe.Pointer(r), unsafe.Sizeof(*r), r)

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

		// Typed out-params for consumer.
		var clientVal Client
		var extraInfoVal DictionaryValue
		var noJSVal bool

		blocked := impl.OnBeforePopup(browser, frame, popupID, targetURL, targetFrameName,
			targetDisposition, userGesture, popupFeatures, windowInfo,
			&clientVal, settings, &extraInfoVal, &noJSVal)

		if !blocked {
			// Write back out-params to C double pointers.
			if clientVal != nil && arg9 != 0 {
				*(*uintptr)(unsafe.Pointer(arg9)) = uintptr(extractOrWrapRawPointer(clientVal, func() any {
					return NewClient(clientVal)
				}))
			}
			if extraInfoVal != nil && arg11 != 0 {
				*(*uintptr)(unsafe.Pointer(arg11)) = uintptr(extractRawPointer(extraInfoVal))
			}
			if noJSVal && arg12 != 0 {
				*(*int32)(unsafe.Pointer(arg12)) = 1
			}
		}

		if blocked {
			return 1
		}
		return 0
	}))

	r.OverrideOnBeforePopupAborted(purego.NewCallback(func(self uintptr, arg0, arg1 uintptr) {
		browser := wrapBrowser(unsafe.Pointer(arg0))
		popupID := int32(arg1)
		impl.OnBeforePopupAborted(browser, popupID)
	}))

	r.OverrideOnBeforeDevToolsPopup(purego.NewCallback(func(self uintptr, arg0, arg1, arg2, arg3, arg4, arg5 uintptr) {
		browser := wrapBrowser(unsafe.Pointer(arg0))
		windowInfo := (*WindowInfo)(unsafe.Pointer(arg1))
		settings := (*BrowserSettings)(unsafe.Pointer(arg3))

		var clientVal Client
		var extraInfoVal DictionaryValue
		var useDefaultVal bool

		impl.OnBeforeDevToolsPopup(browser, windowInfo, &clientVal, settings, &extraInfoVal, &useDefaultVal)

		if clientVal != nil && arg2 != 0 {
			*(*uintptr)(unsafe.Pointer(arg2)) = uintptr(extractOrWrapRawPointer(clientVal, func() any {
				return NewClient(clientVal)
			}))
		}
		if extraInfoVal != nil && arg4 != 0 {
			*(*uintptr)(unsafe.Pointer(arg4)) = uintptr(extractRawPointer(extraInfoVal))
		}
		if useDefaultVal && arg5 != 0 {
			*(*int32)(unsafe.Pointer(arg5)) = 1
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

	w := &lifeSpanHandlerWrapper{rawPtr: r}
	w.LifeSpanHandler = impl
	return w
}

func wrapLifeSpanHandler(ptr unsafe.Pointer) LifeSpanHandler {
	return nil
}

// ---------------------------------------------------------------------------
// AudioHandler — hand-written with safe signatures (skipped by cefgen)
// ---------------------------------------------------------------------------

// AudioHandler handles audio events. OnAudioStreamPacket receives decoded
// [][]float32 data instead of an unsafe.Pointer.
type AudioHandler interface {
	GetAudioParameters(browser Browser, params *AudioParameters) int32
	OnAudioStreamStarted(browser Browser, params *AudioParameters, channels int32)
	OnAudioStreamPacket(browser Browser, data [][]float32, frames int32, pts int64)
	OnAudioStreamStopped(browser Browser)
	OnAudioStreamError(browser Browser, message string)
}

type audioHandlerWrapper struct {
	AudioHandler
	rawPtr   *raw.CEFAudioHandlerT
	mu       sync.Mutex
	channels int32
}

func (w *audioHandlerWrapper) rawPointer() unsafe.Pointer {
	return unsafe.Pointer(w.rawPtr)
}

// NewAudioHandler creates a CEF handler backed by the given implementation.
// It tracks channel count from OnAudioStreamStarted and decodes the raw
// audio data to [][]float32 in OnAudioStreamPacket.
func NewAudioHandler(impl AudioHandler) AudioHandler {
	r := new(raw.CEFAudioHandlerT)
	initRefCount(unsafe.Pointer(r), unsafe.Sizeof(*r), r)

	w := &audioHandlerWrapper{rawPtr: r}
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
		decoded := DecodeAudioPacket(unsafe.Pointer(arg1), ch, frames)
		impl.OnAudioStreamPacket(browser, decoded, frames, pts)
	}))

	r.OverrideOnAudioStreamStopped(purego.NewCallback(func(self uintptr, arg0 uintptr) {
		impl.OnAudioStreamStopped(wrapBrowser(unsafe.Pointer(arg0)))
	}))

	r.OverrideOnAudioStreamError(purego.NewCallback(func(self uintptr, arg0, arg1 uintptr) {
		browser := wrapBrowser(unsafe.Pointer(arg0))
		message := goString(unsafe.Pointer(arg1))
		impl.OnAudioStreamError(browser, message)
	}))

	return w
}

func wrapAudioHandler(ptr unsafe.Pointer) AudioHandler {
	return nil
}
