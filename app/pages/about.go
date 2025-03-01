package pages

import (
	"bubbletea-app/app/components"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AboutModel struct {
	header components.HeaderModel
	footer components.FooterModel
	width  int
	height int
}

func NewAboutModel() AboutModel {
	return AboutModel{
		header: components.NewHeaderModel("About"),
		footer: components.NewFooterModel(),
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
	footerView := m.footer.View(m.width)

	contentHeight := m.height - lipgloss.Height(headerView) - lipgloss.Height(footerView) - 2
	aboutText := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(contentHeight). // SET HEIGHT FOR CONTENT
		Padding(1).
		Render("a nice starting point for a bubbletea terminal app")
	content := fmt.Sprintf(
		"%s\n%s\n%s",
		headerView,
		aboutText,
		footerView,
	)

	contentContainer := lipgloss.NewStyle().
		Height(contentHeight).Render(content)
	return contentContainer
}
