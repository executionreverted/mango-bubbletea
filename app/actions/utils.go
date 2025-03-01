// app/actions/utils.go
package actions

import (
	"strings"
)

// FormatTitle formats a string for display as a title
func FormatTitle(title string) string {
	return strings.ToUpper(title)
}

// TruncateText truncates text to a specified length
func TruncateText(text string, length int) string {
	if len(text) <= length {
		return text
	}
	return text[:length-3] + "..."
}
