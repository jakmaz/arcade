package theme

import "github.com/charmbracelet/lipgloss"

// NewDefaultTheme creates a default theme with sensible colors
func NewDefaultTheme() Theme {
	return &BaseTheme{
		name: "default",

		// UI Colors
		primary:   lipgloss.Color("#ffffff"),
		secondary: lipgloss.Color("#888888"),
		accent:    lipgloss.Color("#0066cc"),
		success:   lipgloss.Color("#22c55e"),
		warning:   lipgloss.Color("#f59e0b"),
		error:     lipgloss.Color("#ef4444"),

		// Game Board Colors
		boardBorder:     lipgloss.Color("#888888"),
		boardBackground: lipgloss.Color(""),
		cellBorder:      lipgloss.Color("#888888"),
		cellBackground:  lipgloss.Color(""),
		selectedCell:    lipgloss.Color("#0066cc"),

		// Game Piece Colors
		player1:   lipgloss.Color("#22c55e"),
		player2:   lipgloss.Color("#ef4444"),
		snakeBody: lipgloss.Color("#22c55e"),
		snakeHead: lipgloss.Color("#0066cc"),
		food:      lipgloss.Color("#ef4444"),

		// Chess Piece Colors
		whitePiece: lipgloss.Color("#ffffff"),
		blackPiece: lipgloss.Color("#444444"),

		// Tetris Colors
		tetrisI: lipgloss.Color("#00f5ff"),
		tetrisO: lipgloss.Color("#ffff00"),
		tetrisT: lipgloss.Color("#800080"),
		tetrisS: lipgloss.Color("#00ff00"),
		tetrisZ: lipgloss.Color("#ff0000"),
		tetrisJ: lipgloss.Color("#0000ff"),
		tetrisL: lipgloss.Color("#ffa500"),

		// Terminal Background
		terminalBackground:    lipgloss.Color("#1a1a1a"),
		useTerminalBackground: true,
	}
}

// NewSystemTheme creates a system theme that adapts to terminal background
func NewSystemTheme() Theme {
	return &BaseTheme{
		name: "system",

		// UI Colors - using adaptive colors
		primary:   lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"},
		secondary: lipgloss.AdaptiveColor{Light: "#585858", Dark: "#a8a8a8"},
		accent:    lipgloss.AdaptiveColor{Light: "#0066cc", Dark: "#66b3ff"},
		success:   lipgloss.AdaptiveColor{Light: "#22c55e", Dark: "#4ade80"},
		warning:   lipgloss.AdaptiveColor{Light: "#f59e0b", Dark: "#fbbf24"},
		error:     lipgloss.AdaptiveColor{Light: "#ef4444", Dark: "#f87171"},

		// Game Board Colors
		boardBorder:     lipgloss.AdaptiveColor{Light: "#585858", Dark: "#a8a8a8"},
		boardBackground: lipgloss.Color(""),
		cellBorder:      lipgloss.AdaptiveColor{Light: "#585858", Dark: "#a8a8a8"},
		cellBackground:  lipgloss.Color(""),
		selectedCell:    lipgloss.AdaptiveColor{Light: "#0066cc", Dark: "#66b3ff"},

		// Game Piece Colors
		player1:   lipgloss.AdaptiveColor{Light: "#22c55e", Dark: "#4ade80"},
		player2:   lipgloss.AdaptiveColor{Light: "#ef4444", Dark: "#f87171"},
		snakeBody: lipgloss.AdaptiveColor{Light: "#22c55e", Dark: "#4ade80"},
		snakeHead: lipgloss.AdaptiveColor{Light: "#0066cc", Dark: "#66b3ff"},
		food:      lipgloss.AdaptiveColor{Light: "#ef4444", Dark: "#f87171"},

		// Chess Piece Colors
		whitePiece: lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#ffffff"},
		blackPiece: lipgloss.AdaptiveColor{Light: "#666666", Dark: "#444444"},

		// Tetris Colors
		tetrisI: lipgloss.AdaptiveColor{Light: "#0891b2", Dark: "#00f5ff"},
		tetrisO: lipgloss.AdaptiveColor{Light: "#eab308", Dark: "#ffff00"},
		tetrisT: lipgloss.AdaptiveColor{Light: "#7c3aed", Dark: "#a855f7"},
		tetrisS: lipgloss.AdaptiveColor{Light: "#16a34a", Dark: "#22c55e"},
		tetrisZ: lipgloss.AdaptiveColor{Light: "#dc2626", Dark: "#ef4444"},
		tetrisJ: lipgloss.AdaptiveColor{Light: "#2563eb", Dark: "#3b82f6"},
		tetrisL: lipgloss.AdaptiveColor{Light: "#ea580c", Dark: "#f97316"},

		// Terminal Background - disabled for system theme
		terminalBackground:    lipgloss.Color(""),
		useTerminalBackground: false,
	}
}
