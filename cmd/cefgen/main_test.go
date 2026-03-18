package main

import "testing"

func TestConfigValidate(t *testing.T) {
	cfg := config{headersDir: "/tmp/headers", outputDir: "/tmp/out", version: "145"}
	if err := cfg.validate(); err != nil {
		t.Fatal(err)
	}
}

func TestConfigValidateErrors(t *testing.T) {
	tests := []struct {
		name string
		cfg  config
	}{
		{"empty headersDir", config{headersDir: "", outputDir: "/tmp/out", version: "145"}},
		{"empty outputDir", config{headersDir: "/tmp/headers", outputDir: "", version: "145"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.validate(); err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}
