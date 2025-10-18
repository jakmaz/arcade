package theme

import "github.com/charmbracelet/lipgloss/v2"

// SystemTheme creates a system theme that adapts to terminal background
func SystemTheme() Theme {
	return &BaseTheme{
		name: "system",

		// UI Colors - using dark colors as default for now
		primary:   lipgloss.Color("#dddddd"),
		secondary: lipgloss.Color("#a8a8a8"),
		accent:    lipgloss.Color("#66b3ff"),
		success:   lipgloss.Color("#4ade80"),
		warning:   lipgloss.Color("#fbbf24"),
		error:     lipgloss.Color("#f87171"),

		// Game Board Colors
		boardBorder:     lipgloss.Color("#a8a8a8"),
		boardBackground: lipgloss.Color(""),
		cellBorder:      lipgloss.Color("#a8a8a8"),
		cellBackground:  lipgloss.Color(""),
		selectedCell:    lipgloss.Color("#66b3ff"),

		// Game Piece Colors
		player1:   lipgloss.Color("#4ade80"),
		player2:   lipgloss.Color("#f87171"),
		snakeBody: lipgloss.Color("#4ade80"),
		snakeHead: lipgloss.Color("#66b3ff"),
		food:      lipgloss.Color("#f87171"),

		// Chess Piece Colors
		whitePiece: lipgloss.Color("#ffffff"),
		blackPiece: lipgloss.Color("#444444"),

		// Tetris Colors
		tetrisI: lipgloss.Color("#00f5ff"),
		tetrisO: lipgloss.Color("#ffff00"),
		tetrisT: lipgloss.Color("#a855f7"),
		tetrisS: lipgloss.Color("#22c55e"),
		tetrisZ: lipgloss.Color("#ef4444"),
		tetrisJ: lipgloss.Color("#3b82f6"),
		tetrisL: lipgloss.Color("#f97316"),

		// Terminal Background - disabled for system theme
		terminalBackground:    lipgloss.Color(""),
		useTerminalBackground: false,
	}
}
