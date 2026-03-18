//go:build cef_integration

package integration

import (
	"runtime"
	"testing"
	"time"

	"github.com/bnema/purego-cef/cef"
)

type testRenderHandler struct {
	painted chan struct{}
}

func (h *testRenderHandler) GetViewRect() cef.Rect {
	return cef.Rect{W: 800, H: 600}
}

func (h *testRenderHandler) OnPaint(event cef.PaintEvent) {
	select {
	case h.painted <- struct{}{}:
	default:
	}
}

func TestOffscreenBrowserPaints(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	rt := cef.NewRuntime("")
	if err := rt.Init(cef.DefaultSettings()); err != nil {
		t.Fatal(err)
	}
	defer rt.Shutdown()

	handler := &testRenderHandler{painted: make(chan struct{}, 1)}
	client := cef.NewClient().WithRenderHandler(handler)

	browser, err := rt.CreateBrowser(cef.BrowserConfig{
		URL:       "data:text/html,<html><body>purego-cef</body></html>",
		Width:     800,
		Height:    600,
		FrameRate: 30,
		Client:    client,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer browser.Close()

	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-handler.painted:
			return
		case <-timeout:
			t.Fatal("timed out waiting for first paint")
		default:
			rt.DoMessageLoopWork()
			time.Sleep(10 * time.Millisecond)
		}
	}
}
