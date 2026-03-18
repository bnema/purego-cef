package cef

import "testing"

func TestMouseEventToC(t *testing.T) {
	event := MouseEvent{X: 10, Y: 20, Modifiers: EventFlagLeftMouseButton}
	raw := event.toC()
	if raw.X != 10 || raw.Y != 20 {
		t.Fatalf("unexpected coordinates: %+v", raw)
	}
	if raw.Modifiers != EventFlagLeftMouseButton {
		t.Fatalf("unexpected modifiers: %d", raw.Modifiers)
	}
}

func TestRenderHandlerFuncDefaultRect(t *testing.T) {
	f := RenderHandlerFunc(func(PaintEvent) {})
	r := f.GetViewRect()
	if r.W != 800 || r.H != 600 {
		t.Fatalf("unexpected default rect: %+v", r)
	}
}
