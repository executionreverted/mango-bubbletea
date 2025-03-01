package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Define common styles
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5A56E0")).
			Bold(true).
			Padding(1).
			MarginBottom(1)

	PageStyle = lipgloss.NewStyle().
			Padding(1)
)
