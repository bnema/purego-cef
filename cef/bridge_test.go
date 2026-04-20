package cef

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/mock"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/core"
	portoutmocks "github.com/bnema/purego-cef/internal/ports/out/mocks"
)

func newBridgeTestEngine(t *testing.T) *core.Engine {
	m := portoutmocks.NewMockCAPI(t)
	m.EXPECT().NewCallback(mock.Anything).Maybe().Return(uintptr(0))
	return core.New(m)
}

func installBridgeTestEngine(t *testing.T) {
	prevEng := eng
	prevCurrentRefManager := currentRefManager
	prevRegisteredRefManagers := append([]*core.RefManager(nil), registeredRefManagers...)

	e := newBridgeTestEngine(t)
	eng = e
	setCurrentRefManager(e.Refs())

	t.Cleanup(func() {
		eng = prevEng
		currentRefManager = prevCurrentRefManager
		registeredRefManagers = prevRegisteredRefManagers
	})
}

type plainRawClientStub struct{}

type nilAudioHandlerStub struct{}
type plainLifeSpanHandlerStub struct{}

func (*nilAudioHandlerStub) GetAudioParameters(Browser, *AudioParameters) int32     { return 0 }
func (*nilAudioHandlerStub) OnAudioStreamStarted(Browser, *AudioParameters, int32)  {}
func (*nilAudioHandlerStub) OnAudioStreamPacket(Browser, [][]float32, int32, int64) {}
func (*nilAudioHandlerStub) OnAudioStreamStopped(Browser)                           {}
func (*nilAudioHandlerStub) OnAudioStreamError(Browser, string)                     {}

func (plainLifeSpanHandlerStub) OnBeforePopup(Browser, Frame, int32, string, string, WindowOpenDisposition, int32, *PopupFeatures, *WindowInfo, *RawClientWriteSlot, *BrowserSettings, *DictionaryValue, *bool) bool {
	return false
}
func (plainLifeSpanHandlerStub) OnBeforePopupAborted(Browser, int32) {}
func (plainLifeSpanHandlerStub) OnBeforeDevToolsPopup(Browser, *WindowInfo, *RawClientWriteSlot, *BrowserSettings, *DictionaryValue, *bool) {
}
func (plainLifeSpanHandlerStub) OnAfterCreated(Browser) {}
func (plainLifeSpanHandlerStub) DoClose(Browser) bool   { return false }
func (plainLifeSpanHandlerStub) OnBeforeClose(Browser)  {}

func (plainRawClientStub) GetAudioHandler() RawAudioHandler          { return nil }
func (plainRawClientStub) GetCommandHandler() CommandHandler         { return nil }
func (plainRawClientStub) GetContextMenuHandler() ContextMenuHandler { return nil }
func (plainRawClientStub) GetDialogHandler() DialogHandler           { return nil }
func (plainRawClientStub) GetDisplayHandler() DisplayHandler         { return nil }
func (plainRawClientStub) GetDownloadHandler() DownloadHandler       { return nil }
func (plainRawClientStub) GetDragHandler() DragHandler               { return nil }
func (plainRawClientStub) GetFindHandler() FindHandler               { return nil }
func (plainRawClientStub) GetFocusHandler() FocusHandler             { return nil }
func (plainRawClientStub) GetFrameHandler() FrameHandler             { return nil }
func (plainRawClientStub) GetPermissionHandler() PermissionHandler   { return nil }
func (plainRawClientStub) GetJsdialogHandler() JsdialogHandler       { return nil }
func (plainRawClientStub) GetKeyboardHandler() KeyboardHandler       { return nil }
func (plainRawClientStub) GetLifeSpanHandler() RawLifeSpanHandler    { return nil }
func (plainRawClientStub) GetLoadHandler() LoadHandler               { return nil }
func (plainRawClientStub) GetPrintHandler() PrintHandler             { return nil }
func (plainRawClientStub) GetRenderHandler() RenderHandler           { return nil }
func (plainRawClientStub) GetRequestHandler() RequestHandler         { return nil }
func (plainRawClientStub) OnProcessMessageReceived(Browser, Frame, ProcessID, ProcessMessage) int32 {
	return 0
}

func TestConstructorsReturnNilForNilImpl(t *testing.T) {
	var rawLifeSpan RawLifeSpanHandler
	if got := NewRawLifeSpanHandler(rawLifeSpan); got != nil {
		t.Fatalf("NewRawLifeSpanHandler(nil) = %#v, want nil", got)
	}

	var lifeSpan LifeSpanHandler
	if got := NewLifeSpanHandler(lifeSpan); got != nil {
		t.Fatalf("NewLifeSpanHandler(nil) = %#v, want nil", got)
	}

	var audio AudioHandler
	if got := NewAudioHandler(audio); got != nil {
		t.Fatalf("NewAudioHandler(nil) = %#v, want nil", got)
	}

	var rawAudio RawAudioHandler
	if got := NewRawAudioHandler(rawAudio); got != nil {
		t.Fatalf("NewRawAudioHandler(nil) = %#v, want nil", got)
	}

	var client Client
	if got := NewClient(client); got != nil {
		t.Fatalf("NewClient(nil) = %#v, want nil", got)
	}
}

func TestConstructorsReturnNilForTypedNilImpl(t *testing.T) {
	var impl AudioHandler = (*nilAudioHandlerStub)(nil)
	if got := NewAudioHandler(impl); got != nil {
		t.Fatalf("NewAudioHandler(typed nil) = %#v, want nil", got)
	}
}

func TestNewLifeSpanHandlerUsesDirectPopupSlotAPI(t *testing.T) {
	installBridgeTestEngine(t)

	if got := NewLifeSpanHandler(plainLifeSpanHandlerStub{}); got == nil {
		t.Fatal("NewLifeSpanHandler(...) = nil, want non-nil raw handler")
	}
}

func TestRawClientWriteSlotPreservesInitialRawPointerUntilTouched(t *testing.T) {
	initial := unsafe.Pointer(&capi.CEFClientT{})
	slot := newRawClientWriteSlot(initial)

	if got := slot.rawPointer(); got != initial {
		t.Fatalf("rawPointer() = %p, want %p", got, initial)
	}
}

func TestRawClientWriteSlotClearWritesNil(t *testing.T) {
	slot := newRawClientWriteSlot(unsafe.Pointer(&capi.CEFClientT{}))
	slot.Clear()

	if got := slot.rawPointer(); got != nil {
		t.Fatalf("rawPointer() = %p, want nil", got)
	}
}

func TestRawClientWriteSlotSetUsesProvidedRawClient(t *testing.T) {
	want := &capi.CEFClientT{}
	slot := newRawClientWriteSlot(unsafe.Pointer(&capi.CEFClientT{}))
	slot.Set(&rawClientWrapper{rawPtr: want})

	if got := slot.rawPointer(); got != unsafe.Pointer(want) {
		t.Fatalf("rawPointer() = %p, want %p", got, want)
	}
}

func TestRawClientWriteSlotSetWrapsPlainRawClient(t *testing.T) {
	installBridgeTestEngine(t)

	initial := unsafe.Pointer(&capi.CEFClientT{})
	slot := newRawClientWriteSlot(initial)
	slot.Set(plainRawClientStub{})

	got := slot.rawPointer()
	if got == nil {
		t.Fatal("rawPointer() = nil, want wrapped raw client pointer")
	}
	if got == initial {
		t.Fatalf("rawPointer() = %p, want wrapped client pointer distinct from initial raw %p", got, initial)
	}
}

func TestAddRefReportsUnknownPointersToDebugHook(t *testing.T) {
	prev := debugRefCountf
	defer func() { debugRefCountf = prev }()

	called := false
	var gotFormat string
	var gotArgs []any
	debugRefCountf = func(format string, args ...any) {
		called = true
		gotFormat = format
		gotArgs = args
	}

	base := unsafe.Pointer(&capi.CEFClientT{})
	addRef(base)

	if !called {
		t.Fatal("addRef(...) did not report unknown pointer")
	}
	if !strings.Contains(gotFormat, "no RefManager found") {
		t.Fatalf("debug format = %q, want no-manager message", gotFormat)
	}
	if len(gotArgs) != 1 || gotArgs[0] != base {
		t.Fatalf("debug args = %#v, want [%p]", gotArgs, base)
	}
}

func TestRawClientWriteSlotSetPanicsOnNilOrUninitializedSlot(t *testing.T) {
	t.Run("nil slot", func(t *testing.T) {
		var slot *RawClientWriteSlot
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic")
			}
			msg, ok := r.(string)
			if !ok || !strings.Contains(msg, "slot is nil") {
				t.Fatalf("panic = %#v, want nil-slot message", r)
			}
		}()
		slot.Set(nil)
	})

	t.Run("zero value slot", func(t *testing.T) {
		slot := &RawClientWriteSlot{}
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic")
			}
			msg, ok := r.(string)
			if !ok || !strings.Contains(msg, "no longer valid") {
				t.Fatalf("panic = %#v, want invalid-slot message", r)
			}
		}()
		slot.Clear()
	})
}

func TestRawClientWriteSlotPanicsAfterInvalidation(t *testing.T) {
	slot := newRawClientWriteSlot(nil)
	slot.invalidate()

	for _, tc := range []struct {
		name string
		run  func()
	}{
		{name: "Set", run: func() { slot.Set(nil) }},
		{name: "Clear", run: func() { slot.Clear() }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("expected panic")
				}
				msg, ok := r.(string)
				if !ok || !strings.Contains(msg, "no longer valid") {
					t.Fatalf("panic = %#v, want invalid-slot message", r)
				}
			}()
			tc.run()
		})
	}
}

func TestAudioHandlerWrapperOnAudioStreamPacketPanics(t *testing.T) {
	w := &audioHandlerWrapper{}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("panic type = %T, want string", r)
		}
		if !strings.Contains(msg, "raw OnAudioStreamPacket called directly") {
			t.Fatalf("panic = %q, want raw misuse message", msg)
		}
	}()

	w.OnAudioStreamPacket(nil, nil, 0, 0)
}
