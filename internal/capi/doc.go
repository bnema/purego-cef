package capi

//go:generate sh -c "go run ../../cmd/cefgen --headers-dir ${CEF_DIR:-$HOME/.local/share/cef}/include --output-dir . --version ${CEF_VERSION:-145}"
