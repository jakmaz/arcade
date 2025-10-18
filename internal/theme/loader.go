package theme

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

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

	// Helper function to resolve color references and create lipgloss.TerminalColor
	resolveColor := func(colorStr string) lipgloss.TerminalColor {
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

	// UI Colors with fallbacks
	if def.UI.Primary != "" {
		theme.primary = resolveColor(def.UI.Primary)
	} else {
		theme.primary = lipgloss.Color("#ffffff")
	}

	if def.UI.Secondary != "" {
		theme.secondary = resolveColor(def.UI.Secondary)
	} else {
		theme.secondary = lipgloss.Color("#888888")
	}

	if def.UI.Accent != "" {
		theme.accent = resolveColor(def.UI.Accent)
	} else {
		theme.accent = lipgloss.Color("#0066cc")
	}

	if def.UI.Success != "" {
		theme.success = resolveColor(def.UI.Success)
	} else {
		theme.success = lipgloss.Color("#22c55e")
	}

	if def.UI.Warning != "" {
		theme.warning = resolveColor(def.UI.Warning)
	} else {
		theme.warning = lipgloss.Color("#f59e0b")
	}

	if def.UI.Error != "" {
		theme.error = resolveColor(def.UI.Error)
	} else {
		theme.error = lipgloss.Color("#ef4444")
	}

	// Board Colors with fallbacks
	if def.Board.Border != "" {
		theme.boardBorder = resolveColor(def.Board.Border)
	} else {
		theme.boardBorder = theme.secondary
	}

	if def.Board.Background != "" {
		theme.boardBackground = resolveColor(def.Board.Background)
	} else {
		theme.boardBackground = lipgloss.Color("")
	}

	if def.Board.CellBorder != "" {
		theme.cellBorder = resolveColor(def.Board.CellBorder)
	} else {
		theme.cellBorder = theme.secondary
	}

	if def.Board.CellBackground != "" {
		theme.cellBackground = resolveColor(def.Board.CellBackground)
	} else {
		theme.cellBackground = lipgloss.Color("")
	}

	if def.Board.SelectedCell != "" {
		theme.selectedCell = resolveColor(def.Board.SelectedCell)
	} else {
		theme.selectedCell = theme.accent
	}

	// Game Colors - TicTacToe
	if def.Games.Tictactoe.Player1 != "" {
		theme.player1 = resolveColor(def.Games.Tictactoe.Player1)
	} else {
		theme.player1 = lipgloss.Color("#22c55e")
	}

	if def.Games.Tictactoe.Player2 != "" {
		theme.player2 = resolveColor(def.Games.Tictactoe.Player2)
	} else {
		theme.player2 = lipgloss.Color("#ef4444")
	}

	// Game Colors - Snake
	if def.Games.Snake.Body != "" {
		theme.snakeBody = resolveColor(def.Games.Snake.Body)
	} else {
		theme.snakeBody = theme.success
	}

	if def.Games.Snake.Head != "" {
		theme.snakeHead = resolveColor(def.Games.Snake.Head)
	} else {
		theme.snakeHead = theme.accent
	}

	if def.Games.Snake.Food != "" {
		theme.food = resolveColor(def.Games.Snake.Food)
	} else {
		theme.food = lipgloss.Color("#ef4444")
	}

	// Game Colors - Chess
	if def.Games.Chess.WhitePieces != "" {
		theme.whitePiece = resolveColor(def.Games.Chess.WhitePieces)
	} else {
		theme.whitePiece = lipgloss.Color("#ffffff")
	}

	if def.Games.Chess.BlackPieces != "" {
		theme.blackPiece = resolveColor(def.Games.Chess.BlackPieces)
	} else {
		theme.blackPiece = lipgloss.Color("#444444")
	}

	// Game Colors - Tetris
	if def.Games.Tetris.IPiece != "" {
		theme.tetrisI = resolveColor(def.Games.Tetris.IPiece)
	} else {
		theme.tetrisI = lipgloss.Color("#00f5ff")
	}

	if def.Games.Tetris.OPiece != "" {
		theme.tetrisO = resolveColor(def.Games.Tetris.OPiece)
	} else {
		theme.tetrisO = lipgloss.Color("#ffff00")
	}

	if def.Games.Tetris.TPiece != "" {
		theme.tetrisT = resolveColor(def.Games.Tetris.TPiece)
	} else {
		theme.tetrisT = lipgloss.Color("#800080")
	}

	if def.Games.Tetris.SPiece != "" {
		theme.tetrisS = resolveColor(def.Games.Tetris.SPiece)
	} else {
		theme.tetrisS = lipgloss.Color("#00ff00")
	}

	if def.Games.Tetris.ZPiece != "" {
		theme.tetrisZ = resolveColor(def.Games.Tetris.ZPiece)
	} else {
		theme.tetrisZ = lipgloss.Color("#ff0000")
	}

	if def.Games.Tetris.JPiece != "" {
		theme.tetrisJ = resolveColor(def.Games.Tetris.JPiece)
	} else {
		theme.tetrisJ = lipgloss.Color("#0000ff")
	}

	if def.Games.Tetris.LPiece != "" {
		theme.tetrisL = resolveColor(def.Games.Tetris.LPiece)
	} else {
		theme.tetrisL = lipgloss.Color("#ffa500")
	}

	// Terminal Background
	if bgColor, exists := def.Palette["bg"]; exists {
		theme.terminalBackground = resolveColor(bgColor)
		theme.useTerminalBackground = true
	} else {
		theme.terminalBackground = lipgloss.Color("")
		theme.useTerminalBackground = false
	}

	return theme, nil
}
