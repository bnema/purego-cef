package mocks_test

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"github.com/bnema/purego-cef/internal/ports/out"
	outmocks "github.com/bnema/purego-cef/internal/ports/out/mocks"
)

var _ out.CAPI = (*outmocks.MockCAPI)(nil)

func TestCuratedMocksOnly(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	dir := filepath.Dir(thisFile)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%q): %v", dir, err)
	}

	var got []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if filepath.Ext(name) != ".go" || (len(name) > len("_test.go") && name[len(name)-len("_test.go"):] == "_test.go") {
			continue
		}
		got = append(got, name)
	}
	slices.Sort(got)

	want := []string{"mock_capi.go"}
	if !slices.Equal(got, want) {
		t.Fatalf("internal/ports/out/mocks files = %v, want %v", got, want)
	}
}
