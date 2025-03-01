// app/pages/settings.go
package pages

import (
	"bubbletea-app/app/components"
	"bubbletea-app/app/config"
	"bubbletea-app/app/global"
	"bubbletea-app/app/styles"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsModel struct {
	inputs     []textinput.Model
	buttons    []components.ButtonModel
	header     components.HeaderModel
	footer     components.FooterModel
	keyMap     config.KeyMap
	width      int
	height     int
	focusIndex int // tracks which element has focus
	focusables int // total focusable elements
}

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
	saveButton := components.NewButtonModel("Save Configuration", func() tea.Msg {
		fmt.Println("Config saved:", hostInput.Value(), portInput.Value(), apiKeyInput.Value())
		return nil
	})
	leaveButton := components.NewButtonModel("QUIT!", func() tea.Msg {
		return tea.Quit()
	})

	return SettingsModel{
		inputs:     []textinput.Model{hostInput, portInput, apiKeyInput},
		buttons:    []components.ButtonModel{saveButton, leaveButton},
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle already focused inputs first
		for i, input := range m.inputs {
			if input.Focused() {
				if key.Matches(msg, m.keyMap.Esc) {
					m.inputs[i].Blur()
					return m, func() tea.Msg {
						return global.InputFocusChangedMsg(false)
					}
				}
				var cmd tea.Cmd
				m.inputs[i], cmd = m.inputs[i].Update(msg)
				return m, cmd
			}
		}

		// Navigation when no inputs are focused
		switch {
		case key.Matches(msg, m.keyMap.Down) || msg.String() == "j":
			m.focusIndex = (m.focusIndex + 1) % m.focusables
			updateFocusState(&m)
			return m, nil

		case key.Matches(msg, m.keyMap.Up) || msg.String() == "k":
			m.focusIndex = (m.focusIndex - 1 + m.focusables) % m.focusables
			updateFocusState(&m)
			return m, nil

		case key.Matches(msg, m.keyMap.Left) || msg.String() == "h":
			// Only navigate left when on a button
			if m.focusIndex >= len(m.inputs) {
				buttonIndex := m.focusIndex - len(m.inputs)
				if buttonIndex > 0 {
					m.focusIndex--
					updateFocusState(&m)
				}
			}
			return m, nil

		case key.Matches(msg, m.keyMap.Right) || msg.String() == "l":
			// Only navigate right when on a button
			if m.focusIndex >= len(m.inputs) {
				buttonIndex := m.focusIndex - len(m.inputs)
				if buttonIndex < len(m.buttons)-1 {
					m.focusIndex++
					updateFocusState(&m)
				}
			}
			return m, nil

		case key.Matches(msg, m.keyMap.Enter):
			// Toggle input focus or activate button
			if m.focusIndex < len(m.inputs) {
				m.inputs[m.focusIndex].Focus()
				return m, tea.Batch(
					textinput.Blink,
					func() tea.Msg {
						return global.InputFocusChangedMsg(true)
					},
				)
			} else {
				buttonIndex := m.focusIndex - len(m.inputs)
				if buttonIndex < len(m.buttons) {
					return m, m.buttons[buttonIndex].OnClick
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
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

	// Alternative with more buttons and dynamic spacing:
	buttonViews := make([]string, len(m.buttons))
	for i, button := range m.buttons {
		buttonViews[i] = button.View()
	}

	buttonsContainer := lipgloss.JoinHorizontal(
		lipgloss.Center,
		buttonViews...,
	)

	// Navigation help
	var navHelp string
	if m.focusIndex < len(m.inputs) && m.inputs[m.focusIndex].Focused() {
		navHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6700")).
			Render("EDIT MODE: Press ESC to exit editing or ENTER to submit")
	} else {
		navHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("Navigate with ↑↓ or j/k • Press Enter to edit/select")
	}
	// Footer view
	footerView := m.footer.View(m.width)

	contentHeight := m.height - 4 // Subtract nav and footer height
	content := fmt.Sprintf(
		"%s\n%s%s\n%s\n%s",
		headerView,
		inputsView,
		buttonsContainer,
		navHelp,
		footerView,
	)
	contentContainer := lipgloss.NewStyle().
		Height(contentHeight).Render(content)

	return styles.PageStyle.Width(m.width - 2).Render(contentContainer)
}

func updateFocusState(m *SettingsModel) {
	// Clear focus on all elements
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
	for i := range m.buttons {
		m.buttons[i].Blur()
	}

	// Set focus only on buttons when they're selected
	if m.focusIndex >= len(m.inputs) {
		buttonIndex := m.focusIndex - len(m.inputs)
		if buttonIndex < len(m.buttons) {
			m.buttons[buttonIndex].Focus()
		}
	}
}
