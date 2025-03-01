package components

import (
	"bubbletea-app/app/config"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ModalModel represents the configuration and state of a modal dialog
type ModalModel struct {
	title            string
	description      string
	isOpen           bool
	onConfirm        func() tea.Cmd
	onCancel         func() tea.Cmd
	buttonFocusIndex int
}

// NewModal creates a new modal with specified configuration
func NewModal(title, description string) ModalModel {
	return ModalModel{
		title:       title,
		description: description,
		isOpen:      false,
		onConfirm:   nil,
		onCancel:    nil,
	}
}

func (m *ModalModel) IsOpen() bool {
	return m.isOpen
}

// Open shows the modal and sets optional confirmation handlers
func (m *ModalModel) Open(title string, description string, onConfirm func() tea.Cmd, onCancel func() tea.Cmd) {
	m.title = title
	m.description = description
	m.isOpen = true
	m.onConfirm = onConfirm
	m.onCancel = onCancel
}

// Close hides the modal
func (m *ModalModel) Close() {
	m.title = ""
	m.description = ""
	m.isOpen = false
	m.onConfirm = nil
	m.onCancel = nil
}

// View renders the modal dialog
func (m ModalModel) View(width, height int) string {
	if !m.isOpen {
		return ""
	}

	modalWidth := width - 10
	modalHeight := height / 2

	modalStyle := lipgloss.NewStyle().
		Width(modalWidth).
		Height(modalHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1).
		Align(lipgloss.Center, lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A9A9A9"))

	buttonsStyle := lipgloss.NewStyle().
		Padding(1, 0)

	buttons := []string{"Confirm", "Cancel"}
	renderedButtons := make([]string, len(buttons))

	for i, button := range buttons {
		buttonStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(1, 1) // More padding

		isFocused := (i == 0 && m.buttonFocusIndex == 0) || (i == 1 && m.buttonFocusIndex == 1)

		if i == 0 { // Confirm button
			if isFocused {
				buttonStyle = buttonStyle.
					Background(lipgloss.Color("#2E8B57")).
					BorderStyle(lipgloss.RoundedBorder()).
					Bold(true)
			} else {
				buttonStyle = buttonStyle.Background(lipgloss.Color("#4CAF50"))
			}
		} else { // Cancel button
			if isFocused {
				buttonStyle = buttonStyle.
					Background(lipgloss.Color("#DC143C")).
					BorderStyle(lipgloss.RoundedBorder()).
					Bold(true)
			} else {
				buttonStyle = buttonStyle.Background(lipgloss.Color("#F44336"))
			}
		}

		renderedButtons[i] = buttonStyle.Render(button)
	}
	modalContent := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		titleStyle.Render(m.title),
		descriptionStyle.Render(m.description),
		buttonsStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				renderedButtons[0],
				" ",
				renderedButtons[1],
			),
		),
	)

	return modalStyle.Render(modalContent)
}

// HandleKey manages modal interactions
func (m *ModalModel) HandleKey(msg tea.KeyMsg, keyMap config.KeyMap) tea.Cmd {
	if !m.isOpen {
		return nil
	}

	switch {
	case key.Matches(msg, keyMap.Left) || msg.String() == "h":
		m.buttonFocusIndex = 0
		return nil

	case key.Matches(msg, keyMap.Right) || msg.String() == "l":
		m.buttonFocusIndex = 1
		return nil

	case key.Matches(msg, keyMap.Enter):
		m.isOpen = false

		if m.buttonFocusIndex == 0 && m.onConfirm != nil {
			// Store the command before closing the modal
			cmd := m.onConfirm()
			return cmd
		}

		if m.buttonFocusIndex == 1 && m.onCancel != nil {
			// Store the command before closing the modal
			cmd := m.onCancel()
			return cmd
		}

		return nil

	case key.Matches(msg, keyMap.Back):
		m.isOpen = false
		if m.onCancel != nil {
			return m.onCancel()
		}
		return nil
	}

	return nil
}
