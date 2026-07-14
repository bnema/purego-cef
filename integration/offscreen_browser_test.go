//go:build cef_integration

package integration

import (
	"testing"

	"github.com/bnema/purego-cef/cef"
)

// testClient is a minimal Client implementation that returns nil for all
// sub-handlers except the configured render and life span handlers.
type testClient struct {
	renderHandler   cef.RenderHandler
	lifeSpanHandler cef.LifeSpanHandler
}

func (c *testClient) GetAudioHandler() cef.AudioHandler             { return nil }
func (c *testClient) GetCommandHandler() cef.CommandHandler         { return nil }
func (c *testClient) GetContextMenuHandler() cef.ContextMenuHandler { return nil }
func (c *testClient) GetDialogHandler() cef.DialogHandler           { return nil }
func (c *testClient) GetDisplayHandler() cef.DisplayHandler         { return nil }
func (c *testClient) GetDownloadHandler() cef.DownloadHandler       { return nil }
func (c *testClient) GetDragHandler() cef.DragHandler               { return nil }
func (c *testClient) GetFindHandler() cef.FindHandler               { return nil }
func (c *testClient) GetFocusHandler() cef.FocusHandler             { return nil }
func (c *testClient) GetFrameHandler() cef.FrameHandler             { return nil }
func (c *testClient) GetPermissionHandler() cef.PermissionHandler   { return nil }
func (c *testClient) GetJsdialogHandler() cef.JsdialogHandler       { return nil }
func (c *testClient) GetKeyboardHandler() cef.KeyboardHandler       { return nil }
func (c *testClient) GetLifeSpanHandler() cef.LifeSpanHandler       { return c.lifeSpanHandler }
func (c *testClient) GetLoadHandler() cef.LoadHandler               { return nil }
func (c *testClient) GetPrintHandler() cef.PrintHandler             { return nil }
func (c *testClient) GetRenderHandler() cef.RenderHandler           { return c.renderHandler }
func (c *testClient) GetRequestHandler() cef.RequestHandler         { return nil }
func (c *testClient) OnProcessMessageReceived(browser cef.Browser, frame cef.Frame, sourceProcess cef.ProcessID, message cef.ProcessMessage) int32 {
	return 0
}

// testRenderHandler is a minimal RenderHandler that signals paint callbacks.
type testRenderHandler struct {
	painted     chan struct{}
	accelerated chan struct{}
}

var (
	_ cef.Client        = (*testClient)(nil)
	_ cef.RenderHandler = (*testRenderHandler)(nil)
)

func (h *testRenderHandler) GetAccessibilityHandler() cef.AccessibilityHandler { return nil }
func (h *testRenderHandler) GetRootScreenRect(browser cef.Browser, rect *cef.Rect) int32 {
	return 0
}
func (h *testRenderHandler) GetViewRect(browser cef.Browser, rect *cef.Rect) {
	if rect == nil {
		return
	}
	rect.Width = 800
	rect.Height = 600
}
func (h *testRenderHandler) GetScreenPoint(browser cef.Browser, viewx, viewy int32, screenx, screeny *int32) int32 {
	return 0
}
func (h *testRenderHandler) GetScreenInfo(browser cef.Browser, screenInfo *cef.ScreenInfo) int32 {
	return 0
}
func (h *testRenderHandler) OnPopupShow(browser cef.Browser, show int32)     {}
func (h *testRenderHandler) OnPopupSize(browser cef.Browser, rect *cef.Rect) {}
func (h *testRenderHandler) OnPaint(browser cef.Browser, type_ cef.PaintElementType, dirtyrects []cef.Rect, buffer []byte, width, height int32) {
	select {
	case h.painted <- struct{}{}:
	default:
	}
}
func (h *testRenderHandler) OnAcceleratedPaint(browser cef.Browser, type_ cef.PaintElementType, dirtyrects []cef.Rect, info *cef.AcceleratedPaintInfo) {
	select {
	case h.accelerated <- struct{}{}:
	default:
	}
}
func (h *testRenderHandler) GetTouchHandleSize(browser cef.Browser, orientation cef.HorizontalAlignment, size *cef.Size) {
}
func (h *testRenderHandler) OnTouchHandleStateChanged(browser cef.Browser, state *cef.TouchHandleState) {
}
func (h *testRenderHandler) StartDragging(browser cef.Browser, dragData cef.DragData, allowedOps cef.DragOperationsMask, x, y int32) int32 {
	return 0
}
func (h *testRenderHandler) UpdateDragCursor(browser cef.Browser, operation cef.DragOperationsMask) {
}
func (h *testRenderHandler) OnScrollOffsetChanged(browser cef.Browser, x, y float64) {}
func (h *testRenderHandler) OnImeCompositionRangeChanged(browser cef.Browser, selectedRange *cef.Range, characterBounds []cef.Rect) {
}
func (h *testRenderHandler) OnTextSelectionChanged(browser cef.Browser, selectedText string, selectedRange *cef.Range) {
}
func (h *testRenderHandler) OnVirtualKeyboardRequested(browser cef.Browser, inputMode cef.TextInputMode) {
}

// TestAPICompiles verifies runtime-free settings behavior. Interface assertions
// above protect the generated handler signatures without initializing CEF.
func TestAPICompiles(t *testing.T) {
	settings := cef.DefaultSettings()
	if !settings.WindowlessRenderingEnabled {
		t.Fatal("expected WindowlessRenderingEnabled to be true")
	}

	t.Log("API compile check passed")
}
