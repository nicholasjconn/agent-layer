package wizard

// approvalModeValues returns the canonical approval mode values in order.
func approvalModeValues() []string {
	values := make([]string, 0, len(ApprovalModeOptions))
	for _, option := range ApprovalModeOptions {
		values = append(values, option.Value)
	}
	return values
}

// approvalModeLabel returns the display label for an approval mode option.
func approvalModeLabel(option ApprovalModeOption) string {
	return option.Value + " - " + option.Description
}

// approvalModeLabels returns the display labels for all approval modes.
func approvalModeLabels() []string {
	labels := make([]string, 0, len(ApprovalModeOptions))
	for _, option := range ApprovalModeOptions {
		labels = append(labels, approvalModeLabel(option))
	}
	return labels
}

// approvalModeLabelForValue returns the display label for a canonical value.
func approvalModeLabelForValue(value string) (string, bool) {
	for _, option := range ApprovalModeOptions {
		if option.Value == value {
			return approvalModeLabel(option), true
		}
	}
	return "", false
}

// approvalModeValueForLabel returns the canonical value for a display label.
func approvalModeValueForLabel(label string) (string, bool) {
	for _, option := range ApprovalModeOptions {
		if approvalModeLabel(option) == label {
			return option.Value, true
		}
	}
	return "", false
}
