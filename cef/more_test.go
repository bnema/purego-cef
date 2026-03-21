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

// --- SafeLifeSpanHandler tests ---

type testSafeLifeSpan struct {
	popupAction    PopupAction
	devToolsAction DevToolsPopupAction
	afterCreated   bool
	closed         bool
}

func (h *testSafeLifeSpan) OnBeforePopup(_ Browser, _ Frame, _ int32, _ string,
	_ string, _ WindowOpenDisposition, _ int32, _ *PopupFeatures, _ *WindowInfo,
	_ *BrowserSettings) PopupAction {
	return h.popupAction
}
func (h *testSafeLifeSpan) OnBeforePopupAborted(_ Browser, _ int32) {}
func (h *testSafeLifeSpan) OnBeforeDevToolsPopup(_ Browser, _ *WindowInfo,
	_ *BrowserSettings) DevToolsPopupAction {
	return h.devToolsAction
}
func (h *testSafeLifeSpan) OnAfterCreated(_ Browser) { h.afterCreated = true }
func (h *testSafeLifeSpan) DoClose(_ Browser) bool   { return false }
func (h *testSafeLifeSpan) OnBeforeClose(_ Browser)  { h.closed = true }

func TestSafeLifeSpanAdapterBlock(t *testing.T) {
	impl := &testSafeLifeSpan{popupAction: PopupAction{Block: true}}
	adapter := &safeLifeSpanAdapter{impl: impl}
	var noJS int32
	blocked := adapter.OnBeforePopup(nil, nil, 0, "", "", 0, 0, nil, nil,
		nil, nil, nil, &noJS)
	if !blocked {
		t.Error("expected popup to be blocked")
	}
}

func TestSafeLifeSpanAdapterAllow(t *testing.T) {
	impl := &testSafeLifeSpan{popupAction: PopupAction{
		Block:              false,
		NoJavascriptAccess: true,
	}}
	adapter := &safeLifeSpanAdapter{impl: impl}
	var noJS int32
	blocked := adapter.OnBeforePopup(nil, nil, 0, "", "", 0, 0, nil, nil,
		nil, nil, nil, &noJS)
	if blocked {
		t.Error("expected popup to be allowed")
	}
	if noJS != 1 {
		t.Errorf("noJavascriptAccess = %d, want 1", noJS)
	}
}

func TestSafeLifeSpanAdapterDevTools(t *testing.T) {
	impl := &testSafeLifeSpan{devToolsAction: DevToolsPopupAction{UseDefaultWindow: true}}
	adapter := &safeLifeSpanAdapter{impl: impl}
	var useDefault int32
	adapter.OnBeforeDevToolsPopup(nil, nil, nil, nil, nil, &useDefault)
	if useDefault != 1 {
		t.Errorf("useDefaultWindow = %d, want 1", useDefault)
	}
}

// --- SafeAudioHandler tests ---

type testSafeAudio struct {
	startedChannels int32
	packets         [][][]float32
}

func (h *testSafeAudio) GetAudioParameters(_ Browser, _ *AudioParameters) int32 { return 1 }
func (h *testSafeAudio) OnAudioStreamStarted(_ Browser, _ *AudioParameters, channels int32) {
	h.startedChannels = channels
}
func (h *testSafeAudio) OnAudioStreamPacket(_ Browser, data [][]float32, _ int32, _ int64) {
	h.packets = append(h.packets, data)
}
func (h *testSafeAudio) OnAudioStreamStopped(_ Browser)         {}
func (h *testSafeAudio) OnAudioStreamError(_ Browser, _ string) {}

func TestSafeAudioAdapterDecodesPacket(t *testing.T) {
	impl := &testSafeAudio{}
	adapter := &safeAudioAdapter{impl: impl}

	// Simulate stream start with 2 channels.
	adapter.OnAudioStreamStarted(nil, nil, 2)
	if impl.startedChannels != 2 {
		t.Fatalf("channels = %d, want 2", impl.startedChannels)
	}

	// Simulate a packet.
	ch0 := [3]float32{0.1, 0.2, 0.3}
	ch1 := [3]float32{0.4, 0.5, 0.6}
	ptrs := [2]*float32{&ch0[0], &ch1[0]}
	adapter.OnAudioStreamPacket(nil, unsafe.Pointer(&ptrs[0]), 3, 0)

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
