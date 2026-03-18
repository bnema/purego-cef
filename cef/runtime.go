package cef

import (
	"fmt"
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/cefstr"
	"github.com/bnema/purego-cef/internal/loader"
)

type runtime struct {
	runtimeDir  string
	handle      uintptr
	initialized bool
}

// NewRuntime creates a new CEF runtime. The runtimeDir parameter specifies
// the CEF runtime directory. An empty string uses the default location
// (~/.local/share/cef or $CEF_DIR).
func NewRuntime(runtimeDir string) Runtime {
	return &runtime{runtimeDir: runtimeDir}
}

func (r *runtime) Init(settings Settings) error {
	if r.initialized {
		return fmt.Errorf("cef: runtime already initialized")
	}

	handle, err := loader.Open(r.runtimeDir)
	if err != nil {
		return fmt.Errorf("cef: open loader: %w", err)
	}
	r.handle = handle

	capi.Register(handle)
	cefstr.Bind(handle)

	args := NewMainArgsFromOS()

	cs, cleanup, err := settings.toC()
	if err != nil {
		return fmt.Errorf("cef: convert settings: %w", err)
	}
	defer cleanup()

	ok := capi.CEFInitialize(args.Ptr(), unsafe.Pointer(&cs), nil, nil)
	if ok != 1 {
		return fmt.Errorf("cef: cef_initialize returned %d", ok)
	}
	r.initialized = true
	return nil
}

func (r *runtime) Shutdown() {
	if !r.initialized {
		return
	}
	capi.CEFShutdown()
	r.initialized = false
}

func (r *runtime) DoMessageLoopWork() {
	capi.CEFDoMessageLoopWork()
}

func (r *runtime) CreateBrowser(cfg BrowserConfig) (Browser, error) {
	if !r.initialized {
		return nil, fmt.Errorf("cef: runtime not initialized, call Init first")
	}
	return createBrowser(cfg)
}
