package themes

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/jakmaz/arcade/internal/theme"
	"gopkg.in/yaml.v3"
)

// TestAllThemesValid ensures all theme files in the themes directory load correctly
func TestAllThemesValid(t *testing.T) {
	themesDir := "."

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		t.Fatalf("Failed to read themes directory: %v", err)
	}

	themeCount := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		themeCount++
		t.Run(entry.Name(), func(t *testing.T) {
			themePath := filepath.Join(themesDir, entry.Name())

			// Test that theme loads without error
			themeObj, err := theme.LoadThemeFromFile(themePath)
			if err != nil {
				t.Fatalf("Failed to load theme %s: %v", entry.Name(), err)
			}

			// Test that theme has a name
			if themeObj.Name() == "" {
				t.Errorf("Theme %s has empty name", entry.Name())
			}

			// Test that theme name matches filename (without extension)
			expectedName := strings.TrimSuffix(entry.Name(), ".yaml")
			if themeObj.Name() != expectedName {
				t.Errorf("Theme name %q doesn't match filename %q", themeObj.Name(), expectedName)
			}
		})
	}

	if themeCount == 0 {
		t.Fatal("No theme files found in themes directory")
	}

	t.Logf("Successfully validated %d theme files", themeCount)
}

// TestThemeHasRequiredFields ensures all themes have the minimum required structure
func TestThemeHasRequiredFields(t *testing.T) {
	themesDir := "."

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		t.Fatalf("Failed to read themes directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			themePath := filepath.Join(themesDir, entry.Name())

			// Read and parse YAML directly to check structure
			data, err := os.ReadFile(themePath)
			if err != nil {
				t.Fatalf("Failed to read theme file: %v", err)
			}

			var def theme.ThemeDefinition
			if err := yaml.Unmarshal(data, &def); err != nil {
				t.Fatalf("Failed to parse theme YAML: %v", err)
			}

			// Check required top-level fields
			if def.Name == "" {
				t.Error("Theme missing required 'name' field")
			}

			// Check UI colors (most critical for app functionality)
			if def.UI.Primary == "" {
				t.Error("Theme missing required 'ui.primary' color")
			}
			if def.UI.Secondary == "" {
				t.Error("Theme missing required 'ui.secondary' color")
			}
			if def.UI.Accent == "" {
				t.Error("Theme missing required 'ui.accent' color")
			}

			// Check that we have at least some palette colors if UI colors reference them
			if len(def.Palette) == 0 && usesColorReferences(def) {
				t.Error("Theme uses color references but has no palette section")
			}
		})
	}
}

// TestThemeColorsValid ensures all color values are valid hex codes or color names
func TestThemeColorsValid(t *testing.T) {
	themesDir := "."

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		t.Fatalf("Failed to read themes directory: %v", err)
	}

	// Valid hex color pattern (3, 4, 6, or 8 characters after #)
	hexColorPattern := regexp.MustCompile(`^#([0-9A-Fa-f]{3}|[0-9A-Fa-f]{4}|[0-9A-Fa-f]{6}|[0-9A-Fa-f]{8})$`)

	// Common named colors that lipgloss supports
	namedColors := map[string]bool{
		"black": true, "red": true, "green": true, "yellow": true,
		"blue": true, "magenta": true, "cyan": true, "white": true,
		"none": true, "": true, // empty string is valid (no color)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			themePath := filepath.Join(themesDir, entry.Name())

			data, err := os.ReadFile(themePath)
			if err != nil {
				t.Fatalf("Failed to read theme file: %v", err)
			}

			var def theme.ThemeDefinition
			if err := yaml.Unmarshal(data, &def); err != nil {
				t.Fatalf("Failed to parse theme YAML: %v", err)
			}

			// Validate palette colors
			for colorName, colorValue := range def.Palette {
				if !isValidColor(colorValue, hexColorPattern, namedColors) {
					t.Errorf("Invalid color in palette '%s': '%s'", colorName, colorValue)
				}
			}

			// Function to validate a color, resolving references
			validateColor := func(colorValue, fieldName string) {
				if colorValue == "" || colorValue == "none" {
					return // These are valid (no color)
				}

				// If it's a reference to palette color, resolve it
				if colorValue[0] != '#' {
					if resolvedColor, exists := def.Palette[colorValue]; exists {
						colorValue = resolvedColor
					} else {
						// It's a color reference but palette doesn't have it
						if !namedColors[colorValue] {
							t.Errorf("Color reference '%s' in field '%s' not found in palette and not a known color name", colorValue, fieldName)
							return
						}
					}
				}

				if !isValidColor(colorValue, hexColorPattern, namedColors) {
					t.Errorf("Invalid color in field '%s': '%s'", fieldName, colorValue)
				}
			}

			// Validate UI colors
			validateColor(def.UI.Primary, "ui.primary")
			validateColor(def.UI.Secondary, "ui.secondary")
			validateColor(def.UI.Accent, "ui.accent")
			validateColor(def.UI.Success, "ui.success")
			validateColor(def.UI.Warning, "ui.warning")
			validateColor(def.UI.Error, "ui.error")

			// Validate board colors
			validateColor(def.Board.Border, "board.border")
			validateColor(def.Board.Background, "board.background")
			validateColor(def.Board.CellBorder, "board.cellBorder")
			validateColor(def.Board.CellBackground, "board.cellBackground")
			validateColor(def.Board.SelectedCell, "board.selectedCell")

			// Validate game colors
			validateColor(def.Games.Chess.WhitePieces, "games.chess.whitePieces")
			validateColor(def.Games.Chess.BlackPieces, "games.chess.blackPieces")
			validateColor(def.Games.Snake.Body, "games.snake.body")
			validateColor(def.Games.Snake.Head, "games.snake.head")
			validateColor(def.Games.Snake.Food, "games.snake.food")
			validateColor(def.Games.Tetris.IPiece, "games.tetris.iPiece")
			validateColor(def.Games.Tetris.OPiece, "games.tetris.oPiece")
			validateColor(def.Games.Tetris.TPiece, "games.tetris.tPiece")
			validateColor(def.Games.Tetris.SPiece, "games.tetris.sPiece")
			validateColor(def.Games.Tetris.ZPiece, "games.tetris.zPiece")
			validateColor(def.Games.Tetris.JPiece, "games.tetris.jPiece")
			validateColor(def.Games.Tetris.LPiece, "games.tetris.lPiece")
			validateColor(def.Games.Tictactoe.Player1, "games.tictactoe.player1")
			validateColor(def.Games.Tictactoe.Player2, "games.tictactoe.player2")
		})
	}
}

// Helper function to check if a theme uses color references
func usesColorReferences(def theme.ThemeDefinition) bool {
	colors := []string{
		def.UI.Primary, def.UI.Secondary, def.UI.Accent,
		def.UI.Success, def.UI.Warning, def.UI.Error,
	}

	for _, color := range colors {
		if color != "" && color[0] != '#' && color != "none" {
			return true
		}
	}
	return false
}

// Helper function to validate if a color value is valid
func isValidColor(color string, hexPattern *regexp.Regexp, namedColors map[string]bool) bool {
	// Check if it's a valid hex color
	if hexPattern.MatchString(color) {
		return true
	}

	// Check if it's a known named color
	return namedColors[strings.ToLower(color)]
}
