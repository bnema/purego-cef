package parser

import (
	"regexp"
	"strings"
)

// docBlock represents a contiguous block of /// comments ending at a specific line.
type docBlock struct {
	lines   []string // raw /// lines (trimmed of leading whitespace and "///")
	endLine int      // 0-based line index of the last /// line in the block
}

// extractDocBlocks scans raw source lines and collects contiguous /// comment blocks.
// Each block records its cleaned lines and the 0-based line number where the block ends.
func extractDocBlocks(raw string) []docBlock {
	lines := strings.Split(raw, "\n")
	var blocks []docBlock
	var cur []string
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if text, ok := strings.CutPrefix(trimmed, "///"); ok {
			// Strip a single leading space if present.
			if len(text) > 0 && text[0] == ' ' {
				text = text[1:]
			}
			cur = append(cur, text)
		} else {
			if len(cur) > 0 {
				blocks = append(blocks, docBlock{
					lines:   cur,
					endLine: i - 1,
				})
				cur = nil
			}
		}
	}
	if len(cur) > 0 {
		blocks = append(blocks, docBlock{
			lines:   cur,
			endLine: len(lines) - 1,
		})
	}
	return blocks
}

var allocRE = regexp.MustCompile(`(?i)NOTE:\s*This\s+struct\s+is\s+allocated\s+(client-side|DLL-side)`)

// classifyKind inspects a doc block's lines for the allocation NOTE.
// Returns "handler" for client-side, "object" for DLL-side, or "" if not found.
func classifyKind(lines []string) string {
	for _, line := range lines {
		if m := allocRE.FindStringSubmatch(line); m != nil {
			switch strings.ToLower(m[1]) {
			case "client-side":
				return "handler"
			case "dll-side":
				return "object"
			}
		}
	}
	return ""
}

// cleanDoc joins doc lines into a single string, stripping NOTE: lines
// and empty-comment delimiter lines (bare "///").
func cleanDoc(lines []string) string {
	var out []string
	for _, line := range lines {
		// Skip NOTE: allocation lines.
		if allocRE.MatchString(line) {
			continue
		}
		// Skip bare empty lines (from "///" with nothing after).
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, " ")
}

// docIndex maps source line numbers to doc blocks for quick lookup.
// Given a declaration starting at line N, the associated doc block
// is the one whose endLine == N-1 (immediately preceding the declaration).
type docIndex struct {
	// byEndLine maps endLine -> index into blocks slice.
	byEndLine map[int]int
	blocks    []docBlock
}

// buildDocIndex creates a docIndex from raw source content.
func buildDocIndex(raw string) *docIndex {
	blocks := extractDocBlocks(raw)
	idx := &docIndex{
		byEndLine: make(map[int]int, len(blocks)),
		blocks:    blocks,
	}
	for i, b := range blocks {
		idx.byEndLine[b.endLine] = i
	}
	return idx
}

// forLine returns the doc block (if any) that immediately precedes the given
// 0-based line number. It searches for a block ending at line-1 (skipping
// blank lines between the comment block and the declaration).
func (d *docIndex) forLine(line int) *docBlock {
	// Check up to 3 lines above for the doc block end
	// (there may be blank lines between the comment and declaration).
	for offset := 1; offset <= 3; offset++ {
		if i, ok := d.byEndLine[line-offset]; ok {
			return &d.blocks[i]
		}
	}
	return nil
}

// findLineOf returns the 0-based line number where needle first appears in raw.
// Returns -1 if not found.
func findLineOf(raw string, needle string) int {
	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		if strings.Contains(line, needle) {
			return i
		}
	}
	return -1
}

// findFieldDoc searches within the struct body in the raw source for a doc block
// preceding a field declaration. fieldName is the C callback/field name.
// structName is the typedef name (e.g., "cef_client_t") used to scope the search.
func findFieldDoc(raw string, structName, fieldName string) string {
	// Find the struct body boundaries.
	lines := strings.Split(raw, "\n")
	structStart := -1
	structEnd := -1
	// Look for "typedef struct _<name>" or similar.
	internalName := "_" + structName
	for i, line := range lines {
		if strings.Contains(line, internalName) && strings.Contains(line, "typedef struct") {
			structStart = i
		}
		if structStart >= 0 && structEnd < 0 && strings.HasPrefix(strings.TrimSpace(line), "}") {
			structEnd = i
			break
		}
	}
	if structStart < 0 || structEnd < 0 {
		return ""
	}

	// Build a doc index for just the struct body region.
	bodyLines := lines[structStart : structEnd+1]
	bodyRaw := strings.Join(bodyLines, "\n")
	idx := buildDocIndex(bodyRaw)

	// Find the field within the body.
	for i, line := range bodyLines {
		if strings.Contains(line, fieldName) {
			if db := idx.forLine(i); db != nil {
				return cleanDoc(db.lines)
			}
			break
		}
	}
	return ""
}
