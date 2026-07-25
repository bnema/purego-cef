package capi

import (
	"testing"
	"unsafe"

	"github.com/bnema/purego"
)

func BenchmarkCallSendMouseMoveEvent(b *testing.B) {
	callback := purego.NewCallback(func(_ *CEFBrowserHostT, _ *CEFMouseEventT, _ int32) {})
	b.Cleanup(func() {
		if err := purego.UnrefCallback(callback); err != nil {
			b.Fatalf("unref callback: %v", err)
		}
	})

	host := &CEFBrowserHostT{}
	host.OverrideSendMouseMoveEvent(callback)
	event := &CEFMouseEventT{}
	eventPointer := unsafe.Pointer(event)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		host.CallSendMouseMoveEvent(eventPointer, 0)
	}
}
