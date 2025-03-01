// app/components/input_wrapper.go
package components

import (
	"github.com/charmbracelet/bubbles/textinput"
)

// InputWrapper wraps textinput.Model to implement FocusableItem
type InputWrapper struct {
	Input *textinput.Model
}

// NewInputWrapper creates a new wrapper for textinput
func NewInputWrapper(input *textinput.Model) *InputWrapper {
	return &InputWrapper{
		Input: input,
	}
}

// Focus sets focus on the input
func (w *InputWrapper) Focus() {
	w.Input.Focus()
}

// Blur removes focus from the input
func (w *InputWrapper) Blur() {
	w.Input.Blur()
}

// IsFocused returns whether the input is focused
func (w *InputWrapper) IsFocused() bool {
	return w.Input.Focused()
}
