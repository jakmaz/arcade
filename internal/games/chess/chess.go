package chess

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

type Model struct {
	board            [8][8]rune
	currentPlayer    string
	cursorX, cursorY int
	width, height    int
}

type Position struct {
	x, y int
}

func New() *Model {
	m := &Model{
		board: [8][8]rune{
			{'♜', '♞', '♝', '♛', '♚', '♝', '♞', '♜'},
			{'♟', '♟', '♟', '♟', '♟', '♟', '♟', '♟'},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			{'♙', '♙', '♙', '♙', '♙', '♙', '♙', '♙'},
			{'♖', '♘', '♗', '♕', '♔', '♗', '♘', '♖'},
		},
		currentPlayer: "White",
		cursorX:       0,
		cursorY:       0,
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
	title := styles.TitleStyle.Render("Chess")

	board := m.renderBoard()

	currentPlayer := styles.SelectedItemStyle.Render("Current Player: " + m.currentPlayer)

	help := styles.HelpStyle.Render("↑ ↓ ← → to move, Enter to select, ESC to return to menu")

	content := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		board,
		"",
		currentPlayer,
		"",
		help,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m *Model) renderBoard() string {
	var rows []string

	for y := range 8 {
		var cells []string
		for x := range 8 {
			piece := m.board[y][x]
			cellContent := " "

			if piece != ' ' {
				if strings.ContainsRune("♔♕♖♗♘♙", piece) {
					cellContent = styles.WhitePieceStyle.Render(string(piece))
				} else {
					cellContent = styles.BlackPieceStyle.Render(string(piece))
				}
			}

			style := styles.CellStyle
			if (x+y)%2 == 1 {
				style = style.Background(lipgloss.Color("#2a2a2a"))
			}
			if m.cursorX == x && m.cursorY == y {
				style = styles.SelectedCellStyle
			}

			cells = append(cells, style.Render(cellContent))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return strings.Join(rows, "\n")
}

type tickMsg struct{}
