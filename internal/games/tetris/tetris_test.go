package tetris

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jakmaz/arcade/internal/theme"
)

func playableModel() *Model {
	m := New()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	return m
}

func press(m *Model, key rune) tea.Cmd {
	_, cmd := m.Update(tea.KeyPressMsg{Code: key})
	return cmd
}

func occupied(m *Model) int {
	count := 0
	for _, row := range m.board {
		for _, cell := range row {
			if cell != empty {
				count++
			}
		}
	}
	return count
}

func TestNewGameAndSevenBag(t *testing.T) {
	m := New()
	if occupied(m) != 0 || m.score != 0 || m.lines != 0 || m.level != 1 || m.paused || m.gameOver {
		t.Fatalf("new game is not empty: score=%d lines=%d level=%d", m.score, m.lines, m.level)
	}
	if !m.fits(m.currentPiece) || m.currentPiece.kind == m.nextPiece {
		t.Fatal("initial pieces are invalid or not from the same seven-bag")
	}
	pieces := []pieceKind{m.currentPiece.kind, m.nextPiece}
	for range 70 - len(pieces) {
		pieces = append(pieces, m.drawPiece())
	}
	for start := 0; start < len(pieces); start += 7 {
		seen := make(map[pieceKind]bool)
		for _, kind := range pieces[start : start+7] {
			if kind < pieceI || kind > pieceL || seen[kind] {
				t.Fatalf("bag %d contains an invalid or duplicate piece: %v", start/7, pieces[start:start+7])
			}
			seen[kind] = true
		}
	}
}

func TestTetrominoRotations(t *testing.T) {
	for kind := pieceI; kind <= pieceL; kind++ {
		t.Run(pieceNames[kind], func(t *testing.T) {
			m := New()
			m.currentPiece = newPiece(kind)
			m.currentPiece.y = 5
			original := m.currentPiece.cells()
			for range 4 {
				cells := make(map[Position]bool)
				for _, cell := range m.currentPiece.cells() {
					cells[cell] = true
				}
				if len(cells) != 4 || !m.fits(m.currentPiece) {
					t.Fatal("rotation lost cells or left the empty board")
				}
				m.rotate()
			}
			if m.currentPiece.cells() != original {
				t.Fatal("four clockwise rotations did not restore the piece")
			}
		})
	}
}

func TestCollisionAndKicks(t *testing.T) {
	t.Run("collision", func(t *testing.T) {
		m := New()
		m.currentPiece = newPiece(pieceO)
		m.currentPiece.x = -1
		before := m.currentPiece
		if !m.fits(before) || m.move(-1, 0) || m.currentPiece != before {
			t.Fatal("left wall collision changed the piece")
		}
		m.currentPiece.x = 7
		if m.move(1, 0) {
			t.Fatal("piece passed through the right wall")
		}
		m.currentPiece.y = 18
		if m.move(0, 1) {
			t.Fatal("piece passed through the floor")
		}
		m.currentPiece = newPiece(pieceO)
		m.board[2][4] = pieceJ
		if m.move(0, 1) {
			t.Fatal("piece passed through a settled block")
		}
		if !m.fits(Piece{kind: pieceI, rotation: 1, x: 3, y: -2}) {
			t.Fatal("rotation space above the board was rejected")
		}
	})

	for _, tc := range []struct {
		name  string
		piece Piece
		x, y  int
	}{
		{"left wall", Piece{kind: pieceI, rotation: 1, x: -2, y: 5}, 0, 5},
		{"right wall", Piece{kind: pieceI, rotation: 1, x: 7, y: 5}, 6, 5},
		{"floor", Piece{kind: pieceI, x: 3, y: 18}, 3, 16},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := New()
			m.currentPiece = tc.piece
			if !m.fits(m.currentPiece) || !m.rotate() || !m.fits(m.currentPiece) {
				t.Fatal("wall/floor kick did not produce a legal rotation")
			}
			if m.currentPiece.x != tc.x || m.currentPiece.y != tc.y || m.currentPiece.rotation != (tc.piece.rotation+1)%4 {
				t.Fatalf("unexpected kick: %+v", m.currentPiece)
			}
		})
	}

	t.Run("blocked rotation", func(t *testing.T) {
		m := New()
		m.currentPiece = Piece{kind: pieceT, x: 3, y: 8}
		for y := range boardHeight {
			for x := range boardWidth {
				m.board[y][x] = pieceJ
			}
		}
		for _, cell := range m.currentPiece.cells() {
			m.board[cell.y][cell.x] = empty
		}
		before := m.currentPiece
		if !m.fits(before) || m.rotate() || m.currentPiece != before {
			t.Fatal("a blocked rotation changed the piece")
		}
	})
}

func TestDropsAndLocking(t *testing.T) {
	t.Run("soft drop and gravity", func(t *testing.T) {
		m := playableModel()
		m.currentPiece = newPiece(pieceO)
		id := m.tickID
		if cmd := press(m, tea.KeyDown); cmd != nil || m.currentPiece.y != 1 || m.score != 1 || m.tickID != id {
			t.Fatal("soft drop did not move one row or changed the gravity schedule")
		}
		m.Update(tickMsg{game: m, id: id})
		if m.currentPiece.y != 2 || m.score != 1 {
			t.Fatal("gravity should move without awarding drop points")
		}
		m.currentPiece.y = 18
		next := m.nextPiece
		if cmd := press(m, tea.KeyDown); cmd == nil || occupied(m) != 4 || m.currentPiece.kind != next || m.currentPiece.y != 0 || m.score != 1 {
			t.Fatal("blocked soft drop did not lock and spawn the next piece")
		}
	})

	t.Run("hard drop uses ghost landing", func(t *testing.T) {
		m := playableModel()
		m.currentPiece = newPiece(pieceT)
		m.board[18][4] = pieceO
		landing := m.landingPiece()
		if !m.fits(landing) || m.fits(Piece{kind: landing.kind, x: landing.x, y: landing.y + 1}) {
			t.Fatal("ghost is not at the last legal row")
		}
		oldTick := tickMsg{game: m, id: m.tickID}
		next := m.nextPiece
		if cmd := press(m, tea.KeySpace); cmd == nil {
			t.Fatal("hard drop did not schedule fresh gravity for the next piece")
		}
		for _, cell := range landing.cells() {
			if m.board[cell.y][cell.x] != pieceT {
				t.Fatal("hard drop did not lock at the ghost position")
			}
		}
		if m.score != landing.y*2 || occupied(m) != 5 || m.currentPiece.kind != next || m.gameOver {
			t.Fatal("hard drop points, board, or next piece are incorrect")
		}
		before := m.currentPiece
		if _, cmd := m.Update(oldTick); cmd != nil || m.currentPiece != before {
			t.Fatal("pre-drop gravity moved the newly spawned piece")
		}
	})

	t.Run("gravity locks on contact", func(t *testing.T) {
		m := playableModel()
		m.currentPiece = Piece{kind: pieceO, x: 3, y: 18}
		next := m.nextPiece
		_, cmd := m.Update(tickMsg{game: m, id: m.tickID})
		if cmd == nil || occupied(m) != 4 || m.currentPiece.kind != next || m.currentPiece.y != 0 || m.score != 0 {
			t.Fatal("gravity did not lock a grounded piece and continue")
		}
	})
}

func TestConsecutiveRowsPreserveBlocksAbove(t *testing.T) {
	m := New()
	for y := 18; y < boardHeight; y++ {
		for x := range boardWidth {
			if x != 4 && x != 5 {
				m.board[y][x] = pieceJ
			}
		}
	}
	m.board[17][8], m.board[16][1], m.board[0][0] = pieceT, pieceZ, pieceL
	m.currentPiece = Piece{kind: pieceO, x: 3, y: 18}
	m.lockPiece()
	if m.lines != 2 || m.score != 300 || m.board[19][8] != pieceT || m.board[18][1] != pieceZ || m.board[2][0] != pieceL || occupied(m) != 3 {
		t.Fatal("consecutive clears lost or misplaced blocks above them")
	}
	if m.board[0] != ([boardWidth]pieceKind{}) || m.board[1] != ([boardWidth]pieceKind{}) {
		t.Fatal("cleared top rows were not emptied")
	}
}

func TestSeparatedRowsCompactInOrder(t *testing.T) {
	m := New()
	for x := range boardWidth {
		m.board[19][x], m.board[17][x] = pieceS, pieceZ
	}
	m.board[18][3], m.board[16][6] = pieceT, pieceL
	if got := m.clearRows(); got != 2 || m.board[19][3] != pieceT || m.board[18][6] != pieceL || occupied(m) != 2 {
		t.Fatal("separated full rows did not compact in order")
	}
}

func TestLineScoringAndLevelProgression(t *testing.T) {
	for count, points := range []int{0, 100, 300, 500, 800} {
		t.Run(fmt.Sprintf("%d lines", count), func(t *testing.T) {
			m := New()
			m.level, m.lines = 2, 10
			for y := boardHeight - count; y < boardHeight; y++ {
				for x := range boardWidth {
					if x != 4 {
						m.board[y][x] = pieceJ
					}
				}
			}
			m.currentPiece = Piece{kind: pieceI, rotation: 1, x: 2, y: 16}
			m.lockPiece()
			if m.score != points*2 || m.lines != 10+count || m.level != 2 {
				t.Fatalf("wrong line award: score=%d lines=%d level=%d", m.score, m.lines, m.level)
			}
		})
	}

	m := New()
	m.lines = 9
	initialGravity := m.gravityInterval()
	for x := range boardWidth {
		if x != 4 && x != 5 {
			m.board[19][x] = pieceJ
		}
	}
	m.currentPiece = Piece{kind: pieceO, x: 3, y: 18}
	m.lockPiece()
	if m.lines != 10 || m.level != 2 || m.score != 100 || m.gravityInterval() >= initialGravity {
		t.Fatal("crossing ten lines did not use the old level's award and increase gravity")
	}
	m.level = 1000
	if m.gravityInterval() != 80*time.Millisecond {
		t.Fatal("gravity should have a practical minimum interval")
	}
}

func TestTopOut(t *testing.T) {
	t.Run("blocked spawn", func(t *testing.T) {
		m := playableModel()
		m.currentPiece = Piece{kind: pieceO, x: 3, y: 18}
		m.nextPiece = pieceT
		m.board[0][4] = pieceJ
		_, cmd := m.Update(tickMsg{game: m, id: m.tickID})
		if !m.gameOver || cmd != nil || m.board[0][4] != pieceJ {
			t.Fatal("spawn collision did not stop gravity and preserve settled blocks")
		}
		before := m.currentPiece
		for _, key := range []rune{tea.KeyLeft, tea.KeyRight, tea.KeyUp, tea.KeyDown, tea.KeySpace, 'p'} {
			press(m, key)
		}
		if m.currentPiece != before || m.paused {
			t.Fatal("game-over controls changed the ended game")
		}
		if cmd := press(m, 'r'); cmd == nil || m.gameOver || occupied(m) != 0 {
			t.Fatal("restart did not recover from top-out")
		}
	})
	t.Run("locking above board", func(t *testing.T) {
		m := New()
		m.currentPiece = Piece{kind: pieceI, rotation: 1, x: 2, y: -2}
		m.lockPiece()
		if !m.gameOver || occupied(m) != 0 {
			t.Fatal("locking above the ceiling did not top out safely")
		}
	})
}

func TestTickPauseRestartAndOwnership(t *testing.T) {
	m := playableModel()
	m.currentPiece = newPiece(pieceT)
	first := tickMsg{game: m, id: m.tickID}
	press(m, tea.KeyLeft)
	press(m, tea.KeyRight)
	press(m, tea.KeyUp)
	if m.tickID != first.id {
		t.Fatal("movement or rotation postponed gravity")
	}
	_, cmd := m.Update(first)
	if cmd == nil || m.currentPiece.y != 1 {
		t.Fatal("valid tick did not move and rearm")
	}
	before := m.currentPiece
	if _, cmd := m.Update(first); cmd != nil || m.currentPiece != before {
		t.Fatal("duplicate tick moved the piece twice")
	}
	foreign := tickMsg{game: New(), id: m.tickID}
	if _, cmd := m.Update(foreign); cmd != nil || m.currentPiece != before {
		t.Fatal("a foreign instance's tick moved this game")
	}
	prePause := tickMsg{game: m, id: m.tickID}
	if cmd := press(m, 'p'); cmd != nil || !m.paused {
		t.Fatal("pause did not stop scheduling")
	}
	for _, key := range []rune{tea.KeyLeft, tea.KeyRight, tea.KeyUp, tea.KeyDown, tea.KeySpace} {
		press(m, key)
	}
	m.Update(prePause)
	if m.currentPiece != before || occupied(m) != 0 || m.score != 0 {
		t.Fatal("paused inputs or queued gravity changed gameplay")
	}
	if cmd := press(m, 'p'); cmd == nil || m.paused {
		t.Fatal("resume did not schedule new gravity")
	}
	m.Update(prePause)
	if m.currentPiece != before {
		t.Fatal("a pre-pause tick moved the resumed game")
	}
	m.Update(tickMsg{game: m, id: m.tickID})
	if m.currentPiece.y != before.y+1 {
		t.Fatal("new gravity did not move after resume")
	}
	old := tickMsg{game: m, id: m.tickID}
	m.score, m.lines, m.level = 450, 12, 2
	m.board[19][0] = pieceL
	press(m, 'p')
	if cmd := press(m, 'r'); cmd == nil || m.paused || m.gameOver || m.score != 0 || m.lines != 0 || m.level != 1 || occupied(m) != 0 {
		t.Fatal("restart did not reset board, statistics, and pause")
	}
	before = m.currentPiece
	m.Update(old)
	if m.currentPiece != before || m.width != 80 || m.height != 24 {
		t.Fatal("restart lost terminal size or accepted an old tick")
	}
}

func TestSmallTerminalSuspendsGameplay(t *testing.T) {
	m := New()
	if m.Init() != nil {
		t.Fatal("game should wait for a usable initial terminal size")
	}
	if _, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24}); cmd == nil {
		t.Fatal("first usable size did not start gravity")
	}
	id := m.tickID
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	if m.tickID != id {
		t.Fatal("usable resize postponed gravity")
	}
	old := tickMsg{game: m, id: m.tickID}
	if _, cmd := m.Update(tea.WindowSizeMsg{Width: viewWidth - 1, Height: viewHeight}); cmd != nil {
		t.Fatal("small terminal kept scheduling gravity")
	}
	before := m.currentPiece
	for _, key := range []rune{tea.KeyLeft, tea.KeyRight, tea.KeyUp, tea.KeyDown, tea.KeySpace} {
		press(m, key)
	}
	m.Update(old)
	if _, cmd := m.Update(tickMsg{game: m, id: m.tickID}); cmd != nil || m.currentPiece != before || occupied(m) != 0 {
		t.Fatal("hidden piece moved, locked, or scheduled gravity")
	}
	if cmd := press(m, 'r'); cmd != nil {
		t.Fatal("restart armed gravity in a small terminal")
	}
	before = m.currentPiece
	if _, cmd := m.Update(tea.WindowSizeMsg{Width: viewWidth, Height: viewHeight}); cmd == nil {
		t.Fatal("restoring terminal size did not resume gravity")
	}
	m.Update(old)
	if m.currentPiece != before {
		t.Fatal("pre-resize gravity affected restored gameplay")
	}
	press(m, 'p')
	m.Update(tea.WindowSizeMsg{Width: 10, Height: 10})
	if _, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24}); cmd != nil || !m.paused {
		t.Fatal("restoring terminal size resumed an explicitly paused game")
	}
}

func TestScheduledTickCapturesInstanceAndGeneration(t *testing.T) {
	m := playableModel()
	m.level = 100
	cmd := m.tick()
	id := m.tickID
	m.reset()
	msg, ok := cmd().(tickMsg)
	if !ok || msg.game != m || msg.id != id || msg.id == m.tickID {
		t.Fatal("timer did not capture its original instance and generation")
	}
	before := m.currentPiece
	if _, cmd := m.Update(msg); cmd != nil || m.currentPiece != before {
		t.Fatal("a timer completed after restart moved the new game")
	}
}

func TestInitAfterInitialSize(t *testing.T) {
	m := playableModel() // The menu wrapper discards this initial resize command.
	resizeTick := tickMsg{game: m, id: m.tickID}
	before := m.currentPiece
	if cmd := m.Init(); cmd == nil || m.tickID == resizeTick.id {
		t.Fatal("Init did not replace the discarded initial resize timer")
	}
	m.Update(resizeTick)
	if m.currentPiece != before {
		t.Fatal("a duplicate startup timer moved the piece")
	}
	m.Update(tickMsg{game: m, id: m.tickID})
	if m.currentPiece.y != before.y+1 {
		t.Fatal("Init's replacement timer did not move the piece")
	}
}

func TestBoardGhostAndPreview(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	m := New()
	for kind := pieceI; kind <= pieceL; kind++ {
		m.currentPiece = newPiece(kind)
		m.nextPiece = kind
		board := m.renderBoard()
		if lipgloss.Height(board) != viewHeight || strings.Count(board, "██") != 4 || strings.Count(board, "░░") != 4 {
			t.Fatalf("%s board is missing blocks or landing ghost", pieceNames[kind])
		}
		for _, row := range strings.Split(board, "\n") {
			if lipgloss.Width(row) != boardWidth*cellWidth+2 {
				t.Fatalf("%s board has an uneven row width", pieceNames[kind])
			}
		}
		preview := m.renderNextPiece()
		if lipgloss.Height(preview) != 2 || lipgloss.Width(preview) > 8 || strings.Count(preview, "██") != 4 {
			t.Fatalf("%s next piece is incomplete or oversized", pieceNames[kind])
		}
	}
	m.gameOver = true
	if strings.Contains(m.renderBoard(), "░░") {
		t.Fatal("ended game still shows a landing ghost")
	}
	if theme.CurrentTheme() == nil {
		t.Fatal("direct rendering did not initialize a theme")
	}
}

func TestViewFitsTerminalAndStates(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	for _, size := range [][2]int{{80, 24}, {48, 22}, {120, 40}, {47, 22}, {48, 21}, {30, 10}, {18, 3}, {1, 1}, {0, 0}, {-1, -1}} {
		t.Run(fmt.Sprintf("%dx%d", size[0], size[1]), func(t *testing.T) {
			m := New()
			m.width, m.height = size[0], size[1]
			view := m.View()
			if m.width <= 0 || m.height <= 0 {
				if view.Content != "" {
					t.Fatal("unknown/empty terminal should render nothing")
				}
				return
			}
			if lipgloss.Height(view.Content) > m.height || lipgloss.Width(view.Content) > m.width {
				t.Fatalf("view is %dx%d, terminal is %dx%d", lipgloss.Width(view.Content), lipgloss.Height(view.Content), m.width, m.height)
			}
			if m.width >= viewWidth && m.height >= viewHeight {
				for _, text := range []string{"T E T R I S", "SCORE 000000", "LEVEL 1 · LINES 0", "NEXT", "PLAYING", "↑ rotate", "Space hard drop", "P pause", "R restart", "Esc back"} {
					if !strings.Contains(view.Content, text) {
						t.Fatalf("view is missing %q", text)
					}
				}
				m.paused = true
				if !strings.Contains(m.View().Content, "PAUSED") || !strings.Contains(m.View().Content, "P to resume") {
					t.Fatal("pause status is missing")
				}
				m.gameOver = true
				if !strings.Contains(m.View().Content, "GAME OVER") || !strings.Contains(m.View().Content, "R to retry") {
					t.Fatal("game-over status is missing")
				}
			} else if m.width >= 18 && !strings.Contains(view.Content, "Terminal too small") {
				t.Fatal("small terminal state is missing")
			}
		})
	}
}
