package cef

import "testing"

func TestSetAsWindowless(t *testing.T) {
	tests := []struct {
		name                  string
		sharedTexture         bool
		expectedSharedTexture int32
	}{
		{name: "with shared texture", sharedTexture: true, expectedSharedTexture: 1},
		{name: "without shared texture", sharedTexture: false, expectedSharedTexture: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := NewWindowInfo()

			SetAsWindowless(&info, WindowHandle(42), tt.sharedTexture)

			if info.ParentWindow != WindowHandle(42) {
				t.Fatalf("expected parent window 42, got %v", info.ParentWindow)
			}
			if info.WindowlessRenderingEnabled != 1 {
				t.Fatalf("expected windowless rendering enabled, got %d", info.WindowlessRenderingEnabled)
			}
			if info.SharedTextureEnabled != tt.expectedSharedTexture {
				t.Fatalf("expected shared texture %d, got %d", tt.expectedSharedTexture, info.SharedTextureEnabled)
			}
		})
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
