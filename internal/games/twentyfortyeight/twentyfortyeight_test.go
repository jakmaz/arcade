package twentyfortyeight

import (
	"math/rand/v2"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestModel() *Model {
	random := rand.New(rand.NewPCG(1, 2))
	return &Model{random: random.IntN, width: 80, height: 24, darkBackground: true}
}

func countTiles(b board) int {
	count := 0
	for _, row := range b {
		for _, value := range row {
			if value != 0 {
				count++
			}
		}
	}
	return count
}

func TestMergeLine(t *testing.T) {
	tests := []struct {
		name   string
		before [4]int
		after  [4]int
		score  int
	}{
		{"two pairs", [4]int{2, 2, 2, 2}, [4]int{4, 4}, 8},
		{"no double merge", [4]int{2, 2, 4}, [4]int{4, 4}, 4},
		{"leading pair", [4]int{4, 4, 4}, [4]int{8, 4}, 8},
		{"different pairs", [4]int{4, 4, 8, 8}, [4]int{8, 16}, 24},
		{"gaps", [4]int{0, 2, 0, 2}, [4]int{4}, 4},
		{"slide only", [4]int{0, 2, 4, 8}, [4]int{2, 4, 8}, 0},
		{"already packed", [4]int{2, 4, 8, 16}, [4]int{2, 4, 8, 16}, 0},
		{"empty", [4]int{}, [4]int{}, 0},
		{"large tiles", [4]int{65536, 65536, 131072, 131072}, [4]int{131072, 262144}, 393216},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			after, score := mergeLine(test.before)
			if after != test.after || score != test.score {
				t.Fatalf("mergeLine(%v) = %v, %d; want %v, %d", test.before, after, score, test.after, test.score)
			}
		})
	}
}

func TestSlideOrientations(t *testing.T) {
	before := board{
		{2, 2, 4, 0},
		{0, 2, 0, 2},
		{4, 0, 4, 4},
		{2, 4, 8, 16},
	}
	tests := []struct {
		name      string
		direction direction
		after     board
		score     int
	}{
		{"left", left, board{{4, 4, 0, 0}, {4, 0, 0, 0}, {8, 4, 0, 0}, {2, 4, 8, 16}}, 16},
		{"right", right, board{{0, 0, 4, 4}, {0, 0, 0, 4}, {0, 0, 4, 8}, {2, 4, 8, 16}}, 16},
		{"up", up, board{{2, 4, 8, 2}, {4, 4, 8, 4}, {2, 0, 0, 16}, {0, 0, 0, 0}}, 12},
		{"down", down, board{{0, 0, 0, 0}, {2, 0, 0, 2}, {4, 4, 8, 4}, {2, 4, 8, 16}}, 12},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			after, score, changed := slide(before, test.direction)
			if after != test.after || score != test.score || !changed {
				t.Fatalf("slide = %v, %d, %t; want %v, %d, true", after, score, changed, test.after, test.score)
			}
		})
	}
	if before[0] != [4]int{2, 2, 4, 0} {
		t.Fatal("slide modified its input board")
	}
}

func TestMergeConservesTileSum(t *testing.T) {
	values := [...]int{0, 2, 4, 8}
	for a := range values {
		for b := range values {
			for c := range values {
				for d := range values {
					line := [4]int{values[a], values[b], values[c], values[d]}
					result, _ := mergeLine(line)
					before, after := 0, 0
					seenEmpty := false
					for i, value := range result {
						before += line[i]
						after += value
						if value == 0 {
							seenEmpty = true
						} else if seenEmpty {
							t.Fatalf("result has an interior gap: %v -> %v", line, result)
						}
					}
					if before != after {
						t.Fatalf("tile sum changed: %v -> %v", line, result)
					}
				}
			}
		}
	}
}

func TestSpawnPolicy(t *testing.T) {
	fours, twos := 0, 0
	for roll := range 10 {
		m := newTestModel()
		m.board = board{{2, 4, 8, 16}}
		before := m.board
		calls := 0
		m.random = func(n int) int {
			calls++
			if calls == 1 {
				if n != 12 {
					t.Fatalf("position draw bound = %d, want 12 empty cells", n)
				}
				return n - 1
			}
			if calls != 2 || n != 10 {
				t.Fatalf("unexpected spawn draw %d with bound %d", calls, n)
			}
			return roll
		}
		if !m.spawnTile() || calls != 2 || countTiles(m.board) != countTiles(before)+1 {
			t.Fatal("spawn did not add exactly one tile with two random draws")
		}
		if m.board[0] != before[0] {
			t.Fatal("spawn overwrote an occupied tile")
		}
		switch value := m.board[3][3]; value {
		case 2:
			twos++
		case 4:
			fours++
		default:
			t.Fatalf("spawn selected the wrong cell or value: %v", m.board)
		}
	}
	if fours != 1 || twos != 9 {
		t.Fatalf("ten equiprobable rolls produced %d twos and %d fours", twos, fours)
	}

	m := newTestModel()
	m.board = deadBoard()
	m.random = func(int) int { t.Fatal("full board consumed randomness"); return 0 }
	if m.spawnTile() || m.board != deadBoard() {
		t.Fatal("spawn changed a full board")
	}
}

func TestEffectiveMoveAndNoOp(t *testing.T) {
	m := newTestModel()
	m.board = board{{2}}
	m.score = 17
	m.random = func(int) int { t.Fatal("no-op consumed randomness"); return 0 }
	if m.move(left) || m.move(up) || m.board != (board{{2}}) || m.score != 17 || m.gameOver {
		t.Fatal("an ineffective move changed the board, score, or playable state")
	}

	m.board = board{{2, 2, 4}}
	calls := 0
	m.random = func(n int) int {
		calls++
		if calls == 1 {
			if n != 14 {
				t.Fatalf("expected 14 empty positions after merging, got %d", n)
			}
			return 0
		}
		if calls != 2 || n != 10 {
			t.Fatal("effective move did not spawn exactly once")
		}
		return 1
	}
	if !m.move(left) || m.board != (board{{4, 4, 2}}) || m.score != 21 || calls != 2 {
		t.Fatalf("unexpected merge/spawn state: board=%v score=%d draws=%d", m.board, m.score, calls)
	}
}

func deadBoard() board {
	return board{{2, 4, 2, 4}, {4, 2, 4, 2}, {2, 4, 2, 4}, {4, 2, 4, 2}}
}

func TestAvailableMovesAndLoss(t *testing.T) {
	if hasMoves(deadBoard()) {
		t.Fatal("alternating full board should have no moves")
	}
	for _, change := range []struct{ row, column, value int }{{3, 3, 0}, {3, 3, 4}, {1, 0, 2}} {
		b := deadBoard()
		b[change.row][change.column] = change.value
		if !hasMoves(b) {
			t.Fatalf("missed an empty cell or neighboring pair: %v", b)
		}
	}

	m := newTestModel()
	m.board = deadBoard()
	m.random = func(int) int { t.Fatal("dead board consumed randomness"); return 0 }
	if m.move(left) || !m.gameOver {
		t.Fatal("ineffective move on a dead board did not detect loss")
	}

	m = newTestModel()
	m.board = deadBoard()
	m.board[3] = [4]int{16, 32, 64, 0}
	m.random = func(int) int { return 0 }
	if !m.move(right) || !m.gameOver || m.board[3] != [4]int{4, 16, 32, 64} {
		t.Fatalf("final spawn did not detect loss: %v", m.board)
	}
	before := m.board
	if m.move(left) || m.board != before {
		t.Fatal("game-over board accepted further movement")
	}
}

func TestWinAndExplicitContinue(t *testing.T) {
	for _, key := range []rune{tea.KeyEnter, tea.KeySpace} {
		t.Run(tea.KeyPressMsg{Code: key}.String(), func(t *testing.T) {
			m := newTestModel()
			m.board = board{{1024, 1024, 2}}
			if !m.move(left) || !m.won || m.continued || m.score != 2048 {
				t.Fatalf("2048 merge did not pause for a win: %+v", m)
			}
			before := m.board
			m.random = func(int) int { t.Fatal("win pause consumed randomness"); return 0 }
			for _, d := range []direction{left, right, up, down} {
				if m.move(d) || m.board != before {
					t.Fatal("win pause accepted movement")
				}
			}
			m.Update(tea.KeyPressMsg{Code: key})
			if !m.continued || m.board != before {
				t.Fatal("continue did not acknowledge the win without moving")
			}
			m.board = board{{2048, 2048}}
			m.random = func(n int) int { return n - 1 }
			if !m.move(left) || m.board[0][0] != 4096 || m.score != 6144 || !m.continued {
				t.Fatal("continued play did not merge and score beyond 2048")
			}
			if !m.move(right) {
				t.Fatal("continued play incorrectly paused for another win")
			}
		})
	}
}

func TestSimultaneousWinAndLoss(t *testing.T) {
	m := newTestModel()
	m.board = board{{1024, 1024, 4, 8}, {16, 32, 64, 128}, {2, 4, 8, 16}, {32, 64, 128, 256}}
	m.random = func(n int) int { return n - 1 }
	if !m.move(left) || !m.won || !m.gameOver || m.continued {
		t.Fatal("final 2048 merge should register both win and exhausted moves")
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !m.continued || !m.gameOver || m.move(right) {
		t.Fatal("acknowledging a win should preserve an exhausted board")
	}
}

func TestStartAndRestart(t *testing.T) {
	m := New()
	if countTiles(m.board) != 2 || m.score != 0 || m.won || m.continued || m.gameOver {
		t.Fatal("new game should start with two tiles, zero score, and no result")
	}
	for _, row := range m.board {
		for _, value := range row {
			if value != 0 && value != 2 && value != 4 {
				t.Fatalf("unexpected starting tile %d", value)
			}
		}
	}
	m.width, m.height, m.darkBackground = 60, 30, false
	m.board, m.score = deadBoard(), 20000
	m.won, m.continued, m.gameOver = true, true, true
	calls := 0
	m.random = func(int) int { calls++; return 0 }
	m.Update(tea.KeyPressMsg{Code: 'r'})
	if m.board != (board{{4, 4}}) || countTiles(m.board) != 2 || m.score != 0 || calls != 4 || m.won || m.continued || m.gameOver {
		t.Fatal("restart did not reset score, board, and result with two new tiles")
	}
	if m.width != 60 || m.height != 30 || m.darkBackground {
		t.Fatal("restart discarded terminal dimensions or background information")
	}
}

func TestArrowControlsAndSmallTerminal(t *testing.T) {
	for _, test := range []struct {
		key      rune
		row, col int
	}{
		{tea.KeyLeft, 1, 0}, {tea.KeyRight, 1, 3},
		{tea.KeyUp, 0, 1}, {tea.KeyDown, 3, 1},
	} {
		m := newTestModel()
		m.board[1][1] = 8
		m.Update(tea.KeyPressMsg{Code: test.key})
		if m.board[test.row][test.col] != 8 || countTiles(m.board) != 2 {
			t.Fatalf("arrow key %v moved incorrectly: %v", test.key, m.board)
		}
	}

	m := newTestModel()
	m.board = board{{2}}
	m.Update(tea.WindowSizeMsg{Width: 10, Height: 5})
	m.random = func(int) int { t.Fatal("hidden board consumed randomness"); return 0 }
	m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if m.board != (board{{2}}) {
		t.Fatal("arrow keys moved a board that could not be displayed")
	}
	for _, key := range []rune{'w', 'a', 's', 'd', tea.KeyEnter, tea.KeySpace} {
		m.Update(tea.KeyPressMsg{Code: key})
	}
	if m.board != (board{{2}}) || m.continued {
		t.Fatal("unrelated controls altered an active game")
	}
}
