package parser

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// PreprocessCEFAPI selects branches guarded by CEF API-version expressions.
// Non-CEF include guards and language/platform wrappers are retained as source
// (with their directives removed later by stripComments), because the parser
// cannot safely evaluate those conditions.
func PreprocessCEFAPI(data []byte, target int) ([]byte, error) {
	lines := bytes.Split(data, []byte("\n"))
	out := make([][]byte, 0, len(lines))
	stack := make([]conditionalFrame, 0)
	active := true
	for lineNo, line := range lines {
		trimmed := strings.TrimSpace(string(line))
		if !strings.HasPrefix(trimmed, "#") {
			if active {
				out = append(out, line)
			}
			continue
		}
		keyword, expression := directive(trimmed)
		switch keyword {
		case "if", "ifdef", "ifndef":
			cef, selected, err := cefCondition(keyword, expression, target)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNo+1, err)
			}
			frame := conditionalFrame{parentActive: active, cef: cef}
			if cef {
				frame.taken = selected
				frame.active = active && selected
			} else {
				frame.active = active
			}
			stack = append(stack, frame)
			active = frame.active
		case "elif":
			if len(stack) == 0 {
				return nil, fmt.Errorf("line %d: #elif without #if", lineNo+1)
			}
			frame := &stack[len(stack)-1]
			if !frame.cef {
				if cefAPICallRE.MatchString(expression) {
					return nil, fmt.Errorf("line %d: CEF #elif without CEF #if", lineNo+1)
				}
				active = frame.parentActive
				frame.active = active
				continue
			}
			if frame.elseSeen {
				return nil, fmt.Errorf("line %d: #elif after #else", lineNo+1)
			}
			_, selected, err := cefCondition("if", expression, target)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNo+1, err)
			}
			frame.active = frame.parentActive && !frame.taken && selected
			frame.taken = frame.taken || selected
			active = frame.active
		case "else":
			if len(stack) == 0 {
				return nil, fmt.Errorf("line %d: #else without #if", lineNo+1)
			}
			frame := &stack[len(stack)-1]
			if frame.elseSeen {
				return nil, fmt.Errorf("line %d: duplicate #else", lineNo+1)
			}
			frame.elseSeen = true
			if frame.cef {
				frame.active = frame.parentActive && !frame.taken
			} else {
				frame.active = frame.parentActive
			}
			active = frame.active
		case "endif":
			if len(stack) == 0 {
				return nil, fmt.Errorf("line %d: #endif without #if", lineNo+1)
			}
			frame := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			active = frame.parentActive
		default:
			// Keep ordinary directives for the next stage. In particular,
			// stripComments expands numeric #defines; include and wrapper
			// directives are discarded there.
			if active {
				out = append(out, line)
			}
		}
	}
	if len(stack) != 0 {
		return nil, fmt.Errorf("unterminated preprocessor conditional")
	}
	return bytes.Join(out, []byte("\n")), nil
}

type conditionalFrame struct{ parentActive, active, cef, taken, elseSeen bool }

func directive(line string) (string, string) {
	line = strings.TrimSpace(strings.TrimPrefix(line, "#"))
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "", ""
	}
	return fields[0], strings.TrimSpace(strings.TrimPrefix(line, fields[0]))
}

var (
	cefExpressionRE = regexp.MustCompile(`^CEF_API_(ADDED|REMOVED|RANGE)\s*\(\s*([0-9]+|CEF_EXPERIMENTAL)\s*(?:,\s*([0-9]+|CEF_EXPERIMENTAL)\s*)?\)$`)
	cefAPICallRE    = regexp.MustCompile(`\bCEF_API_[A-Z_]+\s*\(`)
)

func cefCondition(keyword, expression string, target int) (isCEF, selected bool, err error) {
	if keyword != "if" {
		return false, false, nil
	}
	if !cefAPICallRE.MatchString(expression) {
		return false, false, nil
	}
	m := cefExpressionRE.FindStringSubmatch(expression)
	if m == nil {
		return true, false, fmt.Errorf("unsupported CEF API expression %q", expression)
	}
	version := func(raw string) (int, error) {
		if raw == "CEF_EXPERIMENTAL" {
			return 999999, nil
		}
		return strconv.Atoi(raw)
	}
	first, err := version(m[2])
	if err != nil {
		return true, false, err
	}
	switch m[1] {
	case "ADDED":
		if m[3] != "" {
			return true, false, fmt.Errorf("invalid CEF_API_ADDED arguments")
		}
		return true, target >= first, nil
	case "REMOVED":
		if m[3] != "" {
			return true, false, fmt.Errorf("invalid CEF_API_REMOVED arguments")
		}
		return true, target < first, nil
	case "RANGE":
		if m[3] == "" {
			return true, false, fmt.Errorf("invalid CEF_API_RANGE arguments")
		}
		last, err := version(m[3])
		if err != nil {
			return true, false, err
		}
		return true, target >= first && target < last, nil
	}
	return true, false, fmt.Errorf("unsupported CEF API expression %q", expression)
}
