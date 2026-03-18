package cef

import (
	"testing"
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
)

func TestSettingsSizeField(t *testing.T) {
	// Verify that Size is set correctly when converting to C.
	// We can't call toC() without the runtime bound, but we can check
	// the struct size directly.
	var cs capi.CEFSettingsT
	sz := unsafe.Sizeof(cs)
	if sz == 0 {
		t.Fatal("CEFSettingsT has zero size")
	}
	// On 64-bit Linux the struct should be several hundred bytes.
	if sz < 200 {
		t.Fatalf("CEFSettingsT size %d seems too small", sz)
	}
	t.Logf("CEFSettingsT size = %d", sz)
}

func TestMainArgsPtrNotNil(t *testing.T) {
	args := NewMainArgs([]string{"test"})
	if args.Ptr() == nil {
		t.Fatal("MainArgs.Ptr() returned nil")
	}
}

func TestWindowInfoSize(t *testing.T) {
	var wi capi.CEFWindowInfoT
	sz := unsafe.Sizeof(wi)
	if sz == 0 {
		t.Fatal("CEFWindowInfoT has zero size")
	}
	t.Logf("CEFWindowInfoT size = %d", sz)
}

func TestBrowserSettingsSize(t *testing.T) {
	var bs capi.CEFBrowserSettingsT
	sz := unsafe.Sizeof(bs)
	if sz == 0 {
		t.Fatal("CEFBrowserSettingsT has zero size")
	}
	t.Logf("CEFBrowserSettingsT size = %d", sz)
}

func TestMouseEventSize(t *testing.T) {
	var me capi.CEFMouseEventT
	sz := unsafe.Sizeof(me)
	// 3 fields: int32, int32, uint32 = 12 bytes
	if sz != 12 {
		t.Fatalf("CEFMouseEventT size = %d, want 12", sz)
	}
}
