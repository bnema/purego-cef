package cef

import "testing"

func TestRuntimeInterfaceSatisfied(t *testing.T) {
	var _ Runtime = (*runtime)(nil)
}

func TestBrowserInterfaceSatisfied(t *testing.T) {
	var _ Browser = (*browser)(nil)
}

func TestNewRuntimeReturnsNonNil(t *testing.T) {
	rt := NewRuntime("")
	if rt == nil {
		t.Fatal("NewRuntime returned nil")
	}
}

func TestRuntimeCreateBrowserBeforeInitFails(t *testing.T) {
	rt := NewRuntime("")
	_, err := rt.CreateBrowser(BrowserConfig{})
	if err == nil {
		t.Fatal("expected error when calling CreateBrowser before Init")
	}
}
