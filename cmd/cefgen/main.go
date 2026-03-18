package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/emitter"
	"github.com/bnema/purego-cef/cmd/cefgen/internal/parser"
)

type config struct {
	headersDir string
	outputDir  string
	version    string
}

func (c config) validate() error {
	if c.headersDir == "" || c.outputDir == "" {
		return fmt.Errorf("--headers-dir and --output-dir are required")
	}
	return nil
}

func main() {
	var cfg config
	flag.StringVar(&cfg.headersDir, "headers-dir", "", "CEF include root")
	flag.StringVar(&cfg.outputDir, "output-dir", "", "generated package directory")
	flag.StringVar(&cfg.version, "version", "145", "target major version")
	flag.Parse()
	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cfg config) error {
	if err := cfg.validate(); err != nil {
		return err
	}
	for _, rel := range []string{
		"capi/cef_base_capi.h",
		"capi/cef_app_capi.h",
		"capi/cef_client_capi.h",
		"capi/cef_life_span_handler_capi.h",
		"capi/cef_render_handler_capi.h",
		"capi/cef_browser_capi.h",
	} {
		header, err := parser.ParseFile(filepath.Join(cfg.headersDir, rel))
		if err != nil {
			return err
		}
		header.RegisterName = registerName(filepath.Base(rel))
		code, err := emitter.Emit(header)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(filepath.Base(rel), "_capi.h") + ".go"
		if err := os.WriteFile(filepath.Join(cfg.outputDir, name), []byte(code), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// registerName derives a unique Go register function name from a capi header
// filename, e.g. "cef_app_capi.h" -> "RegisterApp".
func registerName(base string) string {
	// Strip _capi.h suffix
	name := strings.TrimSuffix(base, "_capi.h")
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
