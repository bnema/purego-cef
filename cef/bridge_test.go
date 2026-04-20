package cef

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
)

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
