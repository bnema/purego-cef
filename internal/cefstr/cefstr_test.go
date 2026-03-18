package cefstr

import (
	"testing"
	"unsafe"
)

func TestCEFStringUTF16Layout(t *testing.T) {
	var value CEFStringUTF16
	if unsafe.Sizeof(value) != unsafe.Sizeof(uintptr(0))*3 {
		t.Fatalf("unexpected size: %d", unsafe.Sizeof(value))
	}
}

func TestToGoHandlesNil(t *testing.T) {
	if got := ToGo(nil); got != "" {
		t.Fatalf("got %q", got)
	}
}

func TestToGoReadsUTF16(t *testing.T) {
	buf := []uint16{'h', 'i'}
	value := &CEFStringUTF16{Str: &buf[0], Length: 2}
	if got := ToGo(value); got != "hi" {
		t.Fatalf("got %q", got)
	}
}

func TestFromGoRequiresBoundRuntime(t *testing.T) {
	clearBindingsForTest()
	_, _, err := FromGo("hello")
	if err == nil {
		t.Fatal("expected bind error")
	}
}

func TestFromGoCopiesAndCleansUp(t *testing.T) {
	clearBindingsForTest()
	stringSet = func(src *uint16, srcLen uintptr, out *CEFStringUTF16, copy int32) int32 {
		buf := append([]uint16(nil), unsafe.Slice(src, srcLen)...)
		out.Str = &buf[0]
		out.Length = uintptr(len(buf))
		return 1
	}
	cleared := false
	stringClear = func(out *CEFStringUTF16) {
		cleared = true
		*out = CEFStringUTF16{}
	}
	bound.Store(true)

	value, cleanup, err := FromGo("hello")
	if err != nil {
		t.Fatal(err)
	}
	if ToGo(&value) != "hello" {
		t.Fatalf("got %q", ToGo(&value))
	}
	cleanup()
	if !cleared {
		t.Fatal("cleanup not called")
	}
}
