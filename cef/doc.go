// Package cef provides Go bindings for the Chromium Embedded Framework.
//
// Use [NewRuntime] to create a runtime, [Runtime.Init] to start CEF,
// and [Runtime.CreateBrowser] to open an offscreen browser.
//
// All Runtime and Browser methods must be called from the main OS thread
// (the goroutine that called runtime.LockOSThread).
package cef
