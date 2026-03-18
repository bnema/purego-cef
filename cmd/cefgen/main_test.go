package main

import "testing"

func TestConfigValidate(t *testing.T) {
	cfg := config{headersDir: "/tmp/headers", outputDir: "/tmp/out", version: "145"}
	if err := cfg.validate(); err != nil {
		t.Fatal(err)
	}
}
