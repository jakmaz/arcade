package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jakmaz/arcade/internal/core"
	"github.com/jakmaz/arcade/internal/theme"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

type model struct {
	cursor int
	games  []core.GameInfo
	width  int
	height int
}

func NewMenu() model {
	theme.Initialize()

	return model{
		games: core.AvailableGames(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor == 0 {
				m.cursor = len(m.games) - 1
			} else {
				m.cursor--
			}
		case "down", "j":
			if m.cursor == len(m.games)-1 {
				m.cursor = 0
			} else {
				m.cursor++
			}
		case "left", "h":
			return m, m.cycleToPreviousTheme()
		case "right", "l":
			return m, m.cycleToNextTheme()
		case "enter":
			selected := m.games[m.cursor] // Direct access to GameInfo
			return m, func() tea.Msg {
				return StartGameMsg{GameID: selected.ID}
			}
		}
	}
	return m, nil
}

func (m model) cycleToPreviousTheme() tea.Cmd {
	return func() tea.Msg {
		availableThemes := theme.ListThemes()
		currentName := theme.GetCurrentTheme().Name()

		// Find current theme index
		currentIndex := -1
		for i, name := range availableThemes {
			if name == currentName {
				currentIndex = i
				break
			}
		}

		// Get previous theme (wrap around)
		prevIndex := (currentIndex - 1 + len(availableThemes)) % len(availableThemes)
		newTheme := availableThemes[prevIndex]

		// Set the new theme
		theme.SetCurrentTheme(newTheme)

		return ThemeChangedMsg{ThemeName: newTheme}
	}
}

func (m model) cycleToNextTheme() tea.Cmd {
	return func() tea.Msg {
		availableThemes := theme.ListThemes()
		currentName := theme.GetCurrentTheme().Name()

		// Find current theme index
		currentIndex := -1
		for i, name := range availableThemes {
			if name == currentName {
				currentIndex = i
				break
			}
		}

		// Get next theme (wrap around)
		nextIndex := (currentIndex + 1) % len(availableThemes)
		newTheme := availableThemes[nextIndex]

		// Set the new theme
		theme.SetCurrentTheme(newTheme)

		return ThemeChangedMsg{ThemeName: newTheme}
	}
}

// ThemeChangedMsg indicates a theme was changed
type ThemeChangedMsg struct {
	ThemeName string
}

func (m model) View() string {
	asciiArt := "                             _      \n" +
		"     /\\                     | |     \n" +
		"    /  \\   _ __ ___ __ _  __| | ___ \n" +
		"   / /\\ \\ | '__/ __/ _` |/ _` |/ _ \\\n" +
		"  / ____ \\| | | (_| (_| | (_| |  __/\n" +
		" /_/    \\_\\_|  \\___\\__,_|\\__,_|\\___|"

	title := styles.GetTitleStyle().Render(asciiArt)

	var items []string
	// Render games with cursor
	for i, game := range m.games {
		style := styles.GetMenuItemStyle()
		cursor := " "
		if m.cursor == i {
			style = styles.GetSelectedItemStyle()
			cursor = "> "
		}
		items = append(items, style.Render(cursor+game.Name+" — "+game.Description))
	}

	// Add separator (just UI)
	items = append(items, styles.GetMenuItemStyle().Render("──────────────────────────────────────"))

	// Add theme display (just UI)
	currentTheme := theme.GetCurrentTheme()
	themeDisplay := fmt.Sprintf(" Theme: ← %s → ", currentTheme.Name())
	items = append(items, styles.GetMenuItemStyle().Render(themeDisplay))
	help := styles.GetHelpStyle().Render("↑/↓ to move, ←/→ to change theme, Enter to select, q to quit")

	// Center everything
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		strings.Join(items, "\n"),
		help,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
