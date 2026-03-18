package cef

import "testing"

func TestClientBuilderStoresHandlers(t *testing.T) {
	render := RenderHandlerFunc(func(PaintEvent) {})
	client := NewClient().WithRenderHandler(render)
	if client.renderHandler == nil {
		t.Fatal("render handler not stored")
	}
}

func TestClientBuilderStoresNilLifeSpanHandler(t *testing.T) {
	client := NewClient().WithLifeSpanHandler(nil)
	if client.lifeSpanHandler != nil {
		t.Fatal("expected nil life span handler")
	}
}

type stubLifeSpanHandler struct{}

func (s *stubLifeSpanHandler) OnAfterCreated(Browser) {}
func (s *stubLifeSpanHandler) OnBeforeClose(Browser)  {}

func TestClientBuilderStoresNonNilLifeSpanHandler(t *testing.T) {
	handler := &stubLifeSpanHandler{}
	client := NewClient().WithLifeSpanHandler(handler)
	if client.lifeSpanHandler == nil {
		t.Fatal("life span handler not stored")
	}
}

func TestClientBuilderChaining(t *testing.T) {
	render := RenderHandlerFunc(func(PaintEvent) {})
	client := NewClient().
		WithRenderHandler(render).
		WithLifeSpanHandler(nil)
	if client.renderHandler == nil {
		t.Fatal("render handler lost after chaining")
	}
}
