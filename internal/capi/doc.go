package capi

//go:generate sh -c "go run ../../cmd/cefgen --headers-dir ${DOLLAR}{CEF_DIR:-${DOLLAR}HOME/.local/share/cef}/include --output-dir . --version ${DOLLAR}{CEF_VERSION:-145}"
