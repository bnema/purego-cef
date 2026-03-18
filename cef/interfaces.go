package cef

// Runtime manages the CEF lifecycle. All methods must be called from the main
// OS thread (the one that called runtime.LockOSThread).
type Runtime interface {
	// Init opens the CEF shared library, registers bindings, and calls
	// cef_initialize. Must be called before CreateBrowser.
	Init(settings Settings) error

	// Shutdown releases all CEF resources. Call before process exit.
	Shutdown()

	// DoMessageLoopWork performs a single iteration of CEF message loop work.
	DoMessageLoopWork()

	// CreateBrowser creates an off-screen browser. Pumps the message loop
	// internally until OnAfterCreated fires or a timeout is reached.
	CreateBrowser(cfg BrowserConfig) (Browser, error)
}

// Browser represents a single CEF browser instance.
type Browser interface {
	// Close asks the browser host to close.
	Close()
}
