package cef

import (
	"testing"
	"unsafe"
)

func TestNewKeyEvent(t *testing.T) {
	evt := NewKeyEvent(KeyEventTypeKeyeventChar, 65, 0, 0)
	if evt.Size != unsafe.Sizeof(evt) {
		t.Errorf("Size = %d, want %d", evt.Size, unsafe.Sizeof(evt))
	}
	if evt.WindowsKeyCode != 65 {
		t.Errorf("WindowsKeyCode = %d, want 65", evt.WindowsKeyCode)
	}
	typeVal := *(*int32)(unsafe.Pointer(&evt.Type))
	if typeVal != int32(KeyEventTypeKeyeventChar) {
		t.Errorf("Type = %d, want %d", typeVal, KeyEventTypeKeyeventChar)
	}
}

func TestKeyEventSetType(t *testing.T) {
	evt := NewKeyEvent(KeyEventTypeKeyeventKeydown, 0, 0, 0)
	KeyEventSetType(&evt, KeyEventTypeKeyeventChar)
	typeVal := *(*int32)(unsafe.Pointer(&evt.Type))
	if typeVal != int32(KeyEventTypeKeyeventChar) {
		t.Errorf("Type = %d, want %d", typeVal, KeyEventTypeKeyeventChar)
	}
}

func TestNewScreenInfo(t *testing.T) {
	rect := Rect{}
	si := NewScreenInfo(2.0, 24, 8, false, rect, rect)
	if si.Size != unsafe.Sizeof(si) {
		t.Errorf("Size = %d, want %d", si.Size, unsafe.Sizeof(si))
	}
	if si.DeviceScaleFactor != 2.0 {
		t.Errorf("DeviceScaleFactor = %f, want 2.0", si.DeviceScaleFactor)
	}
	if si.IsMonochrome != 0 {
		t.Errorf("IsMonochrome = %d, want 0", si.IsMonochrome)
	}
}

func TestNewScreenInfoMonochrome(t *testing.T) {
	rect := Rect{}
	si := NewScreenInfo(1.0, 1, 1, true, rect, rect)
	if si.IsMonochrome != 1 {
		t.Errorf("IsMonochrome = %d, want 1", si.IsMonochrome)
	}
}

func TestDecodeAudioPacket(t *testing.T) {
	ch0 := [3]float32{0.1, 0.2, 0.3}
	ch1 := [3]float32{0.4, 0.5, 0.6}
	channels := [2]*float32{&ch0[0], &ch1[0]}
	result := DecodeAudioPacket(unsafe.Pointer(&channels[0]), 2, 3)
	if len(result) != 2 {
		t.Fatalf("expected 2 channels, got %d", len(result))
	}
	if len(result[0]) != 3 {
		t.Fatalf("expected 3 frames, got %d", len(result[0]))
	}
	if result[0][0] != 0.1 || result[0][2] != 0.3 {
		t.Errorf("channel 0: got %v, want [0.1 0.2 0.3]", result[0])
	}
	if result[1][0] != 0.4 || result[1][2] != 0.6 {
		t.Errorf("channel 1: got %v, want [0.4 0.5 0.6]", result[1])
	}
}

func TestDecodeAudioPacketNilData(t *testing.T) {
	result := DecodeAudioPacket(nil, 2, 3)
	if result != nil {
		t.Errorf("expected nil for nil data, got %v", result)
	}
}

func TestDecodeAudioPacketZeroChannels(t *testing.T) {
	ch0 := [3]float32{0.1, 0.2, 0.3}
	channels := [1]*float32{&ch0[0]}
	result := DecodeAudioPacket(unsafe.Pointer(&channels[0]), 0, 3)
	if result != nil {
		t.Errorf("expected nil for zero channels, got %v", result)
	}
}

func TestDecodeAudioPacketZeroFrames(t *testing.T) {
	ch0 := [3]float32{0.1, 0.2, 0.3}
	channels := [1]*float32{&ch0[0]}
	result := DecodeAudioPacket(unsafe.Pointer(&channels[0]), 1, 0)
	if result != nil {
		t.Errorf("expected nil for zero frames, got %v", result)
	}
}
