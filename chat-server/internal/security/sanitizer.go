package security

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	// ANSI escape sequence pattern (CSI + OSC)
	ansiPattern = regexp.MustCompile(`\x1b(?:\[[0-9;]*[a-zA-Z]|\][0-9;]*(?:\\\|[a-zA-Z0-9=:])*\x07|[\(\)].|[\[\]#;?].*?(?:\x1b\\|\x07))`)
	maxMessageLength = 2000
)

func SanitizeInput(input string) string {
	originalLen := len(input)

	if len(input) > maxMessageLength {
		input = input[:maxMessageLength]
	}
	input = ansiPattern.ReplaceAllString(input, "")

	var filtered strings.Builder
	filtered.Grow(len(input))
	for _, r := range input {
		if unicode.IsPrint(r) || r == '\n' || r == '\t' || r == '\r' || r == ' ' {
			if r != '\x1b' {
				filtered.WriteRune(r)
			}
		}
	}

	result := filtered.String()
	fmt.Printf("Sanitized: %d chars -> %d chars\n", originalLen, len(result))
	return result
}