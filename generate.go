// NOTE: The sh -c invocation with shell variable expansion requires
// a POSIX shell (Linux/macOS). Windows users should set CEF_HEADERS
// explicitly or run the generator manually.
//go:generate sh -c "go run ./cmd/cefgen --headers-dir ${CEF_HEADERS:-$HOME/.local/share/cef/include} --capi-dir internal/capi --port-in-dir internal/ports/in --port-out-dir internal/ports/out --public-dir cef"
//go:generate mockery

package puregocef
