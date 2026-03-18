package capi

import "structs"

// CEFStringT matches cef_string_utf16_t (aka cef_string_t on all platforms).
type CEFStringT struct {
	_      structs.HostLayout
	Str    *uint16
	Length uintptr
	Dtor   uintptr
}

// CEFRectT matches cef_rect_t from cef_types_geometry.h.
type CEFRectT struct {
	_          structs.HostLayout
	X, Y, W, H int32
}

// CEFMainArgsT matches cef_main_args_t from cef_types_linux.h.
type CEFMainArgsT struct {
	_    structs.HostLayout
	Argc int32
	_pad [4]byte // padding to align Argv to 8 bytes on 64-bit
	Argv **byte
}

// CEFSettingsT matches cef_settings_t from cef_types.h (CEF 133+).
// Field order must exactly match the C struct.
type CEFSettingsT struct {
	_                                structs.HostLayout
	Size                             uintptr
	NoSandbox                        int32
	_pad0                            [4]byte // align next cef_string_t to 8
	BrowserSubprocessPath            CEFStringT
	FrameworkDirPath                 CEFStringT
	MainBundlePath                   CEFStringT
	MultiThreadedMessageLoop         int32
	ExternalMessagePump              int32
	WindowlessRenderingEnabled       int32
	CommandLineArgsDisabled          int32
	CachePath                        CEFStringT
	RootCachePath                    CEFStringT
	PersistSessionCookies            int32
	_pad1                            [4]byte // align next cef_string_t
	UserAgent                        CEFStringT
	UserAgentProduct                 CEFStringT
	Locale                           CEFStringT
	LogFile                          CEFStringT
	LogSeverity                      int32 // cef_log_severity_t
	LogItems                         int32 // cef_log_items_t
	JavascriptFlags                  CEFStringT
	ResourcesDirPath                 CEFStringT
	LocalesDirPath                   CEFStringT
	RemoteDebuggingPort              int32
	UncaughtExceptionStackSize       int32
	BackgroundColor                  uint32
	_pad2                            [4]byte // align next cef_string_t
	AcceptLanguageList               CEFStringT
	CookieableSchemesList            CEFStringT
	CookieableSchemesExcludeDefaults int32
	_pad3                            [4]byte // align next cef_string_t
	ChromePolicyID                   CEFStringT
	ChromeAppIconID                  int32
	DisableSignalHandlers            int32
}

// CEFWindowInfoT matches cef_window_info_t from cef_types_linux.h.
type CEFWindowInfoT struct {
	_                          structs.HostLayout
	Size                       uintptr
	WindowName                 CEFStringT
	Bounds                     CEFRectT
	ParentWindow               uintptr // cef_window_handle_t = unsigned long
	WindowlessRenderingEnabled int32
	SharedTextureEnabled       int32
	ExternalBeginFrameEnabled  int32
	_pad1                      [4]byte // align Window to 8
	Window                     uintptr // cef_window_handle_t
	RuntimeStyle               int32   // cef_runtime_style_t
	_pad2                      [4]byte // trailing padding
}

// CEFBrowserSettingsT matches cef_browser_settings_t from cef_types.h.
type CEFBrowserSettingsT struct {
	_                          structs.HostLayout
	Size                       uintptr
	WindowlessFrameRate        int32
	_pad0                      [4]byte // align next cef_string_t
	StandardFontFamily         CEFStringT
	FixedFontFamily            CEFStringT
	SerifFontFamily            CEFStringT
	SansSerifFontFamily        CEFStringT
	CursiveFontFamily          CEFStringT
	FantasyFontFamily          CEFStringT
	DefaultFontSize            int32
	DefaultFixedFontSize       int32
	MinimumFontSize            int32
	MinimumLogicalFontSize     int32
	DefaultEncoding            CEFStringT
	RemoteFonts                int32 // cef_state_t
	Javascript                 int32
	JavascriptCloseWindows     int32
	JavascriptAccessClipboard  int32
	JavascriptDomPaste         int32
	ImageLoading               int32
	ImageShrinkStandaloneToFit int32
	TextAreaResize             int32
	TabToLinks                 int32
	LocalStorage               int32
	Databases                  int32
	Webgl                      int32
	BackgroundColor            uint32
	ChromeStatusBubble         int32
	ChromeZoomBubble           int32
	_pad1                      [4]byte // trailing padding for 8-byte alignment
}

// CEFMouseEventT matches cef_mouse_event_t from cef_types.h.
type CEFMouseEventT struct {
	_         structs.HostLayout
	X         int32
	Y         int32
	Modifiers uint32
}
