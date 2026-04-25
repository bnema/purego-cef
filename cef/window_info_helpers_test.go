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

func TestSetAsWindowlessWithoutSharedTexture(t *testing.T) {
	info := NewWindowInfo()

	SetAsWindowless(&info, WindowHandle(42), false)

	if info.ParentWindow != WindowHandle(42) {
		t.Fatalf("expected parent window 42, got %v", info.ParentWindow)
	}
	if info.WindowlessRenderingEnabled != 1 {
		t.Fatalf("expected windowless rendering enabled, got %d", info.WindowlessRenderingEnabled)
	}
	if info.SharedTextureEnabled != 0 {
		t.Fatalf("expected shared texture disabled, got %d", info.SharedTextureEnabled)
	}
}

func TestSetAsWindowlessNilWindowInfoDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected no panic, got %v", r)
		}
	}()

	SetAsWindowless(nil, WindowHandle(1), true)
}
