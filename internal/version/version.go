package version

import (
	"fmt"
	"regexp"
	"strings"
)

var semverPattern = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)$`)

// Normalize validates a semantic version and strips a leading "v".
// It returns the normalized version in "X.Y.Z" form.
func Normalize(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("version is required")
	}
	match := semverPattern.FindStringSubmatch(trimmed)
	if match == nil {
		return "", fmt.Errorf("version %q must be in the form vX.Y.Z or X.Y.Z", raw)
	}
	return fmt.Sprintf("%s.%s.%s", match[1], match[2], match[3]), nil
}

// IsDev reports whether the version string represents a dev build.
func IsDev(raw string) bool {
	return strings.TrimSpace(raw) == "dev"
}
