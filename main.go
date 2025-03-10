package main

import (
	"bubbletea-app/app/components"
	"bubbletea-app/app/config"
	"bubbletea-app/app/global"
	"bubbletea-app/app/pages"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type (
	appModel struct {
		currentPage      string
		modalModel       components.ModalModel
		homeModel        pages.HomeModel
		settingsModel    pages.SettingsModel
		aboutModel       pages.AboutModel
		keyMap           config.KeyMap
		width            int
		height           int
		showHelp         bool
		activeNavigation bool
		inputInFocus     bool // Track if any input has focus globally
	}
)

func initialModel() appModel {
	// Load keybindings
	keyMap, err := config.LoadKeyMap()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading keymap configuration: %v\n", err)
		os.Exit(1)
	}
	width, height, _ := term.GetSize(0)

	return appModel{
		currentPage:   "home",
		modalModel:    components.NewModal("", ""),
		homeModel:     pages.NewHomeModel(keyMap),
		settingsModel: pages.NewSettingsModel(keyMap),
		aboutModel:    pages.NewAboutModel(),
		keyMap:        keyMap,
		showHelp:      false,
		width:         width,
		height:        height,
	}
}

func (m appModel) Init() tea.Cmd {
	return nil
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	// Modal stuff
	case global.KillModalMsg:
		m.modalModel.Close()
		return m, nil
	case global.SpawnModalMsg:
		m.modalModel.Open(
			msg.Title, msg.Description,
			func() tea.Cmd {
				if msg.OnConfirm != nil {
					return func() tea.Msg { return msg.OnConfirm() }
				}
				return nil
			},
			func() tea.Cmd {
				if msg.OnCancel != nil {
					return func() tea.Msg { return msg.OnCancel() }
				}
				return nil
			},
		)
		return m, nil

	case global.InputFocusChangedMsg:
		// Update the global input focus state
		m.inputInFocus = bool(msg)
		return m, nil

	case tea.KeyMsg:
		if m.modalModel.IsOpen() {

			cmd := m.modalModel.HandleKey(msg, m.keyMap)
			return m, cmd
		}

		// If input is in focus, all keypresses should go to the page
		if m.inputInFocus {
			// Updates to focused inputs should be handled by the page
			switch m.currentPage {
			case "home":
				homeModel, homeCmd := m.homeModel.Update(msg)
				m.homeModel = homeModel.(pages.HomeModel)
				return m, homeCmd
			case "settings":
				settingsModel, settingsCmd := m.settingsModel.Update(msg)
				m.settingsModel = settingsModel.(pages.SettingsModel)
				return m, settingsCmd
			case "about":
				aboutModel, aboutCmd := m.aboutModel.Update(msg)
				m.aboutModel = aboutModel.(pages.AboutModel)
				return m, aboutCmd
			}
			return m, nil
		}

		// Toggle help menu when no input is focused
		if key.Matches(msg, m.keyMap.Help) {
			m.showHelp = !m.showHelp
			return m, nil
		}

		// First check if help is shown
		if m.showHelp {
			if key.Matches(msg, m.keyMap.Help) || key.Matches(msg, m.keyMap.Quit) ||
				key.Matches(msg, m.keyMap.Enter) || key.Matches(msg, m.keyMap.Back) {
				m.showHelp = false
			}
			return m, nil
		}

		// Handle other keybindings when no input has focus, maybe do it better later
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Home):
			m.currentPage = "home"
			return m, func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
			}
		case key.Matches(msg, m.keyMap.Settings):
			m.currentPage = "settings"
			return m, func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
			}
		case key.Matches(msg, m.keyMap.About):
			m.currentPage = "about"
			return m, func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.width, Height: m.height}
			}
		}
	}

	// Update the current page model
	switch m.currentPage {
	case "home":
		homeModel, homeCmd := m.homeModel.Update(msg)
		m.homeModel = homeModel.(pages.HomeModel)
		if homeCmd != nil {
			cmds = append(cmds, homeCmd)
		}
	case "settings":
		settingsModel, settingsCmd := m.settingsModel.Update(msg)
		m.settingsModel = settingsModel.(pages.SettingsModel)
		if settingsCmd != nil {
			cmds = append(cmds, settingsCmd)
		}
	case "about":
		aboutModel, aboutCmd := m.aboutModel.Update(msg)
		m.aboutModel = aboutModel.(pages.AboutModel)
		if aboutCmd != nil {
			cmds = append(cmds, aboutCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m appModel) View() string {
	// Navigation header with keybinding info
	navText := fmt.Sprintf(
		"%s: Home • %s: Settings • %s: About • %s: Quit • %s: Help",
		m.keyMap.Home.Help().Key,
		m.keyMap.Settings.Help().Key,
		m.keyMap.About.Help().Key,
		m.keyMap.Quit.Help().Key,
		m.keyMap.Help.Help().Key,
	)

	nav := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#2F4858")).
		Padding(0, 1).
		Width(m.width - 2).
		Render(navText)

	// Show help screen if toggled
	if m.showHelp {
		helpContent := "KEYBOARD SHORTCUTS\n\n"
		helpContent += fmt.Sprintf("%-15s %s\n", m.keyMap.Home.Help().Key, "Go to Home page")
		helpContent += fmt.Sprintf("%-15s %s\n", m.keyMap.Settings.Help().Key, "Go to Settings page")
		helpContent += fmt.Sprintf("%-15s %s\n", m.keyMap.About.Help().Key, "Go to About page")
		helpContent += fmt.Sprintf("%-15s %s\n", m.keyMap.Up.Help().Key, "Move up")
		helpContent += fmt.Sprintf("%-15s %s\n", m.keyMap.Down.Help().Key, "Move down")
		helpContent += fmt.Sprintf("%-15s %s\n", m.keyMap.Enter.Help().Key, "Select/Confirm")
		helpContent += fmt.Sprintf("%-15s %s\n", m.keyMap.Back.Help().Key, "Go back")
		helpContent += fmt.Sprintf("%-15s %s\n", m.keyMap.Help.Help().Key, "Show/hide help")
		helpContent += fmt.Sprintf("%-15s %s\n", m.keyMap.Quit.Help().Key, "Quit application")
		helpContent += fmt.Sprintf("%-15s", "")
		helpContent += fmt.Sprintf("%-15s", "check github to learn about custom keymaps")

		helpBox := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1).
			Width(m.width - 4).
			Height(m.height).
			Render(helpContent)

		return lipgloss.JoinVertical(lipgloss.Left, nav, helpBox)
	}

	// Content based on current page
	var content string
	switch m.currentPage {
	case "home":
		content = m.homeModel.View()
	case "settings":
		content = m.settingsModel.View()
	case "about":
		content = m.aboutModel.View()
	}

	styled := lipgloss.NewStyle().Padding(1).Render(content)

	// Create the full UI by joining the navigation and content vertically
	fullView := lipgloss.JoinVertical(lipgloss.Left, nav, styled)

	// If modal is open, overlay it on top of the existing content
	if m.modalModel.IsOpen() {
		// Get the modal content with proper dimensions
		modalView := m.modalModel.View(m.width, m.height)

		// Create a semi-transparent overlay
		overlay := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Background(lipgloss.Color("#000000")).
			Foreground(lipgloss.Color("#FFFFFF")).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render("")

		// Center the modal in the viewport
		centeredModal := lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			modalView,
		)

		// Use Z-indexing (string concatenation) to layer views
		return fullView + "\n" + overlay + "\n" + centeredModal

	}
	return fullView
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
