package core

import (
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jakmaz/arcade/internal/theme"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

func TestFinalGameLineup(t *testing.T) {
	want := []string{"2048", "snake", "tetris", "tictactoe"}
	var got []string
	for _, game := range AvailableGames() {
		got = append(got, game.ID)
		if CreateGame(game.ID) == nil {
			t.Fatalf("%s has no playable model", game.ID)
		}
	}
	if !slices.Equal(got, want) {
		t.Fatalf("games = %v, want %v", got, want)
	}
	if CreateGame("chess") != nil || CreateGame("missing") != nil {
		t.Fatal("unregistered game returned a model")
	}
}

func TestGamesFitStandardTerminalInEveryTheme(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	if err := theme.Initialize(); err != nil {
		t.Fatal(err)
	}
	for _, name := range theme.ListThemes() {
		t.Run(name, func(t *testing.T) {
			if err := theme.SetCurrentTheme(name); err != nil {
				t.Fatal(err)
			}
			styles.RefreshStyles()
			for _, game := range AvailableGames() {
				t.Run(game.ID, func(t *testing.T) {
					m := CreateGame(game.ID)
					m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
					view := m.View().Content
					if strings.Contains(strings.ToLower(view), "terminal too small") {
						t.Fatal("game cannot be played at 80×24")
					}
					if lipgloss.Width(view) > 80 || lipgloss.Height(view) > 24 {
						t.Fatalf("view overflows 80×24: %d×%d", lipgloss.Width(view), lipgloss.Height(view))
					}
				})
			}
		})
	}
}
