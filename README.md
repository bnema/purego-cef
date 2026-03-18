# purego-cef

Pure Go bindings for the Chromium Embedded Framework C API.

## Status

Milestone 1 targets a minimal offscreen browser on Linux with `CGO_ENABLED=0`.

## Requirements

- Go 1.26+
- `github.com/ebitengine/purego` v0.10.0
- CEF runtime bundle in `$CEF_DIR` or `~/.local/share/cef`

## Architecture

```text
cef/              handwritten public API
internal/capi/    generated raw CEF bindings
internal/loader/  libcef discovery + dlopen
internal/cefstr/  cef_string_t helpers
internal/refcount/ Go-owned refcount callbacks
cmd/cefgen/       header parser + code emitter
```

## Development

```bash
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go generate ./internal/capi
```
