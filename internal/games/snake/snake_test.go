package snake

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestMovementAndFood(t *testing.T) {
	m := New()
	m.Update(tea.KeyPressMsg{Code: tea.KeyLeft}) // A turn directly into the body is ignored.
	if m.nextDirection != Right {
		t.Fatal("reverse turn was accepted")
	}

	for range 6 {
		m.advance()
	}
	if m.snake[0] != (Position{20, 8}) || len(m.snake) != 4 || m.score != 1 {
		t.Fatalf("unexpected state after eating: head=%v length=%d score=%d", m.snake[0], len(m.snake), m.score)
	}
	for _, segment := range m.snake {
		if m.food == segment {
			t.Fatal("food spawned inside the snake")
		}
	}

	m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m.advance()
	if m.snake[0] != (Position{20, 7}) {
		t.Fatalf("arrow turn did not apply on the next step: %v", m.snake[0])
	}
}

func TestCollisions(t *testing.T) {
	t.Run("wall", func(t *testing.T) {
		m := New()
		m.snake = []Position{{boardWidth - 1, 8}, {boardWidth - 2, 8}}
		m.advance()
		if !m.gameOver {
			t.Fatal("moving through the wall did not end the game")
		}
	})

	t.Run("body", func(t *testing.T) {
		m := New()
		m.snake = []Position{{3, 3}, {3, 2}, {4, 2}, {4, 3}, {5, 3}}
		m.advance()
		if !m.gameOver {
			t.Fatal("moving into the body did not end the game")
		}
	})

	t.Run("vacating tail", func(t *testing.T) {
		m := New()
		m.snake = []Position{{3, 3}, {3, 2}, {4, 2}, {4, 3}}
		m.advance()
		if m.gameOver || m.snake[0] != (Position{4, 3}) {
			t.Fatal("moving into a vacating tail should be allowed")
		}
	})
}

func TestTickLifecycle(t *testing.T) {
	m := New()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	firstTick := tickMsg{game: m, id: m.tickID}
	m.Update(firstTick)
	if m.snake[0] != (Position{15, 8}) {
		t.Fatal("a valid tick did not move the snake")
	}

	m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if !m.paused {
		t.Fatal("space did not pause the game")
	}
	m.Update(firstTick)
	if m.snake[0] != (Position{15, 8}) {
		t.Fatal("a stale tick moved a paused snake")
	}

	m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	resumeTick := tickMsg{game: m, id: m.tickID}
	m.Update(resumeTick)
	if m.snake[0] != (Position{16, 8}) {
		t.Fatal("snake did not resume")
	}

	m.Update(tea.KeyPressMsg{Code: 'r'})
	m.Update(resumeTick)
	if m.snake[0] != (Position{14, 8}) || m.score != 0 {
		t.Fatal("restart did not reset the game or reject an old tick")
	}

	otherGame := New()
	otherGame.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	otherGame.Update(tickMsg{game: m, id: otherGame.tickID})
	if otherGame.snake[0] != (Position{14, 8}) {
		t.Fatal("a tick from a previous game moved a new game")
	}
}

func TestBoardRowsHaveEqualWidth(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	m := New()
	for _, row := range strings.Split(m.renderBoard(), "\n") {
		if width := lipgloss.Width(row); width != viewWidth {
			t.Fatalf("board row is %d cells wide, want %d", width, viewWidth)
		}
	}
}

func TestViewFitsTerminal(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	m := New()
	m.width, m.height = 80, 24
	view := m.View()
	lines := strings.Split(view.Content, "\n")
	if len(lines) > m.height {
		t.Fatalf("view uses %d rows in a %d-row terminal", len(lines), m.height)
	}
	for _, line := range lines {
		if width := lipgloss.Width(line); width > m.width {
			t.Fatalf("view uses %d columns in a %d-column terminal", width, m.width)
		}
	}
	if !strings.Contains(view.Content, "SCORE 000") || !strings.Contains(view.Content, "◆") {
		t.Fatal("view is missing the score or food")
	}

	m.width = viewWidth - 1
	if !strings.Contains(m.View().Content, "Terminal too small") {
		t.Fatal("small terminal state was not shown")
	}
}

func BenchmarkView(b *testing.B) {
	b.Setenv("XDG_DATA_HOME", b.TempDir())
	m := New()
	m.width, m.height = 80, 24
	for b.Loop() {
		m.View()
	}
}
