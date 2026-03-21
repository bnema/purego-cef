package cef

import "unsafe"

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
