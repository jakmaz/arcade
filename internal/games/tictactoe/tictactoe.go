package tictactoe

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jakmaz/arcade/internal/theme"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

const (
	cellWidth   = 8
	cellHeight  = 3
	boardWidth  = 3*cellWidth + 4
	boardHeight = 3*cellHeight + 4
	viewWidth   = 42
	viewHeight  = boardHeight + 6
)

var winLines = [8][3]int{
	{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
	{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
	{0, 4, 8}, {2, 4, 6},
}

type Model struct {
	board         [9]rune
	turn          rune
	cursor        int
	winner        rune
	draw          bool
	computer      bool
	notice        string
	width, height int
}

func New() *Model {
	m := &Model{computer: true}
	m.reset()
	return m
}

func (m *Model) reset() {
	m.board = [9]rune{}
	m.turn = 'X'
	m.cursor = 4
	m.winner = 0
	m.draw = false
	m.notice = ""
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = max(0, msg.Width), max(0, msg.Height)
	case tea.KeyPressMsg:
		key := msg.String()
		switch key {
		case "r", "R":
			m.reset()
			return m, nil
		case "c", "C":
			m.computer = !m.computer
			m.reset()
			return m, nil
		}
		if m.width < viewWidth || m.height < viewHeight {
			return m, nil
		}
		switch key {
		case "up":
			if m.cursor >= 3 {
				m.cursor -= 3
			}
			m.notice = ""
		case "down":
			if m.cursor < 6 {
				m.cursor += 3
			}
			m.notice = ""
		case "left":
			if m.cursor%3 > 0 {
				m.cursor--
			}
			m.notice = ""
		case "right":
			if m.cursor%3 < 2 {
				m.cursor++
			}
			m.notice = ""
		case "enter", "space":
			m.play()
		default:
			if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
				m.cursor = int(key[0] - '1')
				m.play()
			}
		}
	}
	return m, nil
}

func (m *Model) finished() bool {
	return m.winner != 0 || m.draw
}

// place applies one legal move. The turn stays on the final player after a win or draw.
func (m *Model) place(cell int) bool {
	if cell < 0 || cell >= len(m.board) || m.finished() || m.board[cell] != 0 {
		return false
	}
	m.board[cell] = m.turn
	m.winner, m.draw = outcome(m.board)
	if !m.finished() {
		if m.turn == 'X' {
			m.turn = 'O'
		} else {
			m.turn = 'X'
		}
	}
	return true
}

func (m *Model) play() {
	if m.finished() {
		return
	}
	if !m.place(m.cursor) {
		m.notice = "Cell occupied · choose another"
		return
	}
	m.notice = ""
	if m.computer && !m.finished() && m.turn == 'O' {
		m.place(computerMove(m.board))
	}
}

func outcome(board [9]rune) (winner rune, draw bool) {
	for _, line := range winLines {
		mark := board[line[0]]
		if mark != 0 && mark == board[line[1]] && mark == board[line[2]] {
			return mark, false
		}
	}
	for _, mark := range board {
		if mark == 0 {
			return 0, false
		}
	}
	return 0, true
}

// computerMove searches the small game tree synchronously, preferring quick wins.
func computerMove(board [9]rune) int {
	if winner, draw := outcome(board); winner != 0 || draw {
		return -1
	}
	bestCell, bestScore := -1, -100
	for _, cell := range [...]int{4, 0, 2, 6, 8, 1, 3, 5, 7} {
		if board[cell] != 0 {
			continue
		}
		board[cell] = 'O'
		score := -minimax(board, 'X', 1)
		board[cell] = 0
		if score > bestScore {
			bestCell, bestScore = cell, score
		}
	}
	return bestCell
}

func minimax(board [9]rune, player rune, depth int) int {
	winner, draw := outcome(board)
	if winner != 0 {
		if winner == player {
			return 10 - depth
		}
		return depth - 10
	}
	if draw {
		return 0
	}
	next := 'X'
	if player == 'X' {
		next = 'O'
	}
	best := -100
	for cell, mark := range board {
		if mark != 0 {
			continue
		}
		board[cell] = player
		best = max(best, -minimax(board, next, depth+1))
		board[cell] = 0
	}
	return best
}

func (m *Model) View() tea.View {
	styles.GetStyles()
	if m.width == 0 || m.height == 0 {
		return tea.NewView("")
	}
	if m.width < viewWidth || m.height < viewHeight {
		message := fmt.Sprintf("Terminal too small\nNeed %d × %d", viewWidth, viewHeight)
		message = styles.GetHelpStyle().MarginTop(0).
			MaxWidth(m.width).MaxHeight(m.height).Render(message)
		return tea.NewView(styles.CenterContent(message, m.width, m.height))
	}

	title := styles.GetTitleStyle().MarginBottom(0).Render("TIC TAC TOE")
	mode := "SOLO · YOU X"
	if !m.computer {
		mode = "LOCAL · X / O"
	}
	mode = styles.GetSelectedItemStyle().Render(mode)
	header := title + strings.Repeat(" ", viewWidth-lipgloss.Width(title)-lipgloss.Width(mode)) + mode

	status := fmt.Sprintf("%c's turn", m.turn)
	statusStyle := styles.GetStyles().Player1Style()
	if m.turn == 'O' {
		statusStyle = styles.GetStyles().Player2Style()
	}
	if m.computer {
		status = "Your turn · X"
	}
	switch {
	case m.winner != 0:
		status = fmt.Sprintf("%c wins · R to play again", m.winner)
		if m.computer {
			status = "You win · R to play again"
			if m.winner == 'O' {
				status = "Computer wins · R to play again"
			}
		}
	case m.draw:
		status = "Draw · R to play again"
		statusStyle = styles.GetHelpStyle().MarginTop(0)
	case m.notice != "":
		status = m.notice
		statusStyle = styles.GetStyles().WarningStyle()
	}

	help := styles.GetHelpStyle().MarginTop(0)
	modeHelp := "C local"
	if !m.computer {
		modeHelp = "C solo"
	}
	content := strings.Join([]string{
		header,
		lipgloss.PlaceHorizontal(viewWidth, lipgloss.Center, m.renderBoard()),
		"",
		statusStyle.Render(status),
		help.Render("↑↓←→ select · Enter/Space place"),
		help.Render("1–9 place · " + modeHelp + " · R restart · Esc back"),
	}, "\n")
	return tea.NewView(styles.CenterContent(content, m.width, m.height))
}

func (m *Model) renderBoard() string {
	themed := styles.GetStyles()
	border := styles.GetBorderStyle()
	selected := lipgloss.NewStyle().Foreground(theme.CurrentTheme().SelectedCell()).Bold(true)
	xArt := [cellHeight][]rune{[]rune(" ▀▄  ▄▀ "), []rune("   ██   "), []rune(" ▄▀  ▀▄ ")}
	oArt := [cellHeight][]rune{[]rune(" ▗▄▄▄▄▖ "), []rune(" ▐    ▌ "), []rune(" ▝▀▀▀▀▘ ")}

	horizontal := strings.Repeat("─", cellWidth)
	gridRows := [4][]rune{
		[]rune("┌" + horizontal + "┬" + horizontal + "┬" + horizontal + "┐"),
		[]rune("├" + horizontal + "┼" + horizontal + "┼" + horizontal + "┤"),
		[]rune("└" + horizontal + "┴" + horizontal + "┴" + horizontal + "┘"),
		[]rune("│" + strings.Repeat(strings.Repeat(" ", cellWidth)+"│", 3)),
	}

	var out strings.Builder
	for y := range boardHeight {
		if y > 0 {
			out.WriteByte('\n')
		}
		grid := gridRows[3]
		if y%(cellHeight+1) == 0 {
			grid = gridRows[1]
			if y == 0 {
				grid = gridRows[0]
			} else if y == boardHeight-1 {
				grid = gridRows[2]
			}
		}
		for x, glyph := range grid {
			style := border
			if outline := m.selectionGlyph(x, y); outline != 0 {
				glyph, style = outline, selected
			} else if y%(cellHeight+1) != 0 && x%(cellWidth+1) != 0 {
				cell := (y/(cellHeight+1))*3 + x/(cellWidth+1)
				row, column := y%(cellHeight+1)-1, x%(cellWidth+1)-1
				switch m.board[cell] {
				case 'X':
					glyph, style = xArt[row][column], themed.Player1Style()
				case 'O':
					glyph, style = oArt[row][column], themed.Player2Style()
				}
			}
			out.WriteString(style.Render(string(glyph)))
		}
	}
	return out.String()
}

func (m *Model) selectionGlyph(x, y int) rune {
	left, top := (m.cursor%3)*(cellWidth+1), (m.cursor/3)*(cellHeight+1)
	right, bottom := left+cellWidth+1, top+cellHeight+1
	if x < left || x > right || y < top || y > bottom {
		return 0
	}
	if y == top || y == bottom {
		if x == left {
			if y == top {
				return '┏'
			}
			return '┗'
		}
		if x == right {
			if y == top {
				return '┓'
			}
			return '┛'
		}
		return '━'
	}
	if x == left || x == right {
		return '┃'
	}
	return 0
}
