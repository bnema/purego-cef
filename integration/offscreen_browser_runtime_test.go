//go:build cef_integration

package integration

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/bnema/purego-cef/cef"
)

const (
	browserCreationTimeout  = 30 * time.Second
	acceleratedPaintTimeout = 60 * time.Second
	shutdownTimeout         = 10 * time.Second
)

// runtimeLifeSpanHandler reports browser creation and final close. It lets the
// runtime test drive CEF's external message pump from its locked UI thread.
type runtimeLifeSpanHandler struct {
	created chan cef.Browser
	closed  chan struct{}
}

func (h *runtimeLifeSpanHandler) OnBeforePopup(cef.Browser, cef.Frame, int32, string, string, cef.WindowOpenDisposition, int32, *cef.PopupFeatures, *cef.WindowInfo, *cef.RawClientWriteSlot, *cef.BrowserSettings, *cef.DictionaryValue, *bool) bool {
	return true
}
func (h *runtimeLifeSpanHandler) OnBeforePopupAborted(cef.Browser, int32) {}
func (h *runtimeLifeSpanHandler) OnBeforeDevToolsPopup(cef.Browser, *cef.WindowInfo, *cef.RawClientWriteSlot, *cef.BrowserSettings, *cef.DictionaryValue, *bool) {
}
func (h *runtimeLifeSpanHandler) OnAfterCreated(browser cef.Browser) {
	select {
	case h.created <- browser:
	default:
	}
}
func (h *runtimeLifeSpanHandler) DoClose(cef.Browser) bool { return false }
func (h *runtimeLifeSpanHandler) OnBeforeClose(cef.Browser) {
	select {
	case h.closed <- struct{}{}:
	default:
	}
}

var _ cef.LifeSpanHandler = (*runtimeLifeSpanHandler)(nil)

// TestAcceleratedSharedTextureOSRAndShutdown is deliberately opt-in: it needs
// a provisioned CEF GPU runtime and must run on CEF's UI thread.
func TestAcceleratedSharedTextureOSRAndShutdown(t *testing.T) {
	if os.Getenv("CEF_RUNTIME_INTEGRATION") != "1" {
		t.Skip("set CEF_RUNTIME_INTEGRATION=1 on a provisioned CEF GPU runner")
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	settings := cef.DefaultSettings()
	settings.CEFDir = os.Getenv("CEF_DIR")
	subprocessPath, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}
	settings.BrowserSubprocessPath = subprocessPath
	settings.LogFile = runtimeLogFile(t)
	if err := cef.InitWithApp(settings, newIntegrationApp()); err != nil {
		t.Fatalf("cef.Init: %v", err)
	}
	defer cef.Shutdown()

	render := &testRenderHandler{
		painted:     make(chan struct{}, 1),
		accelerated: make(chan struct{}, 1),
	}
	lifeSpan := &runtimeLifeSpanHandler{
		created: make(chan cef.Browser, 1),
		closed:  make(chan struct{}, 1),
	}
	client := cef.NewClient(&testClient{
		renderHandler:   render,
		lifeSpanHandler: lifeSpan,
	})
	if client == nil {
		t.Fatal("cef.NewClient returned nil")
	}

	windowInfo := cef.NewWindowInfo()
	cef.SetAsWindowless(&windowInfo, 0, true)
	browserSettings := cef.NewBrowserSettings()
	if got := cef.BrowserHostCreateBrowser(&windowInfo, client, "data:text/html,<title>osr</title><body>accelerated OSR</body>", &browserSettings, nil, nil); got != 1 {
		t.Fatalf("cef.BrowserHostCreateBrowser = %d, want 1", got)
	}

	var browser cef.Browser
	pumpUntil(t, browserCreationTimeout, "browser creation", func() bool {
		select {
		case browser = <-lifeSpan.created:
			return browser != nil
		default:
			return false
		}
	})

	host := browser.GetHost()
	if host == nil {
		t.Fatal("created browser has no host")
	}
	host.WasResized()
	host.Invalidate(cef.PaintElementTypePetView)

	pumpUntil(t, acceleratedPaintTimeout, "accelerated shared-texture paint", func() bool {
		select {
		case <-render.accelerated:
			return true
		default:
			return false
		}
	})

	host.CloseBrowser(1)
	pumpUntil(t, shutdownTimeout, "browser close", func() bool {
		select {
		case <-lifeSpan.closed:
			return true
		default:
			return false
		}
	})

	cef.Shutdown()
}

func pumpUntil(t *testing.T, timeout time.Duration, event string, done func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for !done() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out after %s waiting for %s", timeout, event)
		}
		cef.DoMessageLoopWork()
		time.Sleep(10 * time.Millisecond)
	}
}

func runtimeLogFile(t *testing.T) string {
	t.Helper()
	dir := os.Getenv("CEF_RUNTIME_LOG_DIR")
	if dir == "" {
		dir = t.TempDir()
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create CEF log directory: %v", err)
	}
	return filepath.Join(dir, "cef-runtime.log")
}
