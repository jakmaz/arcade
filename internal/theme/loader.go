package theme

import (
	"embed"
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss/v2"
	"gopkg.in/yaml.v3"
)

//go:embed themes/*.yaml
var themesFS embed.FS

// ThemeDefinition represents the YAML structure for theme files
type ThemeDefinition struct {
	Name    string
	Palette map[string]string
	UI      UIColors
	Board   BoardColors
	Games   GameColors
}

type UIColors struct {
	Primary   string
	Secondary string
	Accent    string
	Success   string
	Warning   string
	Error     string
}

type BoardColors struct {
	Border         string
	Background     string
	CellBorder     string
	CellBackground string
	SelectedCell   string
}

type GameColors struct {
	Chess     ChessColors
	Snake     SnakeColors
	Tetris    TetrisColors
	Tictactoe TicTacToeColors
}

type ChessColors struct {
	WhitePieces string
	BlackPieces string
}

type SnakeColors struct {
	Body string
	Head string
	Food string
}

type TetrisColors struct {
	IPiece string
	OPiece string
	TPiece string
	SPiece string
	ZPiece string
	JPiece string
	LPiece string
}

type TicTacToeColors struct {
	Player1 string
	Player2 string
}

// LoadThemeFromFile loads a theme from a YAML file
func LoadThemeFromFile(path string) (Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file %s: %w", path, err)
	}

	var def ThemeDefinition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("failed to parse theme file %s: %w", path, err)
	}

	return createThemeFromDefinition(&def)
}

// LoadThemesFromDirectories loads themes from user directories in the correct override order.
// The hierarchy is (from lowest to highest priority):
// 1. Built-in themes (embedded)
// 2. USER_CONFIG/opencode/themes/*.yaml
func LoadThemesFromDirectory(dir string) ([]Theme, error) {
	var themes []Theme

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return themes, nil // Directory doesn't exist, return empty slice
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read themes directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		themePath := filepath.Join(dir, entry.Name())
		theme, err := LoadThemeFromFile(themePath)
		if err != nil {
			// Log error but continue loading other themes
			fmt.Fprintf(os.Stderr, "Warning: failed to load theme %s: %v\n", themePath, err)
			continue
		}

		themes = append(themes, theme)
	}

	return themes, nil
}

// createThemeFromDefinition creates a BaseTheme from a ThemeDefinition
func createThemeFromDefinition(def *ThemeDefinition) (*BaseTheme, error) {
	theme := &BaseTheme{name: def.Name}

	// Helper function to resolve color references and create color.Color
	resolveColor := func(colorStr string) color.Color {
		// Resolve color reference if it doesn't start with #
		if colorStr != "" && colorStr[0] != '#' {
			if refColor, exists := def.Palette[colorStr]; exists {
				colorStr = refColor
			}
		}

		if colorStr == "" || colorStr == "none" {
			return lipgloss.Color("")
		}

		return lipgloss.Color(colorStr)
	}

	// UI Colors - no fallbacks, YAML themes must be complete
	theme.primary = resolveColor(def.UI.Primary)
	theme.secondary = resolveColor(def.UI.Secondary)
	theme.accent = resolveColor(def.UI.Accent)
	theme.success = resolveColor(def.UI.Success)
	theme.warning = resolveColor(def.UI.Warning)
	theme.error = resolveColor(def.UI.Error)

	// Board Colors
	theme.boardBorder = resolveColor(def.Board.Border)
	theme.boardBackground = resolveColor(def.Board.Background)
	theme.cellBorder = resolveColor(def.Board.CellBorder)
	theme.cellBackground = resolveColor(def.Board.CellBackground)
	theme.selectedCell = resolveColor(def.Board.SelectedCell)

	// Game Colors - TicTacToe
	theme.player1 = resolveColor(def.Games.Tictactoe.Player1)
	theme.player2 = resolveColor(def.Games.Tictactoe.Player2)

	// Game Colors - Snake
	theme.snakeBody = resolveColor(def.Games.Snake.Body)
	theme.snakeHead = resolveColor(def.Games.Snake.Head)
	theme.food = resolveColor(def.Games.Snake.Food)

	// Game Colors - Chess
	theme.whitePiece = resolveColor(def.Games.Chess.WhitePieces)
	theme.blackPiece = resolveColor(def.Games.Chess.BlackPieces)

	// Game Colors - Tetris
	theme.tetrisI = resolveColor(def.Games.Tetris.IPiece)
	theme.tetrisO = resolveColor(def.Games.Tetris.OPiece)
	theme.tetrisT = resolveColor(def.Games.Tetris.TPiece)
	theme.tetrisS = resolveColor(def.Games.Tetris.SPiece)
	theme.tetrisZ = resolveColor(def.Games.Tetris.ZPiece)
	theme.tetrisJ = resolveColor(def.Games.Tetris.JPiece)
	theme.tetrisL = resolveColor(def.Games.Tetris.LPiece)

	// Terminal Background
	if bgColor, exists := def.Palette["bg"]; exists && bgColor != "" {
		theme.terminalBackground = resolveColor(bgColor)
	} else {
		theme.terminalBackground = lipgloss.Color("")
	}

	return theme, nil
}

// LoadEmbeddedThemes loads all embedded themes from the themes directory
func LoadEmbeddedThemes() ([]Theme, error) {
	var themes []Theme

	entries, err := themesFS.ReadDir("themes")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded themes directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		themePath := filepath.Join("themes", entry.Name())
		data, err := themesFS.ReadFile(themePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read embedded theme %s: %v\n", themePath, err)
			continue
		}

		var def ThemeDefinition
		if err := yaml.Unmarshal(data, &def); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to parse embedded theme %s: %v\n", themePath, err)
			continue
		}

		theme, err := createThemeFromDefinition(&def)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to create theme from %s: %v\n", themePath, err)
			continue
		}

		themes = append(themes, theme)
	}

	return themes, nil
}
