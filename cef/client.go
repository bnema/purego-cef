package cef

import (
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/refcount"
	"github.com/ebitengine/purego"
)

// Client configures handler interfaces used when creating a browser.
// Use NewClient and the With* methods to build one.
type Client struct {
	renderHandler   RenderHandler
	lifeSpanHandler LifeSpanHandler
}

// NewClient returns a new Client builder.
func NewClient() *Client { return &Client{} }

// WithRenderHandler sets the render handler.
func (c *Client) WithRenderHandler(h RenderHandler) *Client {
	c.renderHandler = h
	return c
}

// WithLifeSpanHandler sets the life-span handler.
func (c *Client) WithLifeSpanHandler(h LifeSpanHandler) *Client {
	c.lifeSpanHandler = h
	return c
}

// clientState holds the CEF C structs and the Go handler references together.
// The refcount system pins this object to prevent GC while CEF holds pointers.
type clientState struct {
	client          capi.CEFClientT
	renderHandler   capi.CEFRenderHandlerT
	lifeSpanHandler capi.CEFLifeSpanHandlerT
	owner           *Client
	created         chan *browser
	closed          chan struct{}
}

// newClientState allocates CEF structs, wires refcount and callbacks.
func newClientState(c *Client) *clientState {
	state := &clientState{
		owner:   c,
		created: make(chan *browser, 1),
		closed:  make(chan struct{}),
	}

	// Initialize refcounting for all three structs. The owner keeps the
	// state pinned in the refcount map so GC cannot collect it.
	refcount.Init(unsafe.Pointer(&state.client.Base), unsafe.Sizeof(state.client), state)
	refcount.Init(unsafe.Pointer(&state.renderHandler.Base), unsafe.Sizeof(state.renderHandler), state)
	refcount.Init(unsafe.Pointer(&state.lifeSpanHandler.Base), unsafe.Sizeof(state.lifeSpanHandler), state)

	// Wire client → handler accessors.
	state.client.OverrideGetRenderHandler(purego.NewCallback(func(self unsafe.Pointer) uintptr {
		return uintptr(unsafe.Pointer(&state.renderHandler))
	}))
	state.client.OverrideGetLifeSpanHandler(purego.NewCallback(func(self unsafe.Pointer) uintptr {
		return uintptr(unsafe.Pointer(&state.lifeSpanHandler))
	}))

	// Wire render handler callbacks.
	if c.renderHandler != nil {
		state.renderHandler.OverrideGetViewRect(purego.NewCallback(func(self unsafe.Pointer, browser unsafe.Pointer, rect unsafe.Pointer) {
			r := c.renderHandler.GetViewRect()
			out := (*capi.CEFRectT)(rect)
			out.X = r.X
			out.Y = r.Y
			out.W = r.W
			out.H = r.H
		}))

		state.renderHandler.OverrideOnPaint(purego.NewCallback(func(self unsafe.Pointer, browser unsafe.Pointer, paintType int32, dirtyRectsCount uintptr, dirtyRects unsafe.Pointer, buffer unsafe.Pointer, width int32, height int32) {
			var rects []Rect
			if dirtyRectsCount > 0 {
				raw := unsafe.Slice((*capi.CEFRectT)(dirtyRects), dirtyRectsCount)
				rects = make([]Rect, dirtyRectsCount)
				for i, r := range raw {
					rects[i] = Rect{X: r.X, Y: r.Y, W: r.W, H: r.H}
				}
			}
			evt := PaintEvent{
				Width:      width,
				Height:     height,
				Buffer:     buffer,
				BufferSize: int(width) * int(height) * 4, // BGRA
				DirtyRects: rects,
			}
			c.renderHandler.OnPaint(evt)
		}))
	}

	// Wire life-span handler callbacks.
	if c.lifeSpanHandler != nil {
		state.lifeSpanHandler.OverrideOnAfterCreated(purego.NewCallback(func(self unsafe.Pointer, rawBrowser unsafe.Pointer) {
			b := &browser{raw: (*capi.CEFBrowserT)(rawBrowser)}
			c.lifeSpanHandler.OnAfterCreated(b)
			// Non-blocking send; the CreateBrowser pump reads from this.
			select {
			case state.created <- b:
			default:
			}
		}))

		state.lifeSpanHandler.OverrideOnBeforeClose(purego.NewCallback(func(self unsafe.Pointer, rawBrowser unsafe.Pointer) {
			b := &browser{raw: (*capi.CEFBrowserT)(rawBrowser)}
			c.lifeSpanHandler.OnBeforeClose(b)
			select {
			case <-state.closed:
			default:
				close(state.closed)
			}
		}))
	} else {
		// Even without a user handler we need OnAfterCreated to signal
		// browser creation to CreateBrowser.
		state.lifeSpanHandler.OverrideOnAfterCreated(purego.NewCallback(func(self unsafe.Pointer, rawBrowser unsafe.Pointer) {
			b := &browser{raw: (*capi.CEFBrowserT)(rawBrowser)}
			select {
			case state.created <- b:
			default:
			}
		}))

		state.lifeSpanHandler.OverrideOnBeforeClose(purego.NewCallback(func(self unsafe.Pointer, browser unsafe.Pointer) {
			select {
			case <-state.closed:
			default:
				close(state.closed)
			}
		}))
	}

	return state
}
