package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/emitter"
	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
	"github.com/bnema/purego-cef/cmd/cefgen/internal/parser"
)

type config struct {
	headersDir string
	capiDir    string // internal/capi/
	portInDir  string // internal/ports/in/
	portOutDir string // internal/ports/out/
	publicDir  string // cef/
	version    string
}

func (c config) validate() error {
	if c.headersDir == "" || c.capiDir == "" || c.portInDir == "" || c.portOutDir == "" || c.publicDir == "" {
		return fmt.Errorf("--headers-dir, --capi-dir, --port-in-dir, --port-out-dir, and --public-dir are required")
	}
	return nil
}

func main() {
	var cfg config
	flag.StringVar(&cfg.headersDir, "headers-dir", "", "CEF include root")
	flag.StringVar(&cfg.capiDir, "capi-dir", "", "internal/capi/ output directory")
	flag.StringVar(&cfg.portInDir, "port-in-dir", "", "internal/ports/in/ output directory")
	flag.StringVar(&cfg.portOutDir, "port-out-dir", "", "internal/ports/out/ output directory")
	flag.StringVar(&cfg.publicDir, "public-dir", "", "cef/ output directory")
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
	for _, dir := range []string{cfg.capiDir, cfg.portInDir, cfg.portOutDir, cfg.publicDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create dir %s: %w", dir, err)
		}
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

	// Emit capi, port-out, port-in, and public output for each header.
	for _, e := range entries {
		// CAPI output (was "raw")
		rawCode, err := emitter.EmitRaw(e.header)
		if err != nil {
			return fmt.Errorf("emit capi %s: %w", e.outName, err)
		}
		if err := os.WriteFile(filepath.Join(cfg.capiDir, e.outName), []byte(rawCode), 0o644); err != nil {
			return fmt.Errorf("write capi %s: %w", e.outName, err)
		}

		// Build port data (no skip list — ports need the full surface)
		portData := emitter.BuildPortFileData(e.header, registry)

		// Port out output
		portOutCode, err := emitter.EmitPortOut(portData)
		if err != nil {
			return fmt.Errorf("emit port-out %s: %w", e.outName, err)
		}
		if portOutCode != "" {
			if err := os.WriteFile(filepath.Join(cfg.portOutDir, e.outName), []byte(portOutCode), 0o644); err != nil {
				return fmt.Errorf("write port-out %s: %w", e.outName, err)
			}
		}

		// Port in output
		portInCode, err := emitter.EmitPortIn(portData)
		if err != nil {
			return fmt.Errorf("emit port-in %s: %w", e.outName, err)
		}
		if portInCode != "" {
			if err := os.WriteFile(filepath.Join(cfg.portInDir, e.outName), []byte(portInCode), 0o644); err != nil {
				return fmt.Errorf("write port-in %s: %w", e.outName, err)
			}
		}

		// Public facade output (uses filtered data — skips handwritten types)
		pubData := emitter.BuildPublicFileData(e.header, registry)
		pubCode, err := emitter.EmitPublic(pubData)
		if err != nil {
			return fmt.Errorf("emit public %s: %w", e.outName, err)
		}
		if err := os.WriteFile(filepath.Join(cfg.publicDir, e.outName), []byte(pubCode), 0o644); err != nil {
			return fmt.Errorf("write public %s: %w", e.outName, err)
		}
	}

	// Generate register.go aggregator for capi directory.
	return writeRawRegisterAggregator(cfg.capiDir, registerNames)
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

	// Map CEF's platform suffixes to Go's filename-based build suffixes.
	switch {
	case strings.HasSuffix(name, "_mac"):
		name = strings.TrimSuffix(name, "_mac") + "_darwin"
	case strings.HasSuffix(name, "_win"):
		name = strings.TrimSuffix(name, "_win") + "_windows"
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

	return name + "_gen.go"
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
	discoveredNames, err := discoverRegisterNamesFromDir(dir)
	if err != nil {
		return err
	}
	names = mergeUniqueRegisterNames(names, discoveredNames)
	commonNames, linuxNames, darwinNames, windowsNames := splitRegisterNames(names)

	var sb strings.Builder
	sb.WriteString("package capi\n\n")
	sb.WriteString("// Register loads all CEF C API symbols from the shared library.\n")
	sb.WriteString("// Code generated by cefgen. DO NOT EDIT.\n")
	sb.WriteString("func Register(handle uintptr) {\n")
	for _, name := range commonNames {
		sb.WriteString("\t" + name + "(handle)\n")
	}
	sb.WriteString("\tregisterPlatform(handle)\n")
	sb.WriteString("}\n")
	if err := os.WriteFile(filepath.Join(dir, "register_gen.go"), []byte(sb.String()), 0o644); err != nil {
		return err
	}
	if err := writePlatformRegisterFile(dir, "linux", linuxNames); err != nil {
		return err
	}
	if err := writePlatformRegisterFile(dir, "darwin", darwinNames); err != nil {
		return err
	}
	if err := writePlatformRegisterFile(dir, "windows", windowsNames); err != nil {
		return err
	}
	return writePlatformRegisterDefault(dir)
}

func discoverRegisterNamesFromDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`func (Register[A-Za-z0-9_]+)\(`)
	seen := make(map[string]struct{})
	var names []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasPrefix(name, "register") || name == "doc.go" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		for _, match := range re.FindAllStringSubmatch(string(data), -1) {
			fnName := match[1]
			if _, ok := seen[fnName]; ok {
				continue
			}
			seen[fnName] = struct{}{}
			names = append(names, fnName)
		}
	}

	sort.Strings(names)
	return names, nil
}

func mergeUniqueRegisterNames(groups ...[]string) []string {
	seen := make(map[string]struct{})
	var merged []string
	for _, group := range groups {
		for _, name := range group {
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			merged = append(merged, name)
		}
	}
	sort.Strings(merged)
	return merged
}

func splitRegisterNames(names []string) (common, linux, darwin, windows []string) {
	for _, name := range names {
		switch platformRegisterName(name) {
		case "linux":
			linux = append(linux, name)
		case "darwin":
			darwin = append(darwin, name)
		case "windows":
			windows = append(windows, name)
		default:
			common = append(common, name)
		}
	}
	return common, linux, darwin, windows
}

func platformRegisterName(name string) string {
	switch {
	case strings.HasSuffix(name, "Linux"):
		return "linux"
	case strings.HasSuffix(name, "Mac"):
		return "darwin"
	case strings.HasSuffix(name, "Win"):
		return "windows"
	default:
		return ""
	}
}

func writePlatformRegisterFile(dir, goos string, names []string) error {
	var sb strings.Builder
	sb.WriteString("//go:build " + goos + "\n\n")
	sb.WriteString("package capi\n\n")
	sb.WriteString("// Code generated by cefgen. DO NOT EDIT.\n")
	sb.WriteString("func registerPlatform(handle uintptr) {\n")
	for _, name := range names {
		sb.WriteString("\t" + name + "(handle)\n")
	}
	sb.WriteString("}\n")
	return os.WriteFile(filepath.Join(dir, "register_platform_"+goos+"_gen.go"), []byte(sb.String()), 0o644)
}

func writePlatformRegisterDefault(dir string) error {
	var sb strings.Builder
	sb.WriteString("//go:build !linux && !darwin && !windows\n\n")
	sb.WriteString("package capi\n\n")
	sb.WriteString("// Code generated by cefgen. DO NOT EDIT.\n")
	sb.WriteString("func registerPlatform(_ uintptr) {}\n")
	return os.WriteFile(filepath.Join(dir, "register_platform_default_gen.go"), []byte(sb.String()), 0o644)
}
