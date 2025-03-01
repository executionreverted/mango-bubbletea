package config

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines keybindings for the application
type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Help     key.Binding
	Quit     key.Binding
	Home     key.Binding
	Settings key.Binding
	About    key.Binding
	Enter    key.Binding
	Back     key.Binding
	J        key.Binding
	K        key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Home: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "home page"),
		),
		Settings: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "settings page"),
		),
		About: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "about page"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		J: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "next item"),
		),
		K: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "previous item"),
		),
	}
}

// SaveKeyMap saves custom keybindings to a file
func SaveKeyMap(km KeyMap, filename string) error {
	data, err := json.MarshalIndent(km, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// LoadKeyMap loads custom keybindings from a file
func LoadKeyMap(filename string) (KeyMap, error) {
	km := DefaultKeyMap()

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Save default keybindings
		SaveKeyMap(km, filename)
		return km, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return km, err
	}

	err = json.Unmarshal(data, &km)
	return km, err
}
