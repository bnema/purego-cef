package cef

import "github.com/bnema/purego-cef/internal/capi"

// MouseEvent describes a mouse event to be forwarded to a browser.
type MouseEvent struct {
	X, Y      int32
	Modifiers uint32
}

// Event flag constants matching cef_event_flags_t.
const (
	EventFlagLeftMouseButton   uint32 = 1 << 5
	EventFlagMiddleMouseButton uint32 = 1 << 6
	EventFlagRightMouseButton  uint32 = 1 << 7
)

// toC converts a public MouseEvent to the internal capi type.
func (e MouseEvent) toC() capi.CEFMouseEventT {
	return capi.CEFMouseEventT{X: e.X, Y: e.Y, Modifiers: e.Modifiers}
}
