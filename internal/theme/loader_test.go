package theme

import (
	"image/color"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestEmbeddedGameAndBoardColorsAreLoaded(t *testing.T) {
	themes, err := LoadEmbeddedThemes()
	if err != nil {
		t.Fatal(err)
	}
	if len(themes) == 0 {
		t.Fatal("no embedded themes loaded")
	}
	for _, theme := range themes {
		t.Run(theme.Name(), func(t *testing.T) {
			for name, value := range map[string]color.Color{
				"cellBorder": theme.CellBorder(), "selectedCell": theme.SelectedCell(),
				"I": theme.TetrisI(), "O": theme.TetrisO(), "T": theme.TetrisT(),
				"S": theme.TetrisS(), "Z": theme.TetrisZ(), "J": theme.TetrisJ(), "L": theme.TetrisL(),
			} {
				if _, missing := value.(lipgloss.NoColor); missing || value == nil {
					t.Errorf("%s was not decoded from the theme YAML", name)
				}
			}
		})
	}
}
