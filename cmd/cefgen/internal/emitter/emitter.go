package emitter

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// EmitPublic takes a PublicFileData view model and returns formatted Go source
// for the public API layer.
func EmitPublic(data *PublicFileData) (string, error) {
	funcMap := template.FuncMap{
		"lower": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToLower(s[:1]) + s[1:]
		},
		"zeroVal": func(typ string) string {
			switch typ {
			case "bool":
				return "false"
			case "string":
				return `""`
			case "int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"uintptr", "float32", "float64":
				return "0"
			default:
				if strings.HasPrefix(typ, "*") || strings.HasPrefix(typ, "[]") ||
					strings.HasPrefix(typ, "map[") || typ == "unsafe.Pointer" {
					return "nil"
				}
				// Interface types default to nil.
				return "nil"
			}
		},
		"isEnumReturn": func(ret ReturnData) bool {
			return ret.IsEnum
		},
	}

	tmpl, err := template.New("public").Funcs(funcMap).ParseFS(templateFS,
		"templates/public_file.tmpl",
		"templates/interface.tmpl",
		"templates/object_wrapper.tmpl",
		"templates/handler_constructor.tmpl",
		"templates/data_struct.tmpl",
		"templates/enums.tmpl",
		"templates/free_func.tmpl",
	)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "public_file.tmpl", data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format source: %w\n%s", err, buf.String())
	}
	return string(formatted), nil
}

// EmitRaw takes a parsed Header and returns formatted Go source for the raw layer.
// The raw layer lives in package "raw" and mirrors C struct layouts.
func EmitRaw(header *model.Header) (string, error) {
	tmpl, err := template.New("raw").ParseFS(templateFS, "templates/raw_*.tmpl")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "raw_file.tmpl", header); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format source: %w\n%s", err, buf.String())
	}
	return string(formatted), nil
}
