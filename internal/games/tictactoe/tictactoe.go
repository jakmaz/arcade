package tictactoe

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

func New() *Model {
	m := &Model{
		board: [3][3]rune{
			{'X', ' ', 'O'},
			{' ', 'X', ' '},
			{'O', ' ', 'X'},
		},
		turn:    'X',
		cursorX: 1,
		cursorY: 1,
	}
	return m
}

type Model struct {
	board            [3][3]rune
	turn             rune
	cursorX, cursorY int
	width, height    int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	title := styles.TitleStyle.Render("Tic-Tac-Toe")

	board := m.renderBoard()

	currentPlayer := styles.SelectedItemStyle.Render("Current Player: " + string(m.turn))

	help := styles.HelpStyle.Render("↑ ↓ ← → to move, Enter to place, ESC to return to menu")

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

func (m Model) renderBoard() string {
	xStyle := lipgloss.NewStyle().
		Foreground(styles.Success).
		Bold(true)

	oStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff6b6b")).
		Bold(true)

	var rows []string

	for y := range 3 {
		var cells []string
		for x := range 3 {
			cellContent := " "
			switch m.board[y][x] {
			case 'X':
				cellContent = xStyle.Render("✕")
			case 'O':
				cellContent = oStyle.Render("○")
			}

			style := styles.CellStyle
			if m.cursorX == x && m.cursorY == y {
				style = styles.SelectedCellStyle
			}

			cells = append(cells, style.Render(cellContent))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return strings.Join(rows, "\n")
}
