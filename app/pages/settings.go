// app/pages/settings_page.go
package pages

import (
	"bubbletea-app/app/components"
	"bubbletea-app/app/config"
	"bubbletea-app/app/styles"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	InputFocusChangedMsg bool
	SettingsModel        struct {
		inputs     []textinput.Model
		saveButton components.ButtonModel
		header     components.HeaderModel
		keyMap     config.KeyMap
		width      int
		height     int
		focusIndex int // tracks which element has focus
		focusables int // total focusable elements
	}
)

func NewSettingsModel(keyMap config.KeyMap) SettingsModel {
	// Create inputs
	hostInput := textinput.New()
	hostInput.Placeholder = "Enter host (e.g., localhost)"
	hostInput.Width = 30
	hostInput.Blur()

	portInput := textinput.New()
	portInput.Placeholder = "Enter port (e.g., 8080)"
	portInput.Width = 30
	portInput.Blur()

	apiKeyInput := textinput.New()
	apiKeyInput.Placeholder = "Enter API key"
	apiKeyInput.Width = 30
	apiKeyInput.Blur()

	// Create button with save action
	saveButton := components.NewButtonModel("Save Configuration", nil)

	return SettingsModel{
		inputs:     []textinput.Model{hostInput, portInput, apiKeyInput},
		saveButton: saveButton,
		header:     components.NewHeaderModel("Settings"),
		keyMap:     keyMap,
		focusIndex: 0,
		focusables: 4, // 3 inputs + 1 button
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// Check if we're currently editing an input
		inputInFocus := m.focusIndex < len(m.inputs) && m.inputs[m.focusIndex].Focused()

		// If ESC is pressed while editing, exit edit mode but keep focus on input
		if inputInFocus && msg.String() == "esc" {
			m.inputs[m.focusIndex].Blur()
			return m, nil
		}

		// If an input has focus, don't process navigation keys - let the input handle them
		if inputInFocus {
			var cmd tea.Cmd
			m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
			return m, cmd
		}

		// Navigation only works when no input is in edit mode
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			if m.focusIndex > -1 {
				m.focusCurrent()
			}
		case key.Matches(msg, m.keyMap.Down) || msg.String() == "j":
			m.blurAll()
			m.focusIndex = (m.focusIndex + 1) % m.focusables
			return m, nil

		case key.Matches(msg, m.keyMap.Up) || msg.String() == "k":
			m.blurAll()
			m.focusIndex = (m.focusIndex - 1 + m.focusables) % m.focusables
			return m, nil

		case key.Matches(msg, m.keyMap.Enter):
			// If an input is selected (but not in edit mode), enter edit mode
			if m.focusIndex < len(m.inputs) {
				m.inputs[m.focusIndex].Focus()
				return m, textinput.Blink
			}

			// If button is focused, execute its action
			if m.focusIndex == m.focusables-1 {
				return m, func() tea.Msg {
					return SaveSettingsMsg{
						Host:   m.inputs[0].Value(),
						Port:   m.inputs[1].Value(),
						APIKey: m.inputs[2].Value(),
					}
				}
			}
		}
	}

	return m, nil
}

// Helper to blur all focusable elements
func (m *SettingsModel) blurAll() {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
	m.saveButton.Blur()
}

// Helper to focus current element
func (m *SettingsModel) focusCurrent() {
	if m.focusIndex < len(m.inputs) {
		m.inputs[m.focusIndex].Focus()
	} else {
		m.saveButton.Focus()
	}
}

// Message sent when save button is clicked
type SaveSettingsMsg struct {
	Host   string
	Port   string
	APIKey string
}

func (m SettingsModel) View() string {
	headerView := m.header.View(m.width)

	// Render inputs
	inputsView := ""
	inputLabels := []string{"Host:", "Port:", "API Key:"}
	for i, input := range m.inputs {
		// Visual indication of selection vs actual focus
		selected := m.focusIndex == i
		editing := input.Focused()

		label := lipgloss.NewStyle().
			Bold(selected).
			Foreground(map[bool]lipgloss.Color{
				true:  lipgloss.Color("#25A065"),
				false: lipgloss.Color("#888888"),
			}[selected]).
			Render(inputLabels[i])

		// Style for the input box
		inputBox := input.View()
		if selected && !editing {
			// Selected but not in edit mode - add highlight
			inputBox = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#25A065")).
				Padding(0, 1).
				Render(inputBox)
		} else if editing {
			// Currently being edited
			inputBox = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FF6700")). // Orange for edit mode
				Padding(0, 1).
				Render(inputBox)
		}

		inputsView += fmt.Sprintf("%s\n%s\n\n", label, inputBox)
	}

	// Render button
	buttonStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(map[bool]lipgloss.Color{
			true:  lipgloss.Color("#25A065"),
			false: lipgloss.Color("#AAAAAA"),
		}[m.focusIndex == m.focusables-1])

	buttonView := buttonStyle.Render(m.saveButton.Text)

	// Navigation help
	var navHelp string
	if m.focusIndex < len(m.inputs) && m.inputs[m.focusIndex].Focused() {
		navHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6700")).
			Render("EDIT MODE: Press ESC to exit editing")
	} else {
		navHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("Navigate with ↑↓ or j/k • Press Enter to edit/select")
	}

	content := fmt.Sprintf(
		"%s\n\n%s%s\n\n%s",
		headerView,
		inputsView,
		buttonView,
		navHelp,
	)

	return styles.PageStyle.Width(m.width - 2).Render(content)
}
