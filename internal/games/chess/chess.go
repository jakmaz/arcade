package chess

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

func (m *Model) View() tea.View {
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

	// Center using the global helper
	centered := styles.CenterContent(content, m.width, m.height)

	return tea.NewView(centered)
}

func (m *Model) renderBoard() string {
	var s strings.Builder

	// Top border
	s.WriteString(m.renderTopBorder())

	// Render each rank (row)
	for y := range 8 {
		// Rank number (8 down to 1)
		rankNum := 8 - y
		rankLabel := styles.HelpStyle.Render(fmt.Sprintf(" %d ", rankNum))
		s.WriteString(rankLabel + "│")

		// Render each file (column) in this rank
		for x := range 8 {
			piece := m.board[y][x]
			pieceStr := " "

			// Render piece with appropriate color styling
			if piece != ' ' {
				if strings.ContainsRune("♔♕♖♗♘♙", piece) {
					pieceStr = styles.WhitePieceStyle.Render(string(piece))
				} else {
					pieceStr = styles.BlackPieceStyle.Render(string(piece))
				}
			}

			// Highlight cursor position with cyan color
			if m.cursorX == x && m.cursorY == y {
				// Create cyan highlight style using theme accent color
				highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // Cyan
				pieceStr = highlightStyle.Render(pieceStr)
			}

			// Add piece with padding (3 chars wide per cell)
			fmt.Fprintf(&s, " %s │", pieceStr)
		}
		s.WriteString("\n")

		// Middle border (between ranks) or bottom border
		if y < 7 {
			s.WriteString(m.renderMiddleBorder())
		}
	}

	// Bottom border
	s.WriteString(m.renderBottomBorder())

	// File labels (A-H)
	s.WriteString(m.renderFileLabels())

	return s.String()
}

// renderTopBorder creates the top border of the chess board
// Format:    ┌───┬───┬───┬───┬───┬───┬───┬───┐
func (m *Model) renderTopBorder() string {
	var border strings.Builder
	border.WriteString("   ┌───")
	for range 7 {
		border.WriteString("┬───")
	}
	border.WriteString("┐\n")
	return border.String()
}

// renderMiddleBorder creates the middle borders between ranks
// Format:    ├───┼───┼───┼───┼───┼───┼───┼───┤
func (m *Model) renderMiddleBorder() string {
	var border strings.Builder
	border.WriteString("   ├───")
	for range 7 {
		border.WriteString("┼───")
	}
	border.WriteString("┤\n")
	return border.String()
}

// renderBottomBorder creates the bottom border of the chess board
// Format:    └───┴───┴───┴───┴───┴───┴───┴───┘
func (m *Model) renderBottomBorder() string {
	var border strings.Builder
	border.WriteString("   └───")
	for range 7 {
		border.WriteString("┴───")
	}
	border.WriteString("┘\n")
	return border.String()
}

// renderFileLabels creates the file labels (A-H) at the bottom
// Format:      A   B   C   D   E   F   G   H
func (m *Model) renderFileLabels() string {
	var labels strings.Builder
	labels.WriteString("     ")
	for i := range 8 {
		label := string(rune('A' + i))
		if i > 0 {
			labels.WriteString("   ")
		}
		labels.WriteString(styles.HelpStyle.Render(label))
	}
	return labels.String()
}
