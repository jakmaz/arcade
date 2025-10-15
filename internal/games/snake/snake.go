package snake

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	board    [][]int
	snake    []Position
	food     Position
	score    int
	gameOver bool
	width    int
	height   int
}

type Position struct {
	x, y int
}

func New() *Model {
	return &Model{
		snake: []Position{{10, 10}},
		food:  Position{15, 15},
		score: 0,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "arrow_up":
			// Move snake up
		case "arrow_down":
			// Move snake down
			// ... other controls
		}
	case tickMsg:
		if !m.gameOver {
			// Update game state
			return m, tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
				return tickMsg{}
			})
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *Model) View() string {
	if m.gameOver {
		return fmt.Sprintf("Game Over! Score: %d\nPress ESC to return to menu", m.score)
	}

	// Render snake game board
	return "🐍 Snake Game\n" + m.renderBoard() + fmt.Sprintf("\nScore: %d", m.score)
}

func (m *Model) renderBoard() string {
	// Game rendering logic
	return "█████████████\n█   🍎      █\n█     ●●●   █\n█████████████"
}

type tickMsg struct{}
