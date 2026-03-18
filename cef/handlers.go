package cef

import "unsafe"

// PaintEvent contains the pixel buffer and metadata delivered by CEF's
// OnPaint callback.
type PaintEvent struct {
	Width, Height int32
	Buffer        unsafe.Pointer
	BufferSize    int
	DirtyRects    []Rect
}

// RenderHandler receives paint and geometry callbacks from the browser.
type RenderHandler interface {
	GetViewRect() Rect
	OnPaint(PaintEvent)
}

// LifeSpanHandler receives browser lifecycle notifications.
type LifeSpanHandler interface {
	OnAfterCreated(Browser)
	OnBeforeClose(Browser)
}

// RenderHandlerFunc adapts a simple paint callback into a full RenderHandler.
// GetViewRect returns a default 800x600 rectangle.
type RenderHandlerFunc func(PaintEvent)

// GetViewRect returns a default 800x600 rectangle.
func (f RenderHandlerFunc) GetViewRect() Rect { return Rect{W: 800, H: 600} }

// OnPaint delegates to the wrapped function.
func (f RenderHandlerFunc) OnPaint(e PaintEvent) { f(e) }
