package global

import tea "github.com/charmbracelet/bubbletea"

type SpawnModalMsg struct {
	Title       string
	Description string
	OnConfirm   func() tea.Msg
	OnCancel    func() tea.Msg
}

type KillModalMsg bool

type (
	InputFocusChangedMsg bool
)
