// app/components/button.go
package components

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ButtonModel represents an interactive button
type ButtonModel struct {
	Text    string
	OnClick func() tea.Msg
	OnFocus func()
	focused bool
	keyMap  KeyMap
	width   int
}

// KeyMap defines button-specific keybindings
type KeyMap struct {
	Enter key.Binding
}

// DefaultKeyMap returns default button keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

// NewButtonModel creates a new button
func NewButtonModel(text string, onClick func() tea.Msg) ButtonModel {
	return ButtonModel{
		Text:    text,
		OnClick: onClick,
		keyMap:  DefaultKeyMap(),
		width:   len(text) + 6, // Text + padding
	}
}

// SetOnFocus sets a callback for when button is focused
func (b *ButtonModel) SetOnFocus(onFocus func()) {
	b.OnFocus = onFocus
}

// SetWidth sets the button width
func (b *ButtonModel) SetWidth(width int) {
	b.width = width
}

// Focus sets focus on the button
func (b *ButtonModel) Focus() {
	b.focused = true
	if b.OnFocus != nil {
		b.OnFocus()
	}
}

// Blur removes focus from the button
func (b *ButtonModel) Blur() {
	b.focused = false
}

// IsFocused returns whether the button is focused
func (b ButtonModel) IsFocused() bool {
	return b.focused
}

// Update handles events for the button
func (b ButtonModel) Update(msg tea.Msg) (ButtonModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !b.focused {
			return b, nil
		}

		// Only handle keys when focused
		switch {
		case key.Matches(msg, b.keyMap.Enter):
			if b.OnClick != nil {
				return b, func() tea.Msg { return b.OnClick() }
			}
		}
	}

	return b, nil
}

// View renders the button
func (b ButtonModel) View() string {
	style := lipgloss.NewStyle().
		Width(b.width).
		Align(lipgloss.Center).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder())

	if b.focused {
		style = style.
			BorderForeground(lipgloss.Color("#25A065")).
			Bold(true)
	} else {
		style = style.
			BorderForeground(lipgloss.Color("#AAAAAA"))
	}

	return style.Render(b.Text)
}
