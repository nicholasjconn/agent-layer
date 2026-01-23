package wizard

import "strings"

// approvalModeHelpText returns explanatory text for approval modes.
func approvalModeHelpText() string {
	lines := []string{
		"Approval modes control what runs without prompts:",
	}
	for _, option := range ApprovalModeOptions {
		lines = append(lines, "- "+option.Value+": "+option.Description)
	}
	lines = append(lines, "Support varies by client; Agent Layer applies the closest available behavior.")
	return strings.Join(lines, "\n")
}

// previewModelWarningText returns the warning text shown before preview model selection.
func previewModelWarningText() string {
	return "Preview models are pre-release and can change or be removed without notice."
}

// hasPreviewModels reports whether any model option looks like a preview release.
func hasPreviewModels(options []string) bool {
	for _, option := range options {
		if strings.Contains(option, "preview") {
			return true
		}
	}
	return false
}
