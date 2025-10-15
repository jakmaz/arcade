package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jakmaz/arcade/internal/core"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

type model struct {
	cursor int
	games  []core.Game
	width  int
	height int
}

func NewMenu() model {
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
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.games)-1 {
				m.cursor++
			}
		case "enter":
			selected := m.games[m.cursor]
			return m, tea.Printf("You selected %s!\n", selected.Name)
		}
	}
	return m, nil
}

func (m model) View() string {
	title := styles.TitleStyle.Render("Arcade")

	var items []string
	for i, g := range m.games {
		style := styles.MenuItemStyle
		cursor := " "
		if m.cursor == i {
			style = styles.SelectedItemStyle
			cursor = "> "
		}
		items = append(items, style.Render(cursor+g.Name+" — "+g.Description))
	}

	help := styles.HelpStyle.Render("↑/↓ to move, Enter to select, q to quit")

	// Center everything
	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		strings.Join(items, "\n"),
		help,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
