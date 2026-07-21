package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bnema/purego-cef/internal/cefapi"
)

var (
	apiVersionRE = regexp.MustCompile(fmt.Sprintf(`^\s*#define\s+CEF_API_VERSION_%[1]d\s+%[1]d\s*$`, cefapi.Version))
	apiHashRE    = regexp.MustCompile(fmt.Sprintf(`^\s*#define\s+CEF_API_HASH_%d\s+"([0-9a-f]+)"\s*$`, cefapi.Version))
)

// readAPIContract extracts the exact Linux ABI hash for the target API. It
// deliberately does not infer a newest API: generated bindings have one ABI.
func readAPIContract(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read CEF API versions: %w", err)
	}
	foundVersion, linux := false, false
	var hash string
	for line := range strings.SplitSeq(string(data), "\n") {
		if apiVersionRE.MatchString(line) {
			foundVersion = true
		}
		if strings.Contains(line, "#elif defined(OS_LINUX)") || strings.Contains(line, "#if defined(OS_LINUX)") {
			linux = true
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "#elif") || strings.HasPrefix(strings.TrimSpace(line), "#else") || strings.HasPrefix(strings.TrimSpace(line), "#endif") {
			linux = false
			continue
		}
		if linux {
			if match := apiHashRE.FindStringSubmatch(line); match != nil {
				hash = match[1]
			}
		}
	}
	if !foundVersion {
		return "", fmt.Errorf("headers do not support required CEF API %d", cefapi.Version)
	}
	if hash == "" {
		return "", fmt.Errorf("headers do not provide Linux hash for CEF API %d", cefapi.Version)
	}
	return hash, nil
}
