#!/usr/bin/env sh
set -eu

if [ -z "${CEF_HEADERS:-}" ]; then
  echo "ERROR: CEF_HEADERS must name the CEF 150 include directory (for example /usr/include/cef/include)." >&2
  exit 1
fi
headers_dir="$CEF_HEADERS"

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

