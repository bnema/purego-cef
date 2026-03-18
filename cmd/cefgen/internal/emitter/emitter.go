package emitter

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"text/template"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Emit takes a parsed Header and returns formatted Go source code.
func Emit(header *model.Header) (string, error) {
	tmpl, err := template.New("file").ParseFS(templateFS, "templates/*.tmpl")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "file.tmpl", header); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format source: %w\n%s", err, buf.String())
	}
	return string(formatted), nil
}
