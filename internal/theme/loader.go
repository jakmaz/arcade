package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

// ThemeDefinition represents the JSON structure for theme files
type ThemeDefinition struct {
	Name  string                `json:"name"`
	Defs  map[string]string     `json:"defs"`
	Theme map[string]ColorValue `json:"theme"`
}

// ColorValue can be either a string or an object with dark/light variants
type ColorValue struct {
	Dark  string `json:"dark,omitempty"`
	Light string `json:"light,omitempty"`
	Value string `json:"-"` // For simple string values
}

// UnmarshalJSON handles both string and object color values
func (cv *ColorValue) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		cv.Value = str
		return nil
	}

	// Try to unmarshal as object
	type colorValueAlias ColorValue
	var obj colorValueAlias
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	cv.Dark = obj.Dark
	cv.Light = obj.Light
	return nil
}

// LoadThemeFromFile loads a theme from a JSON file
func LoadThemeFromFile(path string) (Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file %s: %w", path, err)
	}

	var def ThemeDefinition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("failed to parse theme file %s: %w", path, err)
	}

	return createThemeFromDefinition(&def)
}

// LoadThemesFromDirectory loads all theme files from a directory
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
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
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
	resolveColor := func(cv ColorValue) lipgloss.TerminalColor {
		var colorStr string

		if cv.Value != "" {
			colorStr = cv.Value
		} else if cv.Dark != "" {
			// For now, use dark variant (we can enhance this later for light terminal detection)
			colorStr = cv.Dark
		} else if cv.Light != "" {
			colorStr = cv.Light
		}

		// Resolve color reference if it starts with a reference
		if colorStr != "" && colorStr[0] != '#' {
			if refColor, exists := def.Defs[colorStr]; exists {
				colorStr = refColor
			}
		}

		if colorStr == "" || colorStr == "none" {
			return lipgloss.Color("")
		}

		return lipgloss.Color(colorStr)
	}

	// Map theme colors with fallbacks
	if cv, exists := def.Theme["primary"]; exists {
		theme.primary = resolveColor(cv)
	} else {
		theme.primary = lipgloss.Color("#ffffff")
	}

	if cv, exists := def.Theme["secondary"]; exists {
		theme.secondary = resolveColor(cv)
	} else {
		theme.secondary = lipgloss.Color("#888888")
	}

	if cv, exists := def.Theme["accent"]; exists {
		theme.accent = resolveColor(cv)
	} else {
		theme.accent = lipgloss.Color("#0066cc")
	}

	if cv, exists := def.Theme["success"]; exists {
		theme.success = resolveColor(cv)
	} else {
		theme.success = lipgloss.Color("#22c55e")
	}

	if cv, exists := def.Theme["warning"]; exists {
		theme.warning = resolveColor(cv)
	} else {
		theme.warning = lipgloss.Color("#f59e0b")
	}

	if cv, exists := def.Theme["error"]; exists {
		theme.error = resolveColor(cv)
	} else {
		theme.error = lipgloss.Color("#ef4444")
	}

	// Game board colors
	if cv, exists := def.Theme["board_border"]; exists {
		theme.boardBorder = resolveColor(cv)
	} else {
		theme.boardBorder = theme.secondary
	}

	if cv, exists := def.Theme["board_background"]; exists {
		theme.boardBackground = resolveColor(cv)
	} else {
		theme.boardBackground = lipgloss.Color("")
	}

	if cv, exists := def.Theme["cell_border"]; exists {
		theme.cellBorder = resolveColor(cv)
	} else {
		theme.cellBorder = theme.secondary
	}

	if cv, exists := def.Theme["cell_background"]; exists {
		theme.cellBackground = resolveColor(cv)
	} else {
		theme.cellBackground = lipgloss.Color("")
	}

	if cv, exists := def.Theme["selected_cell"]; exists {
		theme.selectedCell = resolveColor(cv)
	} else {
		theme.selectedCell = theme.accent
	}

	// Game piece colors
	if cv, exists := def.Theme["player1"]; exists {
		theme.player1 = resolveColor(cv)
	} else {
		theme.player1 = lipgloss.Color("#22c55e")
	}

	if cv, exists := def.Theme["player2"]; exists {
		theme.player2 = resolveColor(cv)
	} else {
		theme.player2 = lipgloss.Color("#ef4444")
	}

	if cv, exists := def.Theme["snake_body"]; exists {
		theme.snakeBody = resolveColor(cv)
	} else {
		theme.snakeBody = theme.success
	}

	if cv, exists := def.Theme["snake_head"]; exists {
		theme.snakeHead = resolveColor(cv)
	} else {
		theme.snakeHead = theme.accent
	}

	if cv, exists := def.Theme["food"]; exists {
		theme.food = resolveColor(cv)
	} else {
		theme.food = lipgloss.Color("#ef4444")
	}

	// Chess piece colors
	if cv, exists := def.Theme["white_piece"]; exists {
		theme.whitePiece = resolveColor(cv)
	} else {
		theme.whitePiece = lipgloss.Color("#ffffff")
	}

	if cv, exists := def.Theme["black_piece"]; exists {
		theme.blackPiece = resolveColor(cv)
	} else {
		theme.blackPiece = lipgloss.Color("#444444")
	}

	// Tetris colors
	if cv, exists := def.Theme["tetris_i"]; exists {
		theme.tetrisI = resolveColor(cv)
	} else {
		theme.tetrisI = lipgloss.Color("#00f5ff")
	}

	if cv, exists := def.Theme["tetris_o"]; exists {
		theme.tetrisO = resolveColor(cv)
	} else {
		theme.tetrisO = lipgloss.Color("#ffff00")
	}

	if cv, exists := def.Theme["tetris_t"]; exists {
		theme.tetrisT = resolveColor(cv)
	} else {
		theme.tetrisT = lipgloss.Color("#800080")
	}

	if cv, exists := def.Theme["tetris_s"]; exists {
		theme.tetrisS = resolveColor(cv)
	} else {
		theme.tetrisS = lipgloss.Color("#00ff00")
	}

	if cv, exists := def.Theme["tetris_z"]; exists {
		theme.tetrisZ = resolveColor(cv)
	} else {
		theme.tetrisZ = lipgloss.Color("#ff0000")
	}

	if cv, exists := def.Theme["tetris_j"]; exists {
		theme.tetrisJ = resolveColor(cv)
	} else {
		theme.tetrisJ = lipgloss.Color("#0000ff")
	}

	if cv, exists := def.Theme["tetris_l"]; exists {
		theme.tetrisL = resolveColor(cv)
	} else {
		theme.tetrisL = lipgloss.Color("#ffa500")
	}

	return theme, nil
}
