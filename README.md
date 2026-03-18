# purego-cef

Go bindings for the [Chromium Embedded Framework](https://bitbucket.org/chromiumembedded/cef/) C API, using [purego](https://github.com/ebitengine/purego) instead of cgo.

Linux only. `CGO_ENABLED=0`.

## Usage

```go
package main

import (
	"os"
	"runtime"

	"github.com/bnema/purego-cef/cef"
)

func main() {
	runtime.LockOSThread()

	// Subprocess handling (multi-process mode)
	if code, err := cef.MaybeExitSubprocess(""); err == nil && code >= 0 {
		os.Exit(code)
	}

	rt := cef.NewRuntime("")
	if err := rt.Init(cef.DefaultSettings()); err != nil {
		panic(err)
	}
	defer rt.Shutdown()

	client := cef.NewClient().WithRenderHandler(cef.RenderHandlerFunc(func(e cef.PaintEvent) {
		// e.Buffer contains the BGRA pixel data
	}))

	browser, err := rt.CreateBrowser(cef.BrowserConfig{
		URL:       "https://example.com",
		Width:     1280,
		Height:    720,
		FrameRate: 60,
		Client:    client,
	})
	if err != nil {
		panic(err)
	}
	defer browser.Close()

	// Your render loop here
	for {
		rt.DoMessageLoopWork()
	}
}
```

## Requirements

- Go 1.26+
- CEF 145 runtime in `$CEF_DIR` or `~/.local/share/cef`

## Building

```
CGO_ENABLED=0 go build ./...
CGO_ENABLED=0 go test ./...
```

Integration tests need the CEF runtime:

```
CEF_DIR=$HOME/.local/share/cef go test -tags=cef_integration ./integration
```

## Regenerating bindings

The raw C API bindings in `internal/capi/` are generated from CEF headers:

```
CEF_DIR=$HOME/.local/share/cef go generate ./internal/capi
```

## Project layout

```
cef/               public API (interfaces, runtime, browser)
internal/capi/     generated C API bindings
internal/loader/   libcef.so discovery + dlopen
internal/cefstr/   UTF-16 string helpers
internal/refcount/ prevent GC of Go-owned CEF structs
cmd/cefgen/        binding generator (parses headers, emits Go)
```

## License

MIT
