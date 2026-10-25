package twentyfortyeight

import (
	"image/color"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/jakmaz/arcade/internal/theme"
)

func assertViewBounds(t *testing.T, m *Model) {
	t.Helper()
	content := m.View().Content
	if content == "" {
		return
	}
	if height := lipgloss.Height(content); height > m.height {
		t.Fatalf("view uses %d rows in a %d-row terminal", height, m.height)
	}
	for _, line := range strings.Split(content, "\n") {
		if width := lipgloss.Width(line); width > m.width {
			t.Fatalf("view uses %d columns in a %d-column terminal", width, m.width)
		}
	}
}

func TestBoardGeometryAndLargeValues(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	m := newTestModel()
	m.board = board{{2, 4, 8, 16}, {32, 64, 128, 256}, {512, 1024, 2048, 4096}, {8192, 16384, 65536, 131072}}
	content := m.renderBoard()
	rows := strings.Split(content, "\n")
	if len(rows) != boardHeight {
		t.Fatalf("board has %d rows, want %d", len(rows), boardHeight)
	}
	for _, row := range rows {
		if width := lipgloss.Width(row); width != 29 {
			t.Fatalf("board row width = %d, want 29", width)
		}
	}
	plain := ansi.Strip(content)
	for _, value := range []string{"65536", "131072", "2048"} {
		if !strings.Contains(plain, value) {
			t.Fatalf("board omitted or wrapped the exact value %s", value)
		}
	}
	if !strings.Contains(plain, "│  2   │  4   │  8   │  16  │") {
		t.Fatal("tile labels are not centered within the six-column cells")
	}
	m.board = board{}
	if !strings.Contains(ansi.Strip(m.renderBoard()), "├──────┼──────┼──────┼──────┤") {
		t.Fatal("empty board lost its visible cell structure")
	}
}

func TestViewFitsTerminalAndReportsState(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	m := newTestModel()
	m.reset()
	for _, state := range []struct {
		won, continued, lost bool
		text                 string
	}{
		{false, false, false, "Merge tiles. Reach 2048."},
		{true, false, false, "2048! Enter/Space to continue"},
		{true, true, false, "2048 reached. Keep going."},
		{false, false, true, "NO MOVES · R TO RETRY"},
		{true, false, true, "2048! Enter/Space to continue"},
		{true, true, true, "NO MOVES · R TO RETRY"},
	} {
		m.won, m.continued, m.gameOver = state.won, state.continued, state.lost
		for _, size := range [][2]int{{80, 24}, {29, viewHeight}, {40, 24}, {120, 40}} {
			m.width, m.height = size[0], size[1]
			assertViewBounds(t, m)
			content := ansi.Strip(m.View().Content)
			for _, text := range []string{state.text, "SCORE 0", "← ↑ ↓ → move · R restart", "Esc back"} {
				if !strings.Contains(content, text) {
					t.Fatalf("view is missing %q", text)
				}
			}
		}
	}
	m.width, m.height = 28, 24
	if !strings.Contains(ansi.Strip(m.View().Content), "Terminal too small") {
		t.Fatal("narrow terminal did not show a clear too-small state")
	}
	m.width, m.height = 80, viewHeight-1
	if !strings.Contains(ansi.Strip(m.View().Content), "Need 29 × 22 cells") {
		t.Fatal("short terminal did not show the required dimensions")
	}
}

func TestViewBoundsAtEverySmallSize(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	m := newTestModel()
	for width := 0; width <= 31; width++ {
		for height := 0; height <= viewHeight; height++ {
			m.width, m.height = width, height
			assertViewBounds(t, m)
		}
	}
	m.width, m.height = 80, 24
	m.board[0][0] = 1048576
	m.score = int(^uint(0) >> 1)
	assertViewBounds(t, m)
	if !strings.Contains(ansi.Strip(m.View().Content), "1048576") {
		t.Fatal("tile width did not expand for longer values")
	}
}

func TestPaletteContrastAndSystemBackground(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	m := newTestModel()
	themes, err := theme.LoadEmbeddedThemes()
	if err != nil {
		t.Fatal(err)
	}
	for _, current := range themes {
		for value := 2; value <= 131072; value *= 2 {
			style := m.tileStyle(current, value)
			if ratio := contrast(style.GetForeground(), style.GetBackground()); ratio < 4.5 {
				t.Errorf("theme %s tile %d has contrast %.2f, want at least 4.5", current.Name(), value, ratio)
			}
		}
	}

	current := theme.SystemTheme()
	for value := 2; value <= 131072; value *= 2 {
		style := m.tileStyle(current, value)
		if !style.GetReverse() {
			t.Fatal("dark system tile should use the terminal background for its label")
		}
		if _, ok := style.GetForeground().(ansi.BasicColor); !ok {
			t.Fatal("system palette was converted from ANSI to fixed RGB")
		}
		if _, ok := style.GetBackground().(lipgloss.NoColor); !ok {
			t.Fatal("system tile forced a fixed terminal background")
		}
	}
	m.Update(tea.BackgroundColorMsg{Color: color.RGBA{R: 255, G: 255, B: 255, A: 255}})
	style := m.tileStyle(current, 2)
	if m.darkBackground || style.GetReverse() {
		t.Fatal("light system background did not switch to the terminal foreground for labels")
	}
	if _, ok := style.GetForeground().(lipgloss.NoColor); !ok {
		t.Fatal("light system tile replaced the terminal's default foreground")
	}
	if m.View().BackgroundColor != nil {
		t.Fatal("game view forced a terminal background")
	}
}
