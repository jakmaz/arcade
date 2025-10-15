package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jakmaz/arcade/internal/core"
)

type model struct {
	cursor int
	games  []core.Game
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
	var b strings.Builder

	b.WriteString("🎮 Arcade\n\n")
	for i, g := range m.games {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s — %s\n", cursor, g.Name, g.Description))
	}
	b.WriteString("\n↑/↓ to move, Enter to select, q to quit\n")

	return b.String()
}
