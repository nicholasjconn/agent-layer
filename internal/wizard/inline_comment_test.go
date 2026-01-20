package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInlineComment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no comment",
			input:    `key = "value"`,
			expected: "",
		},
		{
			name:     "simple comment",
			input:    `key = "value" # comment`,
			expected: "comment",
		},
		{
			name:     "comment with spaces",
			input:    `key = "value"   #    spaced   `,
			expected: "spaced",
		},
		{
			name:     "hash in string",
			input:    `key = "val#ue" # comment`,
			expected: "comment",
		},
		{
			name:     "hash in single quoted string",
			input:    `key = 'val#ue' # comment`,
			expected: "comment",
		},
		{
			name:     "escaped quote",
			input:    `key = "val\"ue" # comment`,
			expected: "comment",
		},
		// The following cases are expected to FAIL with current implementation
		{
			name:     "hash in triple double quotes",
			input:    `key = """val#ue""" # comment`,
			expected: "comment",
		},
		{
			name:     "hash in triple single quotes",
			input:    `key = '''val#ue''' # comment`,
			expected: "comment",
		},
		{
			name:     "triple quote with inner quote and hash",
			input:    `key = """ " # """ # real comment`,
			expected: "real comment",
		},
		{
			name:     "multiline string start",
			input:    `key = """line1`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inlineComment(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
