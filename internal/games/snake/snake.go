package snake

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

type Model struct {
	board         [20][30]rune
	snake         []Position
	food          Position
	direction     Direction
	score         int
	gameOver      bool
	width, height int
}

type Position struct {
	x, y int
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func New() *Model {
	m := &Model{
		snake: []Position{
			{15, 10},
			{14, 10},
			{13, 10},
		},
		food:      Position{20, 10},
		direction: Right,
		score:     0,
		gameOver:  false,
	}
	m.updateBoard()
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *Model) View() string {
	title := styles.TitleStyle.Render("Snake")

	board := m.renderBoard()

	var status string
	if m.gameOver {
		status = styles.GameOverStyle.Render(fmt.Sprintf("Game Over! Final Score: %d", m.score))
	} else {
		status = styles.SelectedItemStyle.Render(fmt.Sprintf("Score: %d", m.score))
	}

	help := styles.HelpStyle.Render("↑ ↓ ← → to move, Space to pause, ESC to return to menu")

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		board,
		"",
		status,
		"",
		help,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m *Model) updateBoard() {
	for y := range 20 {
		for x := range 30 {
			m.board[y][x] = ' '
		}
	}

	for _, segment := range m.snake {
		if segment.y >= 0 && segment.y < 20 && segment.x >= 0 && segment.x < 30 {
			m.board[segment.y][segment.x] = '●'
		}
	}

	if len(m.snake) > 0 {
		head := m.snake[0]
		if head.y >= 0 && head.y < 20 && head.x >= 0 && head.x < 30 {
			m.board[head.y][head.x] = '◉'
		}
	}

	if m.food.y >= 0 && m.food.y < 20 && m.food.x >= 0 && m.food.x < 30 {
		m.board[m.food.y][m.food.x] = '🍎'
	}
}

func (m *Model) renderBoard() string {
	foodStyle := lipgloss.NewStyle().
		Bold(true)

	var rows []string

	topBorder := styles.BorderStyle.Render("┌" + strings.Repeat("─", 30) + "┐")
	rows = append(rows, topBorder)

	for y := range 20 {
		var rowContent strings.Builder
		rowContent.WriteString(styles.BorderStyle.Render("│"))

		for x := range 30 {
			cell := m.board[y][x]
			switch cell {
			case '◉':
				rowContent.WriteString(styles.SnakeHeadStyle.Render("◉"))
			case '●':
				rowContent.WriteString(styles.SnakeStyle.Render("●"))
			case '🍎':
				rowContent.WriteString(foodStyle.Render("🍎"))
			default:
				rowContent.WriteString(" ")
			}
		}
		rowContent.WriteString(styles.BorderStyle.Render("│"))
		rows = append(rows, rowContent.String())
	}

	bottomBorder := styles.BorderStyle.Render("└" + strings.Repeat("─", 30) + "┘")
	rows = append(rows, bottomBorder)

	return strings.Join(rows, "\n")
}

type tickMsg struct{}
