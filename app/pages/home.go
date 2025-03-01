package pages

import (
	"bubbletea-app/app/actions"
	"bubbletea-app/app/components"
	"bubbletea-app/app/config"
	"bubbletea-app/app/styles"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeModel struct {
	list   list.Model
	header components.HeaderModel
	footer components.FooterModel
	width  int
	height int
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewHomeModel(keyMap config.KeyMap) HomeModel {
	// Create example items
	items := []list.Item{
		item{title: "Task 1", desc: "Description for task 1"},
		item{title: "Task 2", desc: "Description for task 2"},
		item{title: "Task 3", desc: "Description for task 3"},
	}

	// Setup list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Home Page"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.Styles.Title = styles.TitleStyle

	// Initialize with data from API
	apiData := actions.FetchData()
	for _, d := range apiData {
		l.InsertItem(len(items), item{title: d.Title, desc: d.Description})
	}

	return HomeModel{
		list:   l,
		header: components.NewHeaderModel("Home"),
		footer: components.NewFooterModel(),
	}
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-2, msg.Height-7)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m HomeModel) View() string {
	headerView := m.header.View(m.width)
	footerView := m.footer.View(m.width)
	contentHeight := m.height - lipgloss.Height(headerView) - lipgloss.Height(footerView) - 2
	content := fmt.Sprintf(
		"%s\n%s\n%s",
		headerView,
		m.list.View(),
		footerView,
	)
	contentContainer := lipgloss.NewStyle().
		Height(contentHeight).Render(content)

	return contentContainer
}
