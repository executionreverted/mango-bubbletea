package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/charmbracelet/bubbles/key"
)

// KeymapFileName returns the full path to keymap config file
func KeymapFileName() string {
	return filepath.Join(GetConfigPath(), "keymap.json")
}

// GetConfigPath returns OS-specific config directory path
func GetConfigPath() string {
	configDir := ".config/sleek"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to local directory if home can't be found
		return configDir
	}

	// For Windows, use %APPDATA%
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "sleek")
		}
	}

	// For Linux/macOS use ~/.config/sleek
	return filepath.Join(homeDir, configDir)
}

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
	Esc      key.Binding
	Back     key.Binding
	J        key.Binding
	K        key.Binding
	Ctrl     key.Binding
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
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "escape"),
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
		Ctrl: key.NewBinding(
			key.WithKeys("ctrl"),
			key.WithHelp("ctrl", ""),
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

// LoadKeyMap loads custom keybindings from the config file
func LoadKeyMap() (KeyMap, error) {
	km := DefaultKeyMap()
	configDir := GetConfigPath()
	filename := KeymapFileName()

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return km, err
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return km, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return km, err
	}

	var userKeymap map[string]string
	err = json.Unmarshal(data, &userKeymap)
	if err != nil {
		return km, err
	}

	// Apply the keymaps after successful unmarshaling
	km = ApplySimpleKeyMap(km, userKeymap)

	return km, nil
}

// ApplySimpleKeyMap applies user key mappings to the default keymap
func ApplySimpleKeyMap(defaultMap KeyMap, umap map[string]string) KeyMap {
	// Get reflection value of KeyMap struct
	val := reflect.ValueOf(&defaultMap).Elem()

	// Iterate through the user's key mappings
	for fieldName, keyVal := range umap {
		// Find the corresponding field in the KeyMap struct
		field := val.FieldByName(fieldName)

		if field.IsValid() && field.Type() == reflect.TypeOf(key.Binding{}) {
			// Get the description from the original binding
			origBinding := field.Interface().(key.Binding)
			h := origBinding.Help()

			// Create a new binding with the custom key
			newBinding := key.NewBinding(
				key.WithKeys(keyVal),
				key.WithHelp(keyVal, h.Desc),
			)

			// Set the field to the new binding
			field.Set(reflect.ValueOf(newBinding))
		}
	}

	return defaultMap
}
