package version

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{input: "0.5.0", want: "0.5.0"},
		{input: "v0.5.0", want: "0.5.0"},
		{input: "  v1.2.3  ", want: "1.2.3"},
		{input: "", wantErr: true},
		{input: "v1.2", wantErr: true},
		{input: "1.2.3.4", wantErr: true},
		{input: "dev", wantErr: true},
	}
	for _, tt := range tests {
		got, err := Normalize(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Fatalf("Normalize(%q) expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("Normalize(%q) error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Fatalf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsDev(t *testing.T) {
	if !IsDev("dev") {
		t.Fatalf("expected IsDev(dev) to be true")
	}
	if IsDev("v0.5.0") {
		t.Fatalf("expected IsDev(v0.5.0) to be false")
	}
}
