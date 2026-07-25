package cef

import (
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/bnema/purego"
	"github.com/bnema/purego-cef/internal/capi"
)

func TestOutboundRefPtrTransferStaysAliveThroughBlockingDispatch(t *testing.T) {
	const helperEnv = "PUREGO_CEF_REFPTR_TRANSFER_HELPER"
	if os.Getenv(helperEnv) == "" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestOutboundRefPtrTransferStaysAliveThroughBlockingDispatch$", "-test.count=1")
		cmd.Env = append(os.Environ(), helperEnv+"=1")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("refptr transfer helper failed: %v\n%s", err, output)
		}
		return
	}

	// Forced GC must run in an isolated helper process because unrelated
	// constructor tests intentionally leave Go-backed CEF objects pinned.
	var refs atomic.Int32
	var addRefs atomic.Int32
	var releases atomic.Int32
	released := make(chan int32, 4)

	addRefCallback := purego.NewCallback(func(uintptr) uintptr {
		addRefs.Add(1)
		return uintptr(refs.Add(1))
	})
	releaseCallback := purego.NewCallback(func(uintptr) uintptr {
		releases.Add(1)
		remaining := refs.Add(-1)
		released <- remaining
		return uintptr(remaining)
	})
	t.Cleanup(func() {
		if err := purego.UnrefCallback(addRefCallback); err != nil {
			t.Errorf("unref AddRef callback: %v", err)
		}
		if err := purego.UnrefCallback(releaseCallback); err != nil {
			t.Errorf("unref Release callback: %v", err)
		}
	})

	dragRaw := new(capi.CEFDragDataT)
	dragRaw.Base.OverrideAddRef(addRefCallback)
	dragRaw.Base.OverrideRelease(releaseCallback)
	drag := wrapDragData(unsafe.Pointer(dragRaw))
	if got := refs.Load(); got != 1 {
		t.Fatalf("owner wrap refcount = %d, want 1", got)
	}

	entered := make(chan int32, 1)
	allowReturn := make(chan struct{})
	dispatchCallback := purego.NewCallback(func(_, dragData, _, _ uintptr) uintptr {
		entered <- refs.Load()
		<-allowReturn
		(*capi.CEFBaseRefCountedT)(cefCallbackPointer(dragData)).CallRelease()
		return 0
	})
	t.Cleanup(func() {
		if err := purego.UnrefCallback(dispatchCallback); err != nil {
			t.Errorf("unref dispatch callback: %v", err)
		}
	})

	hostRaw := new(capi.CEFBrowserHostT)
	hostRaw.OverrideDragTargetDragEnter(dispatchCallback)
	host := &browserHostImpl{rawPtr: hostRaw}
	done := make(chan struct{})
	go func(arg DragData) {
		host.DragTargetDragEnter(arg, nil, 0)
		close(done)
	}(drag)
	drag = nil

	select {
	case got := <-entered:
		if got != 2 {
			t.Fatalf("refcount at native dispatch entry = %d, want owner + transfer = 2", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for native dispatch entry")
	}

	for range 10 {
		runtime.GC()
		_ = make([]byte, 1<<20)
		runtime.Gosched()
	}
	select {
	case got := <-released:
		t.Fatalf("owner released during blocked native dispatch; remaining refs = %d", got)
	case <-time.After(100 * time.Millisecond):
	}

	close(allowReturn)
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for native dispatch return")
	}
	select {
	case got := <-released:
		if got != 1 {
			t.Fatalf("remaining refs after transfer release = %d, want owner ref 1", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for transfer reference release")
	}

	deadline := time.Now().Add(5 * time.Second)
	for releases.Load() < 2 && time.Now().Before(deadline) {
		runtime.GC()
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
	if got := releases.Load(); got != 2 {
		t.Fatalf("release calls = %d, want transfer and owner releases exactly once each", got)
	}
	if got := addRefs.Load(); got != 2 {
		t.Fatalf("AddRef calls = %d, want owner and transfer AddRef exactly once each", got)
	}
	if got := refs.Load(); got != 0 {
		t.Fatalf("final refcount = %d, want 0", got)
	}
}

func TestOutboundRefPtrMissingFunctionDoesNotTransfer(t *testing.T) {
	var addRefs atomic.Int32
	addRefCallback := purego.NewCallback(func(uintptr) uintptr {
		return uintptr(addRefs.Add(1))
	})
	releaseCallback := purego.NewCallback(func(uintptr) uintptr { return 0 })
	t.Cleanup(func() {
		if err := purego.UnrefCallback(addRefCallback); err != nil {
			t.Errorf("unref AddRef callback: %v", err)
		}
		if err := purego.UnrefCallback(releaseCallback); err != nil {
			t.Errorf("unref Release callback: %v", err)
		}
	})

	dragRaw := new(capi.CEFDragDataT)
	dragRaw.Base.OverrideAddRef(addRefCallback)
	dragRaw.Base.OverrideRelease(releaseCallback)
	drag := wrapDragData(unsafe.Pointer(dragRaw))
	host := &browserHostImpl{rawPtr: new(capi.CEFBrowserHostT)}
	host.DragTargetDragEnter(drag, nil, 0)
	runtime.KeepAlive(drag)

	if got := addRefs.Load(); got != 1 {
		t.Fatalf("AddRef calls with missing native function = %d, want owner AddRef only", got)
	}
	drag.(*dragDataImpl).Release()
}
