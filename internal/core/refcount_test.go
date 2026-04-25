package core

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/mock"

	"github.com/bnema/purego-cef/internal/ports/out/mocks"
)

type testRefCounted struct {
	size             uintptr
	addRef           uintptr
	release          uintptr
	hasOneRef        uintptr
	hasAtLeastOneRef uintptr
}

func newRefManagerForTest(t *testing.T) (*RefManager, *[]uintptr) {
	m := mocks.NewMockCAPI(t)
	var next uintptr = 100
	m.EXPECT().NewCallback(mock.Anything).RunAndReturn(func(any) uintptr {
		next++
		return next
	}).Maybe()
	released := []uintptr{}
	m.EXPECT().UnrefCallback(mock.Anything).RunAndReturn(func(cb uintptr) error {
		released = append(released, cb)
		return nil
	}).Maybe()
	return NewRefManager(m), &released
}

func TestRefManagerReleasesTrackedCallbacksAtZeroRefs(t *testing.T) {
	rm, released := newRefManagerForTest(t)
	obj := &testRefCounted{}
	base := unsafe.Pointer(obj)
	rm.InitRefCount(base, unsafe.Sizeof(*obj), obj)

	rm.TrackCallback(base, 10)
	rm.TrackCallback(base, 11)

	if got := rm.release(base); got != 1 {
		t.Fatalf("release returned %d, want 1", got)
	}
	if len(*released) != 2 || (*released)[0] != 10 || (*released)[1] != 11 {
		t.Fatalf("released callbacks = %v, want [10 11]", *released)
	}
	if rm.Has(base) {
		t.Fatalf("expected ref state to be removed after final release")
	}
}

func TestRefManagerDoesNotReleaseTrackedCallbacksBeforeZeroRefs(t *testing.T) {
	rm, released := newRefManagerForTest(t)
	obj := &testRefCounted{}
	base := unsafe.Pointer(obj)
	rm.InitRefCount(base, unsafe.Sizeof(*obj), obj)
	rm.AddRef(base)

	rm.TrackCallback(base, 10)

	if got := rm.release(base); got != 0 {
		t.Fatalf("first release returned %d, want 0", got)
	}
	if len(*released) != 0 {
		t.Fatalf("released callbacks before final release: %d", len(*released))
	}
	if got := rm.release(base); got != 1 {
		t.Fatalf("second release returned %d, want 1", got)
	}
	if len(*released) != 1 {
		t.Fatalf("released callbacks = %d, want 1", len(*released))
	}
}

func TestRefManagerTrackCallbackUnknownBaseReleasesAndPanics(t *testing.T) {
	rm, released := newRefManagerForTest(t)
	obj := &testRefCounted{}
	base := unsafe.Pointer(obj)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic")
		}
		if len(*released) != 1 {
			t.Fatalf("released callbacks = %d, want 1", len(*released))
		}
	}()

	rm.TrackCallback(base, 10)
}
