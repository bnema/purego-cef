// Package cef exposes the user-facing purego CEF API.
//
// The package intentionally contains two layers:
//
//   - ergonomic APIs for normal use, such as Settings, Init, NewClient,
//     NewLifeSpanHandler, and NewAudioHandler
//   - low-level generated/raw bindings that mirror the CEF C API more directly
//
// As a rule of thumb:
//
//   - prefer Settings over RawSettings
//   - prefer the safe adapter constructors over raw handler interfaces when
//     they exist (for example NewClient, NewLifeSpanHandler, NewAudioHandler)
//   - treat Raw* aliases and generated raw structs as advanced escape hatches
//
// The generated bindings remain public because they are useful for advanced
// callers and for keeping the binding surface complete, but they should not be
// the default choice for new code.
package cef
