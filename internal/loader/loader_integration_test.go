//go:build cef_integration

package loader

import (
	"os"
	"testing"

	"github.com/ebitengine/purego"
)

func TestOpenLoadsLibcef(t *testing.T) {
	cefDir := os.Getenv("CEF_DIR")
	if cefDir == "" {
		t.Skip("CEF_DIR not set")
	}
	handle, err := Open(cefDir)
	if err != nil {
		t.Fatal(err)
	}
	if handle == 0 {
		t.Fatal("handle is zero")
	}
	if err := purego.Dlclose(handle); err != nil {
		t.Fatal(err)
	}
}
