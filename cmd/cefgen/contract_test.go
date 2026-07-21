package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bnema/purego-cef/internal/cefapi"
)

func TestReadAPIContractRequiresTargetLinuxHash(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cef_api_versions.h")
	header := fmt.Sprintf(`#define CEF_API_VERSION_%[1]d %[1]d
#if defined(OS_LINUX)
#define CEF_API_HASH_%[1]d "0123456789abcdef"
#endif
`, cefapi.Version)
	if err := os.WriteFile(path, []byte(header), 0o644); err != nil {
		t.Fatal(err)
	}
	hash, err := readAPIContract(path)
	if err != nil {
		t.Fatal(err)
	}
	if hash != "0123456789abcdef" {
		t.Fatalf("hash = %q", hash)
	}

	versionOnly := fmt.Sprintf(`#define CEF_API_VERSION_%[1]d %[1]d`, cefapi.Version)
	if err := os.WriteFile(path, []byte(versionOnly), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readAPIContract(path); err == nil {
		t.Fatalf("missing Linux hash for API %d accepted", cefapi.Version)
	}

	if err := os.WriteFile(path, []byte(`#define CEF_API_VERSION_14900 14900`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readAPIContract(path); err == nil {
		t.Fatalf("missing API %d accepted", cefapi.Version)
	}
}
