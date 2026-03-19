# purego-cef

Go bindings for the [Chromium Embedded Framework](https://bitbucket.org/chromiumembedded/cef/) C API, using [purego](https://github.com/ebitengine/purego) instead of cgo.

Linux only. `CGO_ENABLED=0`.

## Architecture

The entire public API is **generated** from CEF C headers. The generator parses all 114 CEF headers and produces:

- **Go interfaces** for every CEF type (both handler structs you implement and CEF objects you receive)
- **Typed constructors** for handlers (`NewClient(impl Client)`, `NewRenderHandler(impl RenderHandler)`, etc.)
- **Typed wrappers** for CEF objects with automatic refcount management
- **Doc comments** extracted from CEF C headers

A small hand-written support layer (~200 LOC) handles UTF-16 string conversion, refcount wiring, and CEF lifecycle orchestration.

## Usage

```go
package main

import (
	"os"
	"runtime"
	"unsafe"

	"github.com/bnema/purego-cef/cef"
)

func main() {
	runtime.LockOSThread()
	cef.MaybeExitSubprocess()

	if err := cef.Init(cef.DefaultSettings()); err != nil {
		panic(err)
	}
	defer cef.Shutdown()

	client := cef.NewClient(myClient{})

	// Your render loop here
	for {
		cef.DoMessageLoopWork()
	}
}

// Implement the Client interface — only the methods you need.
type myClient struct{}

func (c myClient) GetLifeSpanHandler() cef.LifeSpanHandler { return nil }
func (c myClient) GetRenderHandler() cef.RenderHandler     { return myRenderer{} }
// ... other handler getters return nil
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

The public API in `cef/` and raw bindings in `cef/internal/raw/` are generated from CEF headers:

```
go run ./cmd/cefgen \
  --headers-dir ~/.local/share/cef/include \
  --raw-dir cef/internal/raw \
  --public-dir cef
```

## Project layout

```
cef/                 generated public API (interfaces, wrappers, constructors)
cef/support.go       hand-written: string conversion, refcount wiring
cef/init.go          hand-written: Init(), Shutdown(), orchestration
cef/internal/raw/    generated C struct layouts + purego bindings
cmd/cefgen/          binding generator (parser, model, emitter, templates)
integration/         integration tests (require CEF runtime)
```

## License

MIT
