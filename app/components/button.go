// app/components/button.go
package components

import tea "github.com/charmbracelet/bubbletea"

// ButtonModel represents a focusable button
type ButtonModel struct {
	Text    string
	focused bool
}

// NewButtonModel creates a new button
func NewButtonModel(text string, onClick func() tea.Cmd) ButtonModel {
	return ButtonModel{
		Text:    text,
		focused: false,
	}
}

// Focus sets focus on the button
func (b *ButtonModel) Focus() {
	b.focused = true
}

// Blur removes focus from the button
func (b *ButtonModel) Blur() {
	b.focused = false
}

// IsFocused returns whether the button is focused
func (b ButtonModel) IsFocused() bool {
	return b.focused
}
