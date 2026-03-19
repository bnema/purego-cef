package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/emitter"
	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
	"github.com/bnema/purego-cef/cmd/cefgen/internal/parser"
)

type config struct {
	headersDir string
	rawDir     string // cef/internal/raw/
	publicDir  string // cef/
	version    string
}

func (c config) validate() error {
	if c.headersDir == "" || c.rawDir == "" || c.publicDir == "" {
		return fmt.Errorf("--headers-dir, --raw-dir, and --public-dir are required")
	}
	return nil
}

func main() {
	var cfg config
	flag.StringVar(&cfg.headersDir, "headers-dir", "", "CEF include root")
	flag.StringVar(&cfg.rawDir, "raw-dir", "", "raw struct output directory")
	flag.StringVar(&cfg.publicDir, "public-dir", "", "public API output directory")
	flag.StringVar(&cfg.version, "version", "145", "target major version")
	flag.Parse()
	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// headerEntry holds a parsed header and its derived output filename.
type headerEntry struct {
	header  *model.Header
	outName string
}

func run(cfg config) error {
	if err := cfg.validate(); err != nil {
		return err
	}

	// Ensure output directories exist.
	if err := os.MkdirAll(cfg.rawDir, 0o755); err != nil {
		return fmt.Errorf("create raw dir: %w", err)
	}
	if err := os.MkdirAll(cfg.publicDir, 0o755); err != nil {
		return fmt.Errorf("create public dir: %w", err)
	}

	var allHeaders []*model.Header
	var entries []headerEntry
	var registerNames []string

	// Pass 1: Type headers (include/internal/cef_types*.h)
	typeHeaders, err := filepath.Glob(filepath.Join(cfg.headersDir, "internal", "cef_types*.h"))
	if err != nil {
		return err
	}
	typeHeaders = filterOut(typeHeaders, "cef_types_wrappers.h")

	for _, path := range typeHeaders {
		header, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		base := filepath.Base(path)
		header.RegisterName = registerName(base)
		registerNames = append(registerNames, header.RegisterName)
		outName := outputName(base)
		allHeaders = append(allHeaders, header)
		entries = append(entries, headerEntry{header: header, outName: outName})
	}

	// Pass 2: CAPI headers (include/capi/**/*_capi.h)
	capiDirs := []string{
		filepath.Join(cfg.headersDir, "capi"),
		filepath.Join(cfg.headersDir, "capi", "views"),
		filepath.Join(cfg.headersDir, "capi", "test"),
	}
	for _, dir := range capiDirs {
		matches, err := filepath.Glob(filepath.Join(dir, "*_capi.h"))
		if err != nil {
			return err
		}
		for _, path := range matches {
			header, err := parser.ParseFile(path)
			if err != nil {
				return fmt.Errorf("parse %s: %w", path, err)
			}
			rel, _ := filepath.Rel(filepath.Join(cfg.headersDir, "capi"), path)
			if rel == "" {
				rel = filepath.Base(path)
			}
			header.RegisterName = registerName(filepath.Base(path))
			registerNames = append(registerNames, header.RegisterName)
			outName := outputName(rel)
			allHeaders = append(allHeaders, header)
			entries = append(entries, headerEntry{header: header, outName: outName})
		}
	}

	// Build type registry from all parsed headers.
	registry := emitter.NewTypeRegistry(allHeaders)

	// Emit both raw and public output for each header.
	for _, e := range entries {
		// Raw output
		rawCode, err := emitter.EmitRaw(e.header)
		if err != nil {
			return fmt.Errorf("emit raw %s: %w", e.outName, err)
		}
		if err := os.WriteFile(filepath.Join(cfg.rawDir, e.outName), []byte(rawCode), 0o644); err != nil {
			return fmt.Errorf("write raw %s: %w", e.outName, err)
		}

		// Public output
		pubData := emitter.BuildPublicFileData(e.header, registry)
		pubCode, err := emitter.EmitPublic(pubData)
		if err != nil {
			return fmt.Errorf("emit public %s: %w", e.outName, err)
		}
		if err := os.WriteFile(filepath.Join(cfg.publicDir, e.outName), []byte(pubCode), 0o644); err != nil {
			return fmt.Errorf("write public %s: %w", e.outName, err)
		}
	}

	// Generate register.go aggregator for raw directory.
	return writeRawRegisterAggregator(cfg.rawDir, registerNames)
}

// registerName derives a unique Go register function name from a header
// filename, e.g. "cef_app_capi.h" -> "RegisterApp",
// "cef_types.h" -> "RegisterTypes".
func registerName(base string) string {
	name := base
	// Strip suffix: try _capi.h first, then .h
	if strings.HasSuffix(name, "_capi.h") {
		name = strings.TrimSuffix(name, "_capi.h")
	} else {
		name = strings.TrimSuffix(name, ".h")
	}
	// Strip cef_ prefix
	name = strings.TrimPrefix(name, "cef_")
	// Title-case each word
	parts := strings.Split(name, "_")
	var sb strings.Builder
	sb.WriteString("Register")
	for _, p := range parts {
		if p == "" {
			continue
		}
		sb.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	return sb.String()
}

// outputName derives the Go output filename from a header path.
// Includes directory prefix for views/ and test/ to avoid collisions.
func outputName(relPath string) string {
	base := filepath.Base(relPath)
	name := strings.TrimSuffix(base, "_capi.h")
	if name == base {
		name = strings.TrimSuffix(base, ".h")
	}

	// Prefix with directory for views/ and test/ headers
	dir := filepath.Dir(relPath)
	switch {
	case strings.HasSuffix(dir, "views"):
		name = "views_" + name
	case strings.HasSuffix(dir, "test"):
		name = "test_" + name
	}

	// Avoid generating files ending in _test.go (Go treats those as test files)
	if strings.HasSuffix(name, "_test") {
		name = name + "_"
	}

	return name + ".go"
}

// filterOut removes paths ending with the given suffix.
func filterOut(paths []string, suffix string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if !strings.HasSuffix(p, suffix) {
			out = append(out, p)
		}
	}
	return out
}

// writeRawRegisterAggregator generates register.go for the raw package that
// calls all per-header register functions.
func writeRawRegisterAggregator(dir string, names []string) error {
	var sb strings.Builder
	sb.WriteString("package raw\n\n")
	sb.WriteString("// Register loads all CEF C API symbols from the shared library.\n")
	sb.WriteString("// Code generated by cefgen. DO NOT EDIT.\n")
	sb.WriteString("func Register(handle uintptr) {\n")
	for _, name := range names {
		sb.WriteString("\t" + name + "(handle)\n")
	}
	sb.WriteString("}\n")
	return os.WriteFile(filepath.Join(dir, "register.go"), []byte(sb.String()), 0o644)
}
