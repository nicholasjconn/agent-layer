package doctor

// Status represents the health status of a check (OK, WARN, FAIL).
type Status string

const (
	// StatusOK indicates the check passed.
	StatusOK Status = "OK"
	// StatusWarn indicates a potential issue that doesn't block functionality.
	StatusWarn Status = "WARN"
	// StatusFail indicates a critical issue that must be resolved.
	StatusFail Status = "FAIL"
)

// Result holds the outcome of a single diagnostic check.
type Result struct {
	Status         Status
	CheckName      string
	Message        string
	Recommendation string
}
