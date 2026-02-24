# purego-cef

CEF (Chromium Embedded Framework) bindings for Go — no CGo, pure `purego`.

## Status

Bootstrap phase.

## Requirements

- Go 1.26+
- CEF runtime bundle (`libcef.so` + resources) in `~/.local/share/cef` or `$CEF_PATH`

## Architecture

```
purego-cef/
  cmd/cefgen/        # code generator: CEF C headers → Go bindings
  internal/
    gen/             # parser, policy, templates
    loader/          # dlopen/dlsym + symbol resolution
    capi/            # generated raw CEF C API bindings (version-pinned)
    runtime/         # refcount, callback registry, thread utilities
    versioning/      # api hash/version validation
  cef/               # public safe API
  examples/
    minimal/
```

## License

MIT
