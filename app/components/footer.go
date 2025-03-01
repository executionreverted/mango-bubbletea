package components

import (
	"github.com/charmbracelet/lipgloss"
)

type FooterModel struct{}

func NewFooterModel() FooterModel {
	return FooterModel{}
}

func (m FooterModel) View(width int) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#2F4858")).
		Padding(0, 1).
		Width(width - 2).
		Render("Bubble Tea App Boilerplate • github.com/yourusername/bubbletea-boilerplate")
}
