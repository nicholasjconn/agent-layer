package wizard

import (
	"fmt"
	"strings"
	"unicode"
)

type multilineState int

const (
	multilineNone multilineState = iota
	multilineBasic
	multilineLiteral
)

// formatTomlNoIndent removes leading indentation from TOML output while preserving multiline strings.
// content is the TOML text to format; returns the formatted TOML or an error if multiline strings are unterminated.
func formatTomlNoIndent(content string) (string, error) {
	lines := strings.Split(content, "\n")
	state := multilineNone

	for i, line := range lines {
		if state == multilineNone {
			lines[i] = strings.TrimLeftFunc(line, unicode.IsSpace)
		}

		lineForScan := stripTomlInlineComment(line)
		nextState, err := advanceMultilineState(lineForScan, state)
		if err != nil {
			return "", err
		}
		state = nextState
	}

	if state != multilineNone {
		return "", fmt.Errorf("unterminated multiline string in TOML output")
	}

	return strings.Join(lines, "\n"), nil
}

// advanceMultilineState updates the multiline string parsing state for a single line.
// line is the TOML line without inline comments; state is the current multiline state; returns the next state or error.
func advanceMultilineState(line string, state multilineState) (multilineState, error) {
	const (
		delimBasic   = `"""`
		delimLiteral = `'''`
	)
	line = strings.TrimSpace(line)
	if line == "" {
		return state, nil
	}

	switch state {
	case multilineBasic:
		if strings.Count(line, delimBasic)%2 == 1 {
			return multilineNone, nil
		}
		return state, nil
	case multilineLiteral:
		if strings.Count(line, delimLiteral)%2 == 1 {
			return multilineNone, nil
		}
		return state, nil
	default:
		basicCount := strings.Count(line, delimBasic)
		literalCount := strings.Count(line, delimLiteral)
		if basicCount%2 == 1 && literalCount%2 == 1 {
			return multilineNone, fmt.Errorf("ambiguous multiline delimiters in TOML output")
		}
		if basicCount%2 == 1 {
			return multilineBasic, nil
		}
		if literalCount%2 == 1 {
			return multilineLiteral, nil
		}
		return multilineNone, nil
	}
}

// stripTomlInlineComment removes inline comments from a TOML line.
// line is a TOML line; returns the substring before the first '#' character.
func stripTomlInlineComment(line string) string {
	if idx := strings.Index(line, "#"); idx >= 0 {
		return line[:idx]
	}
	return line
}
