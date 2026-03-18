package cef

import "testing"

func TestClientBuilderStoresHandlers(t *testing.T) {
	render := RenderHandlerFunc(func(PaintEvent) {})
	client := NewClient().WithRenderHandler(render)
	if client.renderHandler == nil {
		t.Fatal("render handler not stored")
	}
}

func TestClientBuilderStoresLifeSpanHandler(t *testing.T) {
	client := NewClient().WithLifeSpanHandler(nil)
	if client.lifeSpanHandler != nil {
		t.Fatal("expected nil life span handler")
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
