package theme

import "charm.land/lipgloss/v2"

// SystemTheme creates a system theme that adapts to terminal colors using ANSI colors
func SystemTheme() Theme {
	return &BaseTheme{
		name: "system",

		// UI Colors - using ANSI colors that adapt to terminal theme
		primary:   lipgloss.Color(""),  // default foreground
		secondary: lipgloss.Color("8"), // bright black (gray)
		accent:    lipgloss.Color("4"), // blue
		success:   lipgloss.Color("2"), // green
		warning:   lipgloss.Color("3"), // yellow
		error:     lipgloss.Color("1"), // red

		// Game Board Colors
		boardBorder:     lipgloss.Color("8"), // bright black (gray)
		boardBackground: lipgloss.Color(""),  // transparent (use terminal bg)
		cellBorder:      lipgloss.Color("8"), // bright black (gray)
		cellBackground:  lipgloss.Color(""),  // transparent
		selectedCell:    lipgloss.Color("4"), // blue

		// Game Piece Colors
		player1:   lipgloss.Color("2"), // green
		player2:   lipgloss.Color("1"), // red
		snakeBody: lipgloss.Color("2"), // green
		snakeHead: lipgloss.Color("4"), // blue
		food:      lipgloss.Color("1"), // red

		// Chess Piece Colors
		whitePiece: lipgloss.Color("15"), // bright white
		blackPiece: lipgloss.Color("0"),  // black

		// Tetris Colors - using bright ANSI colors for visibility
		tetrisI: lipgloss.Color("6"),  // cyan
		tetrisO: lipgloss.Color("3"),  // yellow
		tetrisT: lipgloss.Color("5"),  // magenta
		tetrisS: lipgloss.Color("2"),  // green
		tetrisZ: lipgloss.Color("1"),  // red
		tetrisJ: lipgloss.Color("4"),  // blue
		tetrisL: lipgloss.Color("11"), // bright yellow

		// Terminal Background - empty for system theme (use terminal's own bg)
		terminalBackground: lipgloss.Color(""),
	}
}
