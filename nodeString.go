package neon

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// StringNode represents a string node in NEON.
type StringNode struct {
	Value string
}

var escapeSequences = map[string]string{
	"t": "\t", "n": "\n", "r": "\r", "f": "\x0C", "b": "\x08",
	`"`: `"`, `\`: `\`, `/`: `/`, "_": "\u00A0",
}

// NewStringNode creates a new StringNode.
func NewStringNode(value string) *StringNode {
	return &StringNode{Value: value}
}

// ToValue returns the string value of the node.
func (n *StringNode) ToValue() string {
	return n.Value
}

// Parse parses a NEON string into its unescaped representation.
func Parse(input string, position Position) (string, error) {
	var res string
	if strings.HasPrefix(input, "'''") || strings.HasPrefix(input, `"""`) { // Multiline
		res = input[3 : len(input)-3]
		lines := strings.Split(res, "\n")
		indent := regexp.MustCompile(`^[\t ]*`)
		for i, line := range lines {
			if i > 0 {
				lines[i] = indent.ReplaceAllString(line, "")
			}
		}
		res = strings.Join(lines, "\n")
		res = strings.Trim(res, "\n")
	} else { // Single-line
		res = input[1 : len(input)-1]
		if strings.HasPrefix(input, "'") {
			res = strings.ReplaceAll(res, "''", "'")
			return res, nil
		}
	}

	// Replace escape sequences
	escapedRegex := regexp.MustCompile(`\\(?:u[0-9a-fA-F]{4}|.)`)
	res = escapedRegex.ReplaceAllStringFunc(res, func(match string) string {
		if val, ok := escapeSequences[string(match[1])]; ok {
			return val
		} else if strings.HasPrefix(match, `\u`) && len(match) == 6 {
			var decoded string
			err := json.Unmarshal([]byte(`"`+match+`"`), &decoded)
			if err == nil {
				return decoded
			}
		}
		panic(fmt.Sprintf("Invalid escape sequence: %s", match))
	})
	return res, nil
}

// ToString converts the string value back to a NEON-compatible string representation.
func (n *StringNode) ToString() string {
	if !strings.Contains(n.Value, "\n") {
		if containsControlChars(n.Value) {
			escaped, _ := json.Marshal(n.Value)
			return string(escaped)
		}
		return "'" + strings.ReplaceAll(n.Value, "'", "''") + "'"
	}

	if containsControlChars(n.Value) || strings.Contains(n.Value, "\n'''") {
		escaped := escapeJSONCompatible(n.Value)
		escaped = strings.ReplaceAll(escaped, `"""`, `""\"`)
		return fmt.Sprintf(`"""%s"""`, escaped)
	}

	// Default multiline string
	lines := strings.Split(n.Value, "\n")
	for i := range lines {
		lines[i] = "\t" + lines[i]
	}
	return fmt.Sprintf("'''\n%s\n'''", strings.Join(lines, "\n"))
}

// Helper to check if a string contains control characters.
func containsControlChars(s string) bool {
	for _, r := range s {
		if r < 0x20 && r != '\n' && r != '\t' {
			return true
		}
	}
	return false
}

// Helper to escape a string compatible with JSON.
func escapeJSONCompatible(s string) string {
	escaped, _ := json.Marshal(s)
	return string(escaped[1 : len(escaped)-1])
}
