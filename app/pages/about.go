package pages

import (
	"bubbletea-app/app/components"
	"bubbletea-app/app/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AboutModel struct {
	header components.HeaderModel
	width  int
	height int
}

func NewAboutModel() AboutModel {
	return AboutModel{
		header: components.NewHeaderModel("About"),
	}
}

func (m AboutModel) Init() tea.Cmd {
	return nil
}

func (m AboutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m AboutModel) View() string {
	headerView := m.header.View(m.width)

	aboutText := lipgloss.NewStyle().
		Padding(1).
		Render("Bubble Tea App Boilerplate\nVersion 1.0.0\n\nA Next.js inspired TUI application structure\nwith pages and actions folders.")

	content := fmt.Sprintf(
		"%s\n%s",
		headerView,
		aboutText,
	)
	contentHeight := m.height - 4 // Subtract nav and footer height
	contentContainer := lipgloss.NewStyle().
		Height(contentHeight).Render(content)
	return styles.PageStyle.Width(m.width - 2).Render(contentContainer)
}
