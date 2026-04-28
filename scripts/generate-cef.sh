#!/usr/bin/env sh
set -eu

headers_dir="${CEF_HEADERS:-${HOME}/.local/share/cef/include}"

go run ./cmd/cefgen \
  --headers-dir "$headers_dir" \
  --capi-dir internal/capi \
  --port-in-dir internal/ports/in \
  --port-out-dir internal/ports/out \
  --public-dir cef

