//go:build cef_integration

package integration

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
	"unsafe"

	"github.com/bnema/purego-cef/cef"
)

// TestMain lets CEF helper processes run before testing.M parses their
// Chromium-only --type argument. Normal tagged suite invocations never load
// CEF here, so the compile-only workflow remains runtime-free.
func TestMain(m *testing.M) {
	if hasCEFSubprocessType(os.Args[1:]) {
		executed, exitCode, err := cef.ExecuteSubprocessWithApp(newIntegrationApp())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "cef.ExecuteSubprocess: %v\n", err)
			os.Exit(1)
		}
		if !executed {
			_, _ = fmt.Fprintln(os.Stderr, "cef.ExecuteSubprocess did not handle CEF helper process")
			os.Exit(1)
		}
		os.Exit(exitCode)
	}

	os.Exit(m.Run())
}

// integrationApp adds only the no-zygote switch required by CEF when this Go
// test binary is used as its browser subprocess. It is never used by library
// callers or the default CEF settings.
type integrationApp struct{}

func newIntegrationApp() cef.App { return integrationApp{} }

func (integrationApp) OnBeforeCommandLineProcessing(_ string, commandLine cef.CommandLine) {
	commandLine.AppendSwitch("no-zygote")
}
func (integrationApp) OnRegisterCustomSchemes(cef.SchemeRegistrar)         {}
func (integrationApp) GetResourceBundleHandler() cef.ResourceBundleHandler { return nil }
func (integrationApp) GetBrowserProcessHandler() cef.BrowserProcessHandler { return nil }
func (integrationApp) GetRenderProcessHandler() cef.RenderProcessHandler   { return nil }

var _ cef.App = integrationApp{}

func hasCEFSubprocessType(args []string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--type=") {
			return true
		}
	}
	return false
}

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

func TestIntegrationAppDisablesZygote(t *testing.T) {
	commandLine := &switchRecordingCommandLine{}
	newIntegrationApp().OnBeforeCommandLineProcessing("", commandLine)

	if got, want := commandLine.switches, []string{"no-zygote"}; !slices.Equal(got, want) {
		t.Fatalf("integration App switches = %q, want %q", got, want)
	}
}

// switchRecordingCommandLine characterizes the one CommandLine operation used
// by the integration-only App without loading a CEF runtime.
type switchRecordingCommandLine struct {
	switches []string
}

func (*switchRecordingCommandLine) IsValid() bool                        { return true }
func (*switchRecordingCommandLine) IsReadOnly() bool                     { return false }
func (*switchRecordingCommandLine) Copy() cef.CommandLine                { return nil }
func (*switchRecordingCommandLine) InitFromArgv(int32, unsafe.Pointer)   {}
func (*switchRecordingCommandLine) InitFromString(string)                {}
func (*switchRecordingCommandLine) Reset()                               {}
func (*switchRecordingCommandLine) GetArgv(cef.StringList)               {}
func (*switchRecordingCommandLine) GetCommandLineString() string         { return "" }
func (*switchRecordingCommandLine) GetProgram() string                   { return "" }
func (*switchRecordingCommandLine) SetProgram(string)                    {}
func (*switchRecordingCommandLine) HasSwitches() bool                    { return false }
func (*switchRecordingCommandLine) HasSwitch(string) bool                { return false }
func (*switchRecordingCommandLine) GetSwitchValue(string) string         { return "" }
func (*switchRecordingCommandLine) GetSwitches(cef.StringMap)            {}
func (c *switchRecordingCommandLine) AppendSwitch(name string)           { c.switches = append(c.switches, name) }
func (*switchRecordingCommandLine) AppendSwitchWithValue(string, string) {}
func (*switchRecordingCommandLine) HasArguments() bool                   { return false }
func (*switchRecordingCommandLine) GetArguments(cef.StringList)          {}
func (*switchRecordingCommandLine) AppendArgument(string)                {}
func (*switchRecordingCommandLine) PrependWrapper(string)                {}
func (*switchRecordingCommandLine) RemoveSwitch(string)                  {}

var _ cef.CommandLine = (*switchRecordingCommandLine)(nil)

func TestHasCEFSubprocessType(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want bool
	}{
		{name: "renderer", args: []string{"--type=renderer"}, want: true},
		{name: "helper with test flags", args: []string{"-test.run=Test", "--type=utility"}, want: true},
		{name: "browser process", args: []string{"-test.run=TestAPICompiles"}},
		{name: "unrelated type prefix", args: []string{"--types=renderer"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasCEFSubprocessType(tc.args); got != tc.want {
				t.Fatalf("hasCEFSubprocessType(%q) = %t, want %t", tc.args, got, tc.want)
			}
		})
	}
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
