package envfile

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "parses export and quoted values",
			input: `
# comment
export TOKEN=value
OTHER = "spaced value"
`,
			want: map[string]string{
				"TOKEN": "value",
				"OTHER": "spaced value",
			},
		},
		{
			name:    "invalid line",
			input:   "INVALID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPatch(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		updates  map[string]string
		contains []string
		absent   []string
	}{
		{
			name:     "add new secret",
			input:    `EXISTING=1`,
			updates:  map[string]string{"NEW": "secret"},
			contains: []string{`NEW=secret`},
		},
		{
			name:     "replace existing secret",
			input:    `KEY=old`,
			updates:  map[string]string{"KEY": "new"},
			contains: []string{`KEY=new`},
		},
		{
			name:    "replace export line",
			input:   `export KEY=old`,
			updates: map[string]string{"KEY": "new"},
			contains: []string{
				`KEY=new`,
			},
			absent: []string{
				`export KEY=old`,
			},
		},
		{
			name:    "replace spaced assignment",
			input:   `KEY = old`,
			updates: map[string]string{"KEY": "new"},
			contains: []string{
				`KEY=new`,
			},
			absent: []string{
				`KEY = old`,
			},
		},
		{
			name:    "dedupe existing key lines",
			input:   "KEY=old\nexport KEY=older\nOTHER=1",
			updates: map[string]string{"KEY": "new"},
			contains: []string{
				`KEY=new`,
				`OTHER=1`,
			},
			absent: []string{
				`export KEY=older`,
			},
		},
		{
			name:     "quote complex secret",
			input:    ``,
			updates:  map[string]string{"COMPLEX": "hash # check"},
			contains: []string{`COMPLEX="hash # check"`},
		},
		{
			name:     "escape quotes and backslashes",
			input:    ``,
			updates:  map[string]string{"COMPLEX": `C:\path\"file"`},
			contains: []string{`COMPLEX="C:\\path\\\"file\""`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Patch(tt.input, tt.updates)
			for _, c := range tt.contains {
				assert.Contains(t, got, c)
			}
			for _, c := range tt.absent {
				assert.NotContains(t, got, c)
			}
			if tt.name == "dedupe existing key lines" {
				assert.Equal(t, 1, strings.Count(got, "KEY="))
			}
		})
	}
}
