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

type plainRawClientStub struct{}

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

func TestRawClientWriteSlotPreservesInitialRawPointerUntilTouched(t *testing.T) {
	initial := unsafe.Pointer(&capi.CEFClientT{})
	slot := &RawClientWriteSlot{initialRaw: initial}

	if got := slot.rawPointer(); got != initial {
		t.Fatalf("rawPointer() = %p, want %p", got, initial)
	}
}

func TestRawClientWriteSlotClearWritesNil(t *testing.T) {
	slot := &RawClientWriteSlot{initialRaw: unsafe.Pointer(&capi.CEFClientT{})}
	slot.Clear()

	if got := slot.rawPointer(); got != nil {
		t.Fatalf("rawPointer() = %p, want nil", got)
	}
}

func TestRawClientWriteSlotSetUsesProvidedRawClient(t *testing.T) {
	want := &capi.CEFClientT{}
	slot := &RawClientWriteSlot{initialRaw: unsafe.Pointer(&capi.CEFClientT{})}
	slot.Set(&rawClientWrapper{rawPtr: want})

	if got := slot.rawPointer(); got != unsafe.Pointer(want) {
		t.Fatalf("rawPointer() = %p, want %p", got, want)
	}
}

func TestRawClientWriteSlotSetWrapsPlainRawClient(t *testing.T) {
	prevEng := eng
	eng = newBridgeTestEngine(t)
	t.Cleanup(func() { eng = prevEng })

	initial := unsafe.Pointer(&capi.CEFClientT{})
	slot := &RawClientWriteSlot{initialRaw: initial}
	slot.Set(plainRawClientStub{})

	got := slot.rawPointer()
	if got == nil {
		t.Fatal("rawPointer() = nil, want wrapped raw client pointer")
	}
	if got == initial {
		t.Fatalf("rawPointer() = %p, want wrapped client pointer distinct from initial raw %p", got, initial)
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
