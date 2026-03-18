package refcount

import (
	"testing"
	"unsafe"

	"github.com/ebitengine/purego"
)

type fakeBase struct {
	Size             uintptr
	AddRef           uintptr
	Release          uintptr
	HasOneRef        uintptr
	HasAtLeastOneRef uintptr
}

func TestInitSetsCallbacksAndPinsOwner(t *testing.T) {
	var base fakeBase
	owner := &struct{ Name string }{"client"}
	state := Init(unsafe.Pointer(&base), unsafe.Sizeof(base), owner)
	if base.AddRef == 0 || base.Release == 0 || base.HasOneRef == 0 {
		t.Fatal("callbacks not initialized")
	}
	if state.Count() != 1 {
		t.Fatalf("count = %d", state.Count())
	}
	if _, ok := Load(unsafe.Pointer(&base)); !ok {
		t.Fatal("owner not pinned")
	}
}

func TestReleaseUnpinsAtZero(t *testing.T) {
	var base fakeBase
	Init(unsafe.Pointer(&base), unsafe.Sizeof(base), &struct{}{})
	AddRef(unsafe.Pointer(&base))
	if got := Release(unsafe.Pointer(&base)); got != 0 {
		t.Fatalf("first release = %d", got)
	}
	if got := Release(unsafe.Pointer(&base)); got != 1 {
		t.Fatalf("second release = %d", got)
	}
	if _, ok := Load(unsafe.Pointer(&base)); ok {
		t.Fatal("owner still pinned")
	}
}

func TestCallbacksAreCallable(t *testing.T) {
	var base fakeBase
	Init(unsafe.Pointer(&base), unsafe.Sizeof(base), &struct{}{})
	var addRef func(unsafe.Pointer)
	var release func(unsafe.Pointer) int32
	purego.RegisterFunc(&addRef, base.AddRef)
	purego.RegisterFunc(&release, base.Release)
	addRef(unsafe.Pointer(&base))
	if got := release(unsafe.Pointer(&base)); got != 0 {
		t.Fatalf("expected release to return 0 (object still alive), got %d", got)
	}
}
