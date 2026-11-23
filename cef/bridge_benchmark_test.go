package cef

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/mock"

	"bnema.dev/purego-cef/internal/capi"
	"bnema.dev/purego-cef/internal/core"
	portoutmocks "bnema.dev/purego-cef/internal/ports/out/mocks"
)

var callbackOwnerBenchmarkSink *callbackOwnerPrimary

func BenchmarkCEFCallbackOwnerAs(b *testing.B) {
	m := portoutmocks.NewMockCAPI(b)
	m.EXPECT().NewCallback(mock.Anything).Maybe().Return(uintptr(0))
	rm := core.New(m).Refs()
	base := &capi.CEFClientT{}
	owner := &callbackOwnerPrimary{}
	rm.InitRefCount(unsafe.Pointer(base), unsafe.Sizeof(*base), owner)

	withCallbackOwnerRegistry(b, func() {
		registerRefManager(rm)
		self := uintptr(unsafe.Pointer(base))
		b.ReportAllocs()
		b.ResetTimer()
		for b.Loop() {
			got, ok := cefCallbackOwnerAs[*callbackOwnerPrimary](self)
			if !ok {
				b.Fatal("cefCallbackOwnerAs() = no owner")
			}
			callbackOwnerBenchmarkSink = got
		}
	})
}
