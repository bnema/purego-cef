package cef

import "testing"

func TestSetAsWindowless(t *testing.T) {
	info := NewWindowInfo()

	SetAsWindowless(&info, WindowHandle(42), true)

	if info.ParentWindow != WindowHandle(42) {
		t.Fatalf("expected parent window 42, got %v", info.ParentWindow)
	}
	if info.WindowlessRenderingEnabled != 1 {
		t.Fatalf("expected windowless rendering enabled, got %d", info.WindowlessRenderingEnabled)
	}
	if info.SharedTextureEnabled != 1 {
		t.Fatalf("expected shared texture enabled, got %d", info.SharedTextureEnabled)
	}
}
