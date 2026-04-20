package emitter

import "github.com/bnema/purego-cef/cmd/cefgen/internal/model"

// typeRenames maps generated public type names to non-colliding exported names.
// Use this when the natural public name would conflict with a handwritten API
// surface that intentionally exposes a safer abstraction.
var typeRenames = map[string]string{
	"Settings": "CEFSettings",
	"MainArgs": "CEFMainArgs",
}

func publicTypeName(name string) string {
	if renamed, ok := typeRenames[name]; ok {
		return renamed
	}
	return name
}

func publicTypeNameForCName(cName string) string {
	return publicTypeName(model.PublicName(cName))
}
