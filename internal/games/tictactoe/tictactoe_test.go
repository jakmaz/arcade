package tictactoe

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/jakmaz/arcade/internal/theme"
)

func localGame() *Model {
	m := New()
	m.computer = false
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	return m
}

func TestAllWinLines(t *testing.T) {
	lines := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}
	for _, player := range []rune{'X', 'O'} {
		for _, line := range lines {
			t.Run(fmt.Sprintf("%c/%v", player, line), func(t *testing.T) {
				m := localGame()
				m.turn = player
				m.board[line[0]], m.board[line[1]] = player, player
				if !m.place(line[2]) || m.winner != player || m.draw || !m.finished() {
					t.Fatalf("winning move was not recognized: %+v", m)
				}
				before := m.board
				for cell := range m.board {
					if m.place(cell) {
						t.Fatal("move was accepted after a win")
					}
				}
				if m.board != before || m.turn != player {
					t.Fatal("finished board or final turn changed")
				}
			})
		}
	}
}

func TestAlternationAndIllegalMoves(t *testing.T) {
	m := localGame()
	if !m.place(0) || m.board[0] != 'X' || m.turn != 'O' {
		t.Fatal("X's move did not pass the turn to O")
	}
	before := m.board
	for _, cell := range []int{-1, 0, 9, 100} {
		if m.place(cell) || m.board != before || m.turn != 'O' {
			t.Fatalf("illegal move %d changed the board or turn", cell)
		}
	}
	if !m.place(4) || m.board[4] != 'O' || m.turn != 'X' || m.finished() {
		t.Fatal("O's move did not pass the turn to X")
	}
}

func TestDrawAndFinalSquareWin(t *testing.T) {
	t.Run("draw", func(t *testing.T) {
		m := localGame()
		for move, cell := range []int{0, 1, 2, 4, 3, 5, 7, 6, 8} {
			if !m.place(cell) {
				t.Fatalf("legal move %d was rejected", move)
			}
			if move < 8 && m.finished() {
				t.Fatalf("game finished early at move %d", move)
			}
		}
		if !m.draw || m.winner != 0 || !m.finished() {
			t.Fatal("full board was not a draw")
		}
		before := m.board
		if m.place(0) || m.board != before {
			t.Fatal("move was accepted after a draw")
		}
	})
	t.Run("win on ninth move", func(t *testing.T) {
		m := localGame()
		m.board = [9]rune{'X', 'O', 'X', 'O', 'X', 'O', 'O', 'X', 0}
		if !m.place(8) || m.winner != 'X' || m.draw {
			t.Fatal("a win on the last square was treated as a draw")
		}
	})
}

func TestKeyboardInput(t *testing.T) {
	for _, tc := range []struct {
		name   string
		start  int
		key    rune
		finish int
	}{
		{"up", 4, tea.KeyUp, 1}, {"down", 4, tea.KeyDown, 7},
		{"left", 4, tea.KeyLeft, 3}, {"right", 4, tea.KeyRight, 5},
		{"top edge", 0, tea.KeyUp, 0}, {"bottom edge", 8, tea.KeyDown, 8},
		{"left edge", 3, tea.KeyLeft, 3}, {"right edge", 5, tea.KeyRight, 5},
		{"ignored letter", 4, 'w', 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := localGame()
			m.cursor = tc.start
			m.Update(tea.KeyPressMsg{Code: tc.key})
			if m.cursor != tc.finish || m.board != ([9]rune{}) {
				t.Fatalf("unexpected cursor or board: %+v", m)
			}
		})
	}
	for _, key := range []rune{tea.KeyEnter, tea.KeySpace} {
		m := localGame()
		m.Update(tea.KeyPressMsg{Code: key})
		if m.board[4] != 'X' || m.turn != 'O' {
			t.Fatalf("placement key %v did not place X", key)
		}
		m.Update(tea.KeyPressMsg{Code: key})
		if m.board[4] != 'X' || m.turn != 'O' || m.notice == "" {
			t.Fatal("occupied square was overwritten or gave no feedback")
		}
		m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
		if m.notice != "" {
			t.Fatal("occupied-square feedback did not clear on navigation")
		}
		m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
		if m.cursor != 4 {
			t.Fatal("cursor could not return to an occupied cell")
		}
	}
	for key := '1'; key <= '9'; key++ {
		m := localGame()
		m.Update(tea.KeyPressMsg{Code: key})
		cell := int(key - '1')
		if m.cursor != cell || m.board[cell] != 'X' || m.turn != 'O' {
			t.Fatalf("shortcut %c did not use row-major order", key)
		}
	}
}

func TestRestartAndModeSwitch(t *testing.T) {
	m := New()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if !m.computer || m.board != ([9]rune{}) || m.turn != 'X' || m.Init() != nil {
		t.Fatal("new game is not an empty synchronous solo game")
	}
	_, command := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if command != nil || m.board[4] != 'X' || m.turn != 'X' {
		t.Fatal("computer did not respond synchronously")
	}
	_, command = m.Update(tea.KeyPressMsg{Code: 'r'})
	if command != nil || m.board != ([9]rune{}) || !m.computer || m.turn != 'X' {
		t.Fatal("restart left a move or pending computer response")
	}
	m.Update(tea.KeyPressMsg{Code: 'c'})
	if m.computer {
		t.Fatal("C did not select local mode")
	}
	m.Update(tea.KeyPressMsg{Code: '1'})
	if m.board[0] != 'X' || m.turn != 'O' {
		t.Fatal("local mode did not allow alternating players")
	}
	m.winner, m.draw, m.notice = 'X', true, "occupied"
	m.Update(tea.KeyPressMsg{Code: 'R'})
	if m.board != ([9]rune{}) || m.finished() || m.notice != "" || m.cursor != 4 || m.turn != 'X' || m.computer {
		t.Fatal("restart did not clear the round and preserve local mode")
	}
	if m.width != 80 || m.height != 24 {
		t.Fatal("restart lost the terminal dimensions")
	}
	m.Update(tea.KeyPressMsg{Code: '1'})
	m.Update(tea.KeyPressMsg{Code: 'C'})
	if !m.computer || m.board != ([9]rune{}) || m.turn != 'X' {
		t.Fatal("changing to solo mode did not start a fresh round")
	}
}

func TestComputerTactics(t *testing.T) {
	for _, tc := range []struct {
		name  string
		board [9]rune
		want  int
	}{
		{"wins before blocking", [9]rune{'O', 'O', 0, 'X', 'X', 0, 'X', 0, 0}, 2},
		{"blocks a row", [9]rune{'X', 'X', 0, 0, 'O', 0, 0, 0, 0}, 2},
		{"blocks a diagonal", [9]rune{'X', 'O', 0, 0, 'X', 0, 0, 0, 0}, 8},
		{"takes center", [9]rune{'X', 0, 0, 0, 0, 0, 0, 0, 0}, 4},
		{"finished win", [9]rune{'X', 'X', 'X', 0, 'O', 0, 0, 0, 'O'}, -1},
		{"finished draw", [9]rune{'X', 'O', 'X', 'X', 'O', 'O', 'O', 'X', 'X'}, -1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := tc.board
			if got := computerMove(tc.board); got != tc.want {
				t.Fatalf("computer chose %d, want %d", got, tc.want)
			}
			if tc.board != before {
				t.Fatal("computer search mutated the caller's board")
			}
		})
	}
}

func TestComputerCannotLose(t *testing.T) {
	var explore func(Model)
	explore = func(m Model) {
		if m.winner == 'X' {
			t.Fatalf("computer lost: %v", m.board)
		}
		if m.finished() {
			return
		}
		if m.turn != 'X' {
			t.Fatal("computer did not return the turn to the human")
		}
		for cell, mark := range m.board {
			if mark != 0 {
				continue
			}
			next := m
			next.cursor = cell
			next.play()
			if next.board[cell] != 'X' {
				t.Fatal("computer overwrote the human move")
			}
			xMoves, oMoves := 0, 0
			for i, before := range m.board {
				if before != 0 && next.board[i] != before {
					t.Fatal("computer overwrote an occupied square")
				}
				if before == 0 {
					switch next.board[i] {
					case 'X':
						xMoves++
					case 'O':
						oMoves++
					}
				}
			}
			if xMoves != 1 || oMoves > 1 || (!next.finished() && oMoves != 1) {
				t.Fatal("computer made an illegal number of moves")
			}
			explore(next)
		}
	}
	explore(*New())
}

func TestBoardDimensionsAndOccupiedSelection(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	m := localGame()
	m.board = [9]rune{'X', 'O', 'X', 'O', 0, 'O', 'X', 'O', 'X'}
	for cell := range m.board {
		m.cursor = cell
		board := ansi.Strip(m.renderBoard())
		rows := strings.Split(board, "\n")
		if len(rows) != boardHeight {
			t.Fatalf("board has %d rows, want %d", len(rows), boardHeight)
		}
		for _, row := range rows {
			if got := lipgloss.Width(row); got != boardWidth {
				t.Fatalf("board row has %d columns, want %d", got, boardWidth)
			}
		}
		left, top := (cell%3)*(cellWidth+1), (cell/3)*(cellHeight+1)
		if []rune(rows[top])[left] != '┏' || []rune(rows[top+cellHeight+1])[left] != '┗' {
			t.Fatalf("cell %d has no visible selection outline", cell)
		}
	}
}

func TestViewFitsTerminal(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	for _, size := range [][2]int{
		{80, 24}, {viewWidth, viewHeight}, {viewWidth - 1, viewHeight},
		{viewWidth, viewHeight - 1}, {20, 8}, {8, 2}, {1, 1}, {0, 0}, {-1, -1},
	} {
		t.Run(fmt.Sprintf("%dx%d", size[0], size[1]), func(t *testing.T) {
			m := New()
			m.Update(tea.WindowSizeMsg{Width: size[0], Height: size[1]})
			content := ansi.Strip(m.View().Content)
			if m.width == 0 || m.height == 0 {
				if content != "" {
					t.Fatal("zero-sized terminal has content")
				}
				return
			}
			if len(strings.Split(content, "\n")) > m.height {
				t.Fatal("view is taller than the terminal")
			}
			for _, line := range strings.Split(content, "\n") {
				if lipgloss.Width(line) > m.width {
					t.Fatal("view is wider than the terminal")
				}
			}
			if m.width >= viewWidth && m.height >= viewHeight {
				if !strings.Contains(content, "SOLO · YOU X") || !strings.Contains(content, "Your turn · X") || strings.Contains(content, "Preview") {
					t.Fatal("initial view does not describe a fresh playable game")
				}
			} else if m.width >= 18 && !strings.Contains(content, "Terminal too small") {
				t.Fatal("small terminal feedback is missing")
			}
		})
	}
	if theme.CurrentTheme() == nil {
		t.Fatal("rendering a direct launch did not initialize the theme")
	}
}

func TestSmallTerminalDoesNotAcceptHiddenMoves(t *testing.T) {
	m := localGame()
	m.Update(tea.WindowSizeMsg{Width: viewWidth - 1, Height: viewHeight})
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m.Update(tea.KeyPressMsg{Code: '1'})
	if m.board != ([9]rune{}) {
		t.Fatal("an invisible move was accepted in a small terminal")
	}
	m.Update(tea.WindowSizeMsg{Width: viewWidth, Height: viewHeight})
	m.Update(tea.KeyPressMsg{Code: '1'})
	if m.board[0] != 'X' {
		t.Fatal("game did not resume accepting moves after a resize")
	}
}

func TestStatusAndModeLabels(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	for _, tc := range []struct {
		name     string
		computer bool
		turn     rune
		winner   rune
		draw     bool
		want     string
	}{
		{"local X turn", false, 'X', 0, false, "X's turn"},
		{"local O turn", false, 'O', 0, false, "O's turn"},
		{"local X win", false, 'X', 'X', false, "X wins"},
		{"local O win", false, 'O', 'O', false, "O wins"},
		{"solo turn", true, 'X', 0, false, "Your turn · X"},
		{"human win", true, 'X', 'X', false, "You win"},
		{"computer win", true, 'O', 'O', false, "Computer wins"},
		{"draw", true, 'X', 0, true, "Draw"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := localGame()
			m.width, m.height = viewWidth, viewHeight
			m.computer, m.turn, m.winner, m.draw = tc.computer, tc.turn, tc.winner, tc.draw
			content := ansi.Strip(m.View().Content)
			if !strings.Contains(content, tc.want) {
				t.Fatalf("status %q is missing", tc.want)
			}
			label, help := "LOCAL · X / O", "C solo"
			if tc.computer {
				label, help = "SOLO · YOU X", "C local"
			}
			if !strings.Contains(content, label) || !strings.Contains(content, help) {
				t.Fatal("mode or toggle destination is not visible")
			}
			for _, line := range strings.Split(content, "\n") {
				if lipgloss.Width(line) > m.width {
					t.Fatal("status or help exceeds the minimum terminal width")
				}
			}
		})
	}
}
