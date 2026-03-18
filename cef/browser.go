package cef

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/cefstr"
)

// Browser wraps a CEF browser instance.
type Browser struct {
	raw   *capi.CEFBrowserT
	state *clientState
}

// BrowserConfig configures browser creation.
type BrowserConfig struct {
	URL       string
	Width     int32
	Height    int32
	FrameRate int32
	Client    *Client
}

// CreateBrowser creates an off-screen browser asynchronously and pumps
// the message loop until OnAfterCreated fires or a timeout is reached.
func CreateBrowser(cfg BrowserConfig) (*Browser, error) {
	if cfg.Client == nil {
		cfg.Client = NewClient()
	}
	if cfg.Width == 0 {
		cfg.Width = 800
	}
	if cfg.Height == 0 {
		cfg.Height = 600
	}
	if cfg.FrameRate == 0 {
		cfg.FrameRate = 30
	}
	if cfg.URL == "" {
		cfg.URL = "about:blank"
	}

	state := newClientState(cfg.Client)

	var window capi.CEFWindowInfoT
	window.Size = uintptr(unsafe.Sizeof(window))
	window.Bounds = capi.CEFRectT{W: cfg.Width, H: cfg.Height}
	window.WindowlessRenderingEnabled = 1

	var settings capi.CEFBrowserSettingsT
	settings.Size = uintptr(unsafe.Sizeof(settings))
	settings.WindowlessFrameRate = cfg.FrameRate

	url, cleanup, err := cefstr.FromGo(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("cef: create browser url string: %w", err)
	}
	defer cleanup()

	urlCapi := capi.CEFStringT{Str: url.Str, Length: url.Length, Dtor: url.Dtor}

	ok := capi.CEFBrowserHostCreateBrowser(
		unsafe.Pointer(&window),
		unsafe.Pointer(&state.client),
		unsafe.Pointer(&urlCapi),
		unsafe.Pointer(&settings),
		nil, nil,
	)
	if ok != 1 {
		return nil, fmt.Errorf("cef: cef_browser_host_create_browser failed (returned %d)", ok)
	}

	// Pump the message loop until OnAfterCreated delivers the browser.
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	for {
		select {
		case browser := <-state.created:
			browser.state = state
			return browser, nil
		case <-timer.C:
			return nil, fmt.Errorf("cef: timed out waiting for browser creation")
		default:
			capi.CEFDoMessageLoopWork()
			time.Sleep(time.Millisecond)
		}
	}
}

// Close asks the browser host to close.
func (b *Browser) Close() {
	host := (*capi.CEFBrowserHostT)(unsafe.Pointer(b.raw.CallGetHost()))
	host.CallCloseBrowser(1) // force_close = true
}

// Raw returns the underlying CEF browser pointer for advanced use.
func (b *Browser) Raw() *capi.CEFBrowserT { return b.raw }
