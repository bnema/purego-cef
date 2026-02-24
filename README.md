# purego-cef

CEF (Chromium Embedded Framework) bindings for Go — no CGo, pure `purego`.

> Provides a Chromium-based webview (CEF OSR) embedded in a GTK4 window via purego — no CGo.

## Status

Bootstrap phase. See `docs/architecture.md` for the design.

## Requirements

- Go 1.26+
- CEF runtime bundle (libcef.so + resources) in `~/.local/share/cef` or `CEF_PATH`

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
  gtk4osr/           # GTK4 DrawingArea adapter + event translation
  runtimeassets/     # CEF runtime path discovery and validation
  examples/
    gtk4-osr-minimal/
```

## Related projects

- [`bnema/purego`](https://github.com/bnema/purego) — unified purego fork (base runtime)
- [`bnema/puregotk`](https://github.com/bnema/puregotk) — GTK4 bindings
- [`bnema/dumber`](https://github.com/bnema/dumber) — the app using purego-cef as its webview backend

## License

MIT
