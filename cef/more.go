package cef

import (
	"sync"
	"unsafe"
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
// SafeLifeSpanHandler — popup callbacks without unsafe.Pointer
// ---------------------------------------------------------------------------

// PopupAction is the return value from SafeLifeSpanHandler.OnBeforePopup.
// Set Block to true to prevent the popup from opening.
// When Block is false, the remaining fields configure the popup.
type PopupAction struct {
	Block              bool
	Client             Client          // nil keeps the default client
	ExtraInfo          DictionaryValue // nil keeps the default extra info
	NoJavascriptAccess bool
}

// DevToolsPopupAction is the return value from SafeLifeSpanHandler.OnBeforeDevToolsPopup.
type DevToolsPopupAction struct {
	Client           Client          // nil keeps the default client
	ExtraInfo        DictionaryValue // nil keeps the default extra info
	UseDefaultWindow bool
}

// SafeLifeSpanHandler is a consumer-friendly alternative to LifeSpanHandler
// that replaces unsafe.Pointer params with typed fields.
type SafeLifeSpanHandler interface {
	OnBeforePopup(browser Browser, frame Frame, popupID int32, targetURL string,
		targetFrameName string, targetDisposition WindowOpenDisposition,
		userGesture int32, popupFeatures *PopupFeatures, windowInfo *WindowInfo,
		settings *BrowserSettings) PopupAction
	OnBeforePopupAborted(browser Browser, popupID int32)
	OnBeforeDevToolsPopup(browser Browser, windowInfo *WindowInfo,
		settings *BrowserSettings) DevToolsPopupAction
	OnAfterCreated(browser Browser)
	DoClose(browser Browser) bool
	OnBeforeClose(browser Browser)
}

// safeLifeSpanAdapter adapts SafeLifeSpanHandler to LifeSpanHandler.
type safeLifeSpanAdapter struct {
	impl SafeLifeSpanHandler
}

func (a *safeLifeSpanAdapter) OnBeforePopup(browser Browser, frame Frame, popupID int32,
	targetURL string, targetFrameName string, targetDisposition WindowOpenDisposition,
	userGesture int32, popupFeatures *PopupFeatures, windowInfo *WindowInfo,
	client unsafe.Pointer, settings *BrowserSettings, extraInfo unsafe.Pointer,
	noJavascriptAccess *int32) bool {

	action := a.impl.OnBeforePopup(browser, frame, popupID, targetURL, targetFrameName,
		targetDisposition, userGesture, popupFeatures, windowInfo, settings)

	if action.Block {
		return true
	}
	if action.Client != nil && client != nil {
		*(*uintptr)(client) = uintptr(extractRawPointer(action.Client))
	}
	if action.ExtraInfo != nil && extraInfo != nil {
		*(*uintptr)(extraInfo) = uintptr(extractRawPointer(action.ExtraInfo))
	}
	if noJavascriptAccess != nil {
		if action.NoJavascriptAccess {
			*noJavascriptAccess = 1
		} else {
			*noJavascriptAccess = 0
		}
	}
	return false
}

func (a *safeLifeSpanAdapter) OnBeforePopupAborted(browser Browser, popupID int32) {
	a.impl.OnBeforePopupAborted(browser, popupID)
}

func (a *safeLifeSpanAdapter) OnBeforeDevToolsPopup(browser Browser,
	windowInfo *WindowInfo, client unsafe.Pointer, settings *BrowserSettings,
	extraInfo unsafe.Pointer, useDefaultWindow *int32) {

	action := a.impl.OnBeforeDevToolsPopup(browser, windowInfo, settings)

	if action.Client != nil && client != nil {
		*(*uintptr)(client) = uintptr(extractRawPointer(action.Client))
	}
	if action.ExtraInfo != nil && extraInfo != nil {
		*(*uintptr)(extraInfo) = uintptr(extractRawPointer(action.ExtraInfo))
	}
	if useDefaultWindow != nil {
		if action.UseDefaultWindow {
			*useDefaultWindow = 1
		} else {
			*useDefaultWindow = 0
		}
	}
}

func (a *safeLifeSpanAdapter) OnAfterCreated(browser Browser) {
	a.impl.OnAfterCreated(browser)
}

func (a *safeLifeSpanAdapter) DoClose(browser Browser) bool {
	return a.impl.DoClose(browser)
}

func (a *safeLifeSpanAdapter) OnBeforeClose(browser Browser) {
	a.impl.OnBeforeClose(browser)
}

// NewSafeLifeSpanHandler creates a CEF LifeSpanHandler from a SafeLifeSpanHandler.
func NewSafeLifeSpanHandler(impl SafeLifeSpanHandler) LifeSpanHandler {
	return NewLifeSpanHandler(&safeLifeSpanAdapter{impl: impl})
}

// ---------------------------------------------------------------------------
// SafeAudioHandler — audio callbacks without unsafe.Pointer
// ---------------------------------------------------------------------------

// SafeAudioHandler is a consumer-friendly alternative to AudioHandler
// that passes decoded [][]float32 audio data instead of unsafe.Pointer.
type SafeAudioHandler interface {
	GetAudioParameters(browser Browser, params *AudioParameters) int32
	OnAudioStreamStarted(browser Browser, params *AudioParameters, channels int32)
	OnAudioStreamPacket(browser Browser, data [][]float32, frames int32, pts int64)
	OnAudioStreamStopped(browser Browser)
	OnAudioStreamError(browser Browser, message string)
}

// safeAudioAdapter adapts SafeAudioHandler to AudioHandler.
// It captures the channel count from OnAudioStreamStarted and uses it
// to decode the float** data in OnAudioStreamPacket.
type safeAudioAdapter struct {
	impl     SafeAudioHandler
	mu       sync.Mutex
	channels int32
}

func (a *safeAudioAdapter) GetAudioParameters(browser Browser, params *AudioParameters) int32 {
	return a.impl.GetAudioParameters(browser, params)
}

func (a *safeAudioAdapter) OnAudioStreamStarted(browser Browser, params *AudioParameters, channels int32) {
	a.mu.Lock()
	a.channels = channels
	a.mu.Unlock()
	a.impl.OnAudioStreamStarted(browser, params, channels)
}

func (a *safeAudioAdapter) OnAudioStreamPacket(browser Browser, data unsafe.Pointer, frames int32, pts int64) {
	a.mu.Lock()
	ch := a.channels
	a.mu.Unlock()
	decoded := DecodeAudioPacket(data, ch, frames)
	a.impl.OnAudioStreamPacket(browser, decoded, frames, pts)
}

func (a *safeAudioAdapter) OnAudioStreamStopped(browser Browser) {
	a.impl.OnAudioStreamStopped(browser)
}

func (a *safeAudioAdapter) OnAudioStreamError(browser Browser, message string) {
	a.impl.OnAudioStreamError(browser, message)
}

// NewSafeAudioHandler creates a CEF AudioHandler from a SafeAudioHandler.
func NewSafeAudioHandler(impl SafeAudioHandler) AudioHandler {
	return NewAudioHandler(&safeAudioAdapter{impl: impl})
}
