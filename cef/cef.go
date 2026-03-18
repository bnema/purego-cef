// Package cef provides the public API for embedding the Chromium Embedded
// Framework via purego (no cgo).
package cef

import (
	"fmt"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/cefstr"
	"github.com/bnema/purego-cef/internal/loader"
)

// MaybeExitSubprocess calls cef_execute_process and returns the exit code.
// If the exit code is >= 0, the caller should os.Exit with that code (this
// process was a CEF sub-process). If < 0, continue as the browser process.
//
// This is a process-level function that runs before any Runtime is created.
func MaybeExitSubprocess(runtimeDir string) (int, error) {
	handle, err := loader.Open(runtimeDir)
	if err != nil {
		return -1, fmt.Errorf("cef: open loader: %w", err)
	}

	capi.Register(handle)
	cefstr.Bind(handle)

	args := NewMainArgsFromOS()
	code := capi.CEFExecuteProcess(args.Ptr(), nil, nil)
	return int(code), nil
}
