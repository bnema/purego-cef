//go:build cef_integration

package integration

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/bnema/purego-cef/cef"
)

// testClient is a minimal Client implementation that returns nil for all
// sub-handlers except the render handler.
type testClient struct {
	renderHandler cef.RenderHandler
}

func (c *testClient) GetAudioHandler() cef.AudioHandler                { return nil }
func (c *testClient) GetCommandHandler() cef.CommandHandler            { return nil }
func (c *testClient) GetContextMenuHandler() cef.ContextMenuHandler    { return nil }
func (c *testClient) GetDialogHandler() cef.DialogHandler              { return nil }
func (c *testClient) GetDisplayHandler() cef.DisplayHandler            { return nil }
func (c *testClient) GetDownloadHandler() cef.DownloadHandler          { return nil }
func (c *testClient) GetDragHandler() cef.DragHandler                  { return nil }
func (c *testClient) GetFindHandler() cef.FindHandler                  { return nil }
func (c *testClient) GetFocusHandler() cef.FocusHandler                { return nil }
func (c *testClient) GetFrameHandler() cef.FrameHandler                { return nil }
func (c *testClient) GetPermissionHandler() cef.PermissionHandler      { return nil }
func (c *testClient) GetJsdialogHandler() cef.JsdialogHandler          { return nil }
func (c *testClient) GetKeyboardHandler() cef.KeyboardHandler          { return nil }
func (c *testClient) GetLifeSpanHandler() cef.LifeSpanHandler          { return nil }
func (c *testClient) GetLoadHandler() cef.LoadHandler                  { return nil }
func (c *testClient) GetPrintHandler() cef.PrintHandler                { return nil }
func (c *testClient) GetRenderHandler() cef.RenderHandler              { return c.renderHandler }
func (c *testClient) GetRequestHandler() cef.RequestHandler            { return nil }
func (c *testClient) OnProcessMessageReceived(browser cef.Browser, frame cef.Frame, sourceProcess cef.ProcessID, message cef.ProcessMessage) int32 {
	return 0
}

// testRenderHandler is a minimal RenderHandler that signals when OnPaint fires.
type testRenderHandler struct {
	painted chan struct{}
}

func (h *testRenderHandler) GetAccessibilityHandler() cef.AccessibilityHandler { return nil }
func (h *testRenderHandler) GetRootScreenRect(browser cef.Browser, rect uintptr) int32 {
	return 0
}
func (h *testRenderHandler) GetViewRect(browser cef.Browser, rect uintptr) {}
func (h *testRenderHandler) GetScreenPoint(browser cef.Browser, viewx, viewy int32, screenx, screeny unsafe.Pointer) int32 {
	return 0
}
func (h *testRenderHandler) GetScreenInfo(browser cef.Browser, screenInfo uintptr) int32 { return 0 }
func (h *testRenderHandler) OnPopupShow(browser cef.Browser, show int32)                 {}
func (h *testRenderHandler) OnPopupSize(browser cef.Browser, rect uintptr)                {}
func (h *testRenderHandler) OnPaint(browser cef.Browser, type_ cef.PaintElementType, dirtyrectscount int, dirtyrects uintptr, buffer unsafe.Pointer, width, height int32) {
	select {
	case h.painted <- struct{}{}:
	default:
	}
}
func (h *testRenderHandler) OnAcceleratedPaint(browser cef.Browser, type_ cef.PaintElementType, dirtyrectscount int, dirtyrects uintptr, info uintptr) {
}
func (h *testRenderHandler) GetTouchHandleSize(browser cef.Browser, orientation cef.HorizontalAlignment, size uintptr) {
}
func (h *testRenderHandler) OnTouchHandleStateChanged(browser cef.Browser, state uintptr)        {}
func (h *testRenderHandler) StartDragging(browser cef.Browser, dragData cef.DragData, allowedOps cef.DragOperationsMask, x, y int32) int32 {
	return 0
}
func (h *testRenderHandler) UpdateDragCursor(browser cef.Browser, operation cef.DragOperationsMask) {
}
func (h *testRenderHandler) OnScrollOffsetChanged(browser cef.Browser, x, y float64) {}
func (h *testRenderHandler) OnImeCompositionRangeChanged(browser cef.Browser, selectedRange uintptr, characterBoundscount int, characterBounds uintptr) {
}
func (h *testRenderHandler) OnTextSelectionChanged(browser cef.Browser, selectedText string, selectedRange uintptr) {
}
func (h *testRenderHandler) OnVirtualKeyboardRequested(browser cef.Browser, inputMode cef.TextInputMode) {
}

// TestAPICompiles verifies that the generated public API types compose
// correctly. This is a compile-time check; the actual CEF runtime is
// required for a functional test (run with -tags cef_integration on a
// machine that has CEF installed).
func TestAPICompiles(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	handler := &testRenderHandler{painted: make(chan struct{}, 1)}
	client := &testClient{renderHandler: handler}

	// Verify NewClient accepts our implementation.
	clientPtr := cef.NewClient(client)
	if clientPtr == nil {
		t.Fatal("NewClient returned nil")
	}

	// Verify NewRenderHandler accepts our implementation.
	rhPtr := cef.NewRenderHandler(handler)
	if rhPtr == nil {
		t.Fatal("NewRenderHandler returned nil")
	}

	// Verify settings conversion compiles.
	settings := cef.DefaultSettings()
	if !settings.WindowlessRenderingEnabled {
		t.Fatal("expected WindowlessRenderingEnabled to be true")
	}

	t.Log("API compile check passed")
}
