#!/usr/bin/env sh
set -eu

if [ -n "${CEF_HEADERS:-}" ]; then
  headers_dir="$CEF_HEADERS"
elif [ -n "${HOME:-}" ]; then
  headers_dir="${HOME}/.local/share/cef/include"
else
  echo "ERROR: Neither CEF_HEADERS nor HOME is set. Cannot determine headers directory." >&2
  exit 1
fi

if [ ! -d "$headers_dir" ]; then
  echo "ERROR: headers directory does not exist: $headers_dir" >&2
  exit 1
fi

if [ ! -r "$headers_dir" ]; then
  echo "ERROR: headers directory is not readable: $headers_dir" >&2
  exit 1
fi

go run ./cmd/cefgen \
  --headers-dir "$headers_dir" \
  --capi-dir internal/capi \
  --port-in-dir internal/ports/in \
  --port-out-dir internal/ports/out \
  --public-dir cef

