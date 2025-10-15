package tetris

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

type Model struct {
	board         [20][10]int
	currentPiece  Piece
	nextPiece     Piece
	score         int
	level         int
	lines         int
	gameOver      bool
	width, height int
}

type Piece struct {
	shape [][]int
	x, y  int
	color int
}

type Position struct {
	x, y int
}

func New() *Model {
	m := &Model{
		currentPiece: Piece{
			shape: [][]int{
				{1, 1, 1, 1},
			},
			x:     3,
			y:     0,
			color: 1,
		},
		nextPiece: Piece{
			shape: [][]int{
				{1, 1},
				{1, 1},
			},
			x:     0,
			y:     0,
			color: 2,
		},
		score:    12450,
		level:    3,
		lines:    8,
		gameOver: false,
	}

	for y := 17; y < 20; y++ {
		for x := 0; x < 10; x++ {
			if x < 2 || x > 7 || y == 19 {
				m.board[y][x] = (x+y)%7 + 1
			}
		}
	}

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
	title := styles.TitleStyle.Render("Tetris")

	gameArea := lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderBoard(),
		"  ",
		m.renderSidebar(),
	)

	var status string
	if m.gameOver {
		status = styles.GameOverStyle.Render("Game Over!")
	} else {
		status = styles.SelectedItemStyle.Render("Playing...")
	}

	help := styles.HelpStyle.Render("← → to move, ↓ to drop, ↑ to rotate, Space to hard drop, ESC to return to menu")

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		gameArea,
		"",
		status,
		"",
		help,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m *Model) renderBoard() string {
	colors := []lipgloss.Color{
		"",
		"#ff6b6b", // Red
		"#4ecdc4", // Teal
		"#45b7d1", // Blue
		"#f9ca24", // Yellow
		"#f0932b", // Orange
		"#eb4d4b", // Dark Red
		"#6c5ce7", // Purple
	}

	var board [20][10]int
	copy(board[:], m.board[:])

	for i, row := range m.currentPiece.shape {
		for j, cell := range row {
			if cell != 0 {
				boardY := m.currentPiece.y + i
				boardX := m.currentPiece.x + j
				if boardY >= 0 && boardY < 20 && boardX >= 0 && boardX < 10 {
					board[boardY][boardX] = m.currentPiece.color
				}
			}
		}
	}

	var rows []string

	topBorder := styles.BorderStyle.Render("┌" + strings.Repeat("─", 20) + "┐")
	rows = append(rows, topBorder)

	for y := range 20 {
		var rowContent strings.Builder
		rowContent.WriteString(styles.BorderStyle.Render("│"))

		for x := range 10 {
			cell := board[y][x]
			if cell == 0 {
				rowContent.WriteString("  ")
			} else {
				blockStyle := lipgloss.NewStyle().
					Background(colors[cell]).
					Foreground(colors[cell])
				rowContent.WriteString(blockStyle.Render("██"))
			}
		}
		rowContent.WriteString(styles.BorderStyle.Render("│"))
		rows = append(rows, rowContent.String())
	}

	bottomBorder := styles.BorderStyle.Render("└" + strings.Repeat("─", 20) + "┘")
	rows = append(rows, bottomBorder)

	return strings.Join(rows, "\n")
}

func (m *Model) renderSidebar() string {
	stats := fmt.Sprintf("Score: %d\n\nLevel: %d\n\nLines: %d", m.score, m.level, m.lines)

	nextPieceTitle := styles.SelectedItemStyle.Render("Next:")
	nextPiece := m.renderNextPiece()

	content := lipgloss.JoinVertical(lipgloss.Left,
		stats,
		"",
		nextPieceTitle,
		nextPiece,
	)

	return styles.SidebarStyle.Render(content)
}

func (m *Model) renderNextPiece() string {
	colors := []lipgloss.Color{
		"",
		"#ff6b6b", // Red
		"#4ecdc4", // Teal
		"#45b7d1", // Blue
		"#f9ca24", // Yellow
		"#f0932b", // Orange
		"#eb4d4b", // Dark Red
		"#6c5ce7", // Purple
	}

	var rows []string
	for _, row := range m.nextPiece.shape {
		var rowContent strings.Builder
		for _, cell := range row {
			if cell == 0 {
				rowContent.WriteString("  ")
			} else {
				blockStyle := lipgloss.NewStyle().
					Background(colors[m.nextPiece.color]).
					Foreground(colors[m.nextPiece.color])
				rowContent.WriteString(blockStyle.Render("██"))
			}
		}
		rows = append(rows, rowContent.String())
	}

	return strings.Join(rows, "\n")
}

type tickMsg struct{}
