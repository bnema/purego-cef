package mocks_test

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/bnema/purego-cef/cef"
	cefmocks "github.com/bnema/purego-cef/cef/mocks"
)

var (
	_ cef.App                   = (*cefmocks.MockApp)(nil)
	_ cef.Browser               = (*cefmocks.MockBrowser)(nil)
	_ cef.BrowserProcessHandler = (*cefmocks.MockBrowserProcessHandler)(nil)
	_ cef.Client                = (*cefmocks.MockClient)(nil)
	_ cef.ContextMenuHandler    = (*cefmocks.MockContextMenuHandler)(nil)
	_ cef.DialogHandler         = (*cefmocks.MockDialogHandler)(nil)
	_ cef.DisplayHandler        = (*cefmocks.MockDisplayHandler)(nil)
	_ cef.FocusHandler          = (*cefmocks.MockFocusHandler)(nil)
	_ cef.Frame                 = (*cefmocks.MockFrame)(nil)
	_ cef.KeyboardHandler       = (*cefmocks.MockKeyboardHandler)(nil)
	_ cef.LifeSpanHandler       = (*cefmocks.MockLifeSpanHandler)(nil)
	_ cef.LoadHandler           = (*cefmocks.MockLoadHandler)(nil)
	_ cef.ProcessMessage        = (*cefmocks.MockProcessMessage)(nil)
	_ cef.RenderHandler         = (*cefmocks.MockRenderHandler)(nil)
	_ cef.RenderProcessHandler  = (*cefmocks.MockRenderProcessHandler)(nil)
	_ cef.Request               = (*cefmocks.MockRequest)(nil)
	_ cef.RequestContext        = (*cefmocks.MockRequestContext)(nil)
	_ cef.RequestHandler        = (*cefmocks.MockRequestHandler)(nil)
	_ cef.Response              = (*cefmocks.MockResponse)(nil)
)

func TestCuratedMocksOnly(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	dir := filepath.Dir(thisFile)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", dir, err)
	}

	var got []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if filepath.Ext(name) != ".go" || hasTestSuffix(name) {
			continue
		}
		got = append(got, name)
	}
	slices.Sort(got)

	want := []string{
		"mock_app.go",
		"mock_browser.go",
		"mock_browser_process_handler.go",
		"mock_client.go",
		"mock_context_menu_handler.go",
		"mock_dialog_handler.go",
		"mock_display_handler.go",
		"mock_focus_handler.go",
		"mock_frame.go",
		"mock_keyboard_handler.go",
		"mock_life_span_handler.go",
		"mock_load_handler.go",
		"mock_process_message.go",
		"mock_render_handler.go",
		"mock_render_process_handler.go",
		"mock_request.go",
		"mock_request_context.go",
		"mock_request_handler.go",
		"mock_response.go",
	}

	if !slices.Equal(got, want) {
		t.Fatalf("cef/mocks files = %v, want %v", got, want)
	}
}

func hasTestSuffix(name string) bool {
	return len(name) > len("_test.go") && name[len(name)-len("_test.go"):] == "_test.go"
}
