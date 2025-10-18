package theme

import (
	"image/color"
)

// Theme represents a complete color scheme for the arcade games
type Theme interface {
	Name() string

	// UI Colors
	Primary() color.Color
	Secondary() color.Color
	Accent() color.Color
	Success() color.Color
	Warning() color.Color
	Error() color.Color

	// Game Board Colors
	BoardBorder() color.Color
	BoardBackground() color.Color
	CellBorder() color.Color
	CellBackground() color.Color
	SelectedCell() color.Color

	// Game Piece Colors
	Player1() color.Color
	Player2() color.Color
	SnakeBody() color.Color
	SnakeHead() color.Color
	Food() color.Color

	// Chess Piece Colors
	WhitePiece() color.Color
	BlackPiece() color.Color

	// Tetris Colors
	TetrisI() color.Color
	TetrisO() color.Color
	TetrisT() color.Color
	TetrisS() color.Color
	TetrisZ() color.Color
	TetrisJ() color.Color
	TetrisL() color.Color

	// Terminal Background
	TerminalBackground() color.Color
	ShouldUseTerminalBackground() bool
}

// BaseTheme provides default implementations
type BaseTheme struct {
	name string

	// Color definitions
	primary   color.Color
	secondary color.Color
	accent    color.Color
	success   color.Color
	warning   color.Color
	error     color.Color

	boardBorder     color.Color
	boardBackground color.Color
	cellBorder      color.Color
	cellBackground  color.Color
	selectedCell    color.Color

	player1   color.Color
	player2   color.Color
	snakeBody color.Color
	snakeHead color.Color
	food      color.Color

	whitePiece color.Color
	blackPiece color.Color

	tetrisI color.Color
	tetrisO color.Color
	tetrisT color.Color
	tetrisS color.Color
	tetrisZ color.Color
	tetrisJ color.Color
	tetrisL color.Color

	terminalBackground    color.Color
	useTerminalBackground bool
}

func (t *BaseTheme) Name() string { return t.name }

// UI Colors
func (t *BaseTheme) Primary() color.Color   { return t.primary }
func (t *BaseTheme) Secondary() color.Color { return t.secondary }
func (t *BaseTheme) Accent() color.Color    { return t.accent }
func (t *BaseTheme) Success() color.Color   { return t.success }
func (t *BaseTheme) Warning() color.Color   { return t.warning }
func (t *BaseTheme) Error() color.Color     { return t.error }

// Game Board Colors
func (t *BaseTheme) BoardBorder() color.Color     { return t.boardBorder }
func (t *BaseTheme) BoardBackground() color.Color { return t.boardBackground }
func (t *BaseTheme) CellBorder() color.Color      { return t.cellBorder }
func (t *BaseTheme) CellBackground() color.Color  { return t.cellBackground }
func (t *BaseTheme) SelectedCell() color.Color    { return t.selectedCell }

// Game Piece Colors
func (t *BaseTheme) Player1() color.Color   { return t.player1 }
func (t *BaseTheme) Player2() color.Color   { return t.player2 }
func (t *BaseTheme) SnakeBody() color.Color { return t.snakeBody }
func (t *BaseTheme) SnakeHead() color.Color { return t.snakeHead }
func (t *BaseTheme) Food() color.Color      { return t.food }

// Chess Piece Colors
func (t *BaseTheme) WhitePiece() color.Color { return t.whitePiece }
func (t *BaseTheme) BlackPiece() color.Color { return t.blackPiece }

// Tetris Colors
func (t *BaseTheme) TetrisI() color.Color { return t.tetrisI }
func (t *BaseTheme) TetrisO() color.Color { return t.tetrisO }
func (t *BaseTheme) TetrisT() color.Color { return t.tetrisT }
func (t *BaseTheme) TetrisS() color.Color { return t.tetrisS }
func (t *BaseTheme) TetrisZ() color.Color { return t.tetrisZ }
func (t *BaseTheme) TetrisJ() color.Color { return t.tetrisJ }
func (t *BaseTheme) TetrisL() color.Color { return t.tetrisL }

// Terminal Background
func (t *BaseTheme) TerminalBackground() color.Color   { return t.terminalBackground }
func (t *BaseTheme) ShouldUseTerminalBackground() bool { return t.useTerminalBackground }
