package capi

// Register loads all CEF C API symbols from the shared library identified by
// handle. It must be called after the CEF library has been opened with
// internal/loader.
func Register(handle uintptr) {
	RegisterBase(handle)
	RegisterApp(handle)
	RegisterClient(handle)
	RegisterLifeSpanHandler(handle)
	RegisterRenderHandler(handle)
	RegisterBrowser(handle)
}
