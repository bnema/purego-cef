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

// --- LifeSpanHandler tests ---

type testLifeSpan struct {
	blocked    bool
	noJS       bool
	useDefault bool
}

func (h *testLifeSpan) OnBeforePopup(_ Browser, _ Frame, _ int32, _ string,
	_ string, _ WindowOpenDisposition, _ int32, _ *PopupFeatures, _ *WindowInfo,
	_ *Client, _ *BrowserSettings, _ *DictionaryValue, noJS *bool) bool {
	if noJS != nil {
		*noJS = h.noJS
	}
	return h.blocked
}
func (h *testLifeSpan) OnBeforePopupAborted(_ Browser, _ int32) {}
func (h *testLifeSpan) OnBeforeDevToolsPopup(_ Browser, _ *WindowInfo,
	_ *Client, _ *BrowserSettings, _ *DictionaryValue, useDefault *bool) {
	if useDefault != nil {
		*useDefault = h.useDefault
	}
}
func (h *testLifeSpan) OnAfterCreated(_ Browser) {}
func (h *testLifeSpan) DoClose(_ Browser) bool   { return false }
func (h *testLifeSpan) OnBeforeClose(_ Browser)  {}

func TestLifeSpanHandlerBlock(t *testing.T) {
	impl := &testLifeSpan{blocked: true}
	var noJS bool
	blocked := impl.OnBeforePopup(nil, nil, 0, "", "", 0, 0, nil, nil,
		nil, nil, nil, &noJS)
	if !blocked {
		t.Error("expected popup to be blocked")
	}
}

func TestLifeSpanHandlerNoJS(t *testing.T) {
	impl := &testLifeSpan{blocked: false, noJS: true}
	var noJS bool
	blocked := impl.OnBeforePopup(nil, nil, 0, "", "", 0, 0, nil, nil,
		nil, nil, nil, &noJS)
	if blocked {
		t.Error("expected popup to be allowed")
	}
	if !noJS {
		t.Error("expected noJavascriptAccess to be true")
	}
}

func TestLifeSpanHandlerDevTools(t *testing.T) {
	impl := &testLifeSpan{useDefault: true}
	var useDefault bool
	impl.OnBeforeDevToolsPopup(nil, nil, nil, nil, nil, &useDefault)
	if !useDefault {
		t.Error("expected useDefaultWindow to be true")
	}
}

// --- AudioHandler tests ---

type testAudio struct {
	startedChannels int32
	packets         [][][]float32
}

func (h *testAudio) GetAudioParameters(_ Browser, _ *AudioParameters) int32 { return 1 }
func (h *testAudio) OnAudioStreamStarted(_ Browser, _ *AudioParameters, channels int32) {
	h.startedChannels = channels
}
func (h *testAudio) OnAudioStreamPacket(_ Browser, data [][]float32, _ int32, _ int64) {
	h.packets = append(h.packets, data)
}
func (h *testAudio) OnAudioStreamStopped(_ Browser)         {}
func (h *testAudio) OnAudioStreamError(_ Browser, _ string) {}

func TestAudioHandlerDecodesPacket(t *testing.T) {
	impl := &testAudio{}
	w := NewAudioHandler(impl).(*audioHandlerWrapper)

	// Simulate stream start with 2 channels.
	w.mu.Lock()
	w.channels = 2
	w.mu.Unlock()
	impl.OnAudioStreamStarted(nil, nil, 2)
	if impl.startedChannels != 2 {
		t.Fatalf("channels = %d, want 2", impl.startedChannels)
	}

	// Test the decode path directly.
	ch0 := [3]float32{0.1, 0.2, 0.3}
	ch1 := [3]float32{0.4, 0.5, 0.6}
	ptrs := [2]*float32{&ch0[0], &ch1[0]}
	decoded := DecodeAudioPacket(unsafe.Pointer(&ptrs[0]), 2, 3)
	impl.OnAudioStreamPacket(nil, decoded, 3, 0)

	if len(impl.packets) != 1 {
		t.Fatalf("expected 1 packet, got %d", len(impl.packets))
	}
	pkt := impl.packets[0]
	if len(pkt) != 2 {
		t.Fatalf("expected 2 channels, got %d", len(pkt))
	}
	if pkt[0][0] != 0.1 || pkt[1][2] != 0.6 {
		t.Errorf("unexpected data: ch0=%v ch1=%v", pkt[0], pkt[1])
	}
}
