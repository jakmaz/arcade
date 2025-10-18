package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a complete color scheme for the arcade games
type Theme interface {
	Name() string

	// UI Colors
	Primary() lipgloss.TerminalColor
	Secondary() lipgloss.TerminalColor
	Accent() lipgloss.TerminalColor
	Success() lipgloss.TerminalColor
	Warning() lipgloss.TerminalColor
	Error() lipgloss.TerminalColor

	// Game Board Colors
	BoardBorder() lipgloss.TerminalColor
	BoardBackground() lipgloss.TerminalColor
	CellBorder() lipgloss.TerminalColor
	CellBackground() lipgloss.TerminalColor
	SelectedCell() lipgloss.TerminalColor

	// Game Piece Colors
	Player1() lipgloss.TerminalColor
	Player2() lipgloss.TerminalColor
	SnakeBody() lipgloss.TerminalColor
	SnakeHead() lipgloss.TerminalColor
	Food() lipgloss.TerminalColor

	// Chess Piece Colors
	WhitePiece() lipgloss.TerminalColor
	BlackPiece() lipgloss.TerminalColor

	// Tetris Colors
	TetrisI() lipgloss.TerminalColor
	TetrisO() lipgloss.TerminalColor
	TetrisT() lipgloss.TerminalColor
	TetrisS() lipgloss.TerminalColor
	TetrisZ() lipgloss.TerminalColor
	TetrisJ() lipgloss.TerminalColor
	TetrisL() lipgloss.TerminalColor
}

// BaseTheme provides default implementations
type BaseTheme struct {
	name string

	// Color definitions
	primary   lipgloss.TerminalColor
	secondary lipgloss.TerminalColor
	accent    lipgloss.TerminalColor
	success   lipgloss.TerminalColor
	warning   lipgloss.TerminalColor
	error     lipgloss.TerminalColor

	boardBorder     lipgloss.TerminalColor
	boardBackground lipgloss.TerminalColor
	cellBorder      lipgloss.TerminalColor
	cellBackground  lipgloss.TerminalColor
	selectedCell    lipgloss.TerminalColor

	player1   lipgloss.TerminalColor
	player2   lipgloss.TerminalColor
	snakeBody lipgloss.TerminalColor
	snakeHead lipgloss.TerminalColor
	food      lipgloss.TerminalColor

	whitePiece lipgloss.TerminalColor
	blackPiece lipgloss.TerminalColor

	tetrisI lipgloss.TerminalColor
	tetrisO lipgloss.TerminalColor
	tetrisT lipgloss.TerminalColor
	tetrisS lipgloss.TerminalColor
	tetrisZ lipgloss.TerminalColor
	tetrisJ lipgloss.TerminalColor
	tetrisL lipgloss.TerminalColor
}

func (t *BaseTheme) Name() string { return t.name }

// UI Colors
func (t *BaseTheme) Primary() lipgloss.TerminalColor   { return t.primary }
func (t *BaseTheme) Secondary() lipgloss.TerminalColor { return t.secondary }
func (t *BaseTheme) Accent() lipgloss.TerminalColor    { return t.accent }
func (t *BaseTheme) Success() lipgloss.TerminalColor   { return t.success }
func (t *BaseTheme) Warning() lipgloss.TerminalColor   { return t.warning }
func (t *BaseTheme) Error() lipgloss.TerminalColor     { return t.error }

// Game Board Colors
func (t *BaseTheme) BoardBorder() lipgloss.TerminalColor     { return t.boardBorder }
func (t *BaseTheme) BoardBackground() lipgloss.TerminalColor { return t.boardBackground }
func (t *BaseTheme) CellBorder() lipgloss.TerminalColor      { return t.cellBorder }
func (t *BaseTheme) CellBackground() lipgloss.TerminalColor  { return t.cellBackground }
func (t *BaseTheme) SelectedCell() lipgloss.TerminalColor    { return t.selectedCell }

// Game Piece Colors
func (t *BaseTheme) Player1() lipgloss.TerminalColor   { return t.player1 }
func (t *BaseTheme) Player2() lipgloss.TerminalColor   { return t.player2 }
func (t *BaseTheme) SnakeBody() lipgloss.TerminalColor { return t.snakeBody }
func (t *BaseTheme) SnakeHead() lipgloss.TerminalColor { return t.snakeHead }
func (t *BaseTheme) Food() lipgloss.TerminalColor      { return t.food }

// Chess Piece Colors
func (t *BaseTheme) WhitePiece() lipgloss.TerminalColor { return t.whitePiece }
func (t *BaseTheme) BlackPiece() lipgloss.TerminalColor { return t.blackPiece }

// Tetris Colors
func (t *BaseTheme) TetrisI() lipgloss.TerminalColor { return t.tetrisI }
func (t *BaseTheme) TetrisO() lipgloss.TerminalColor { return t.tetrisO }
func (t *BaseTheme) TetrisT() lipgloss.TerminalColor { return t.tetrisT }
func (t *BaseTheme) TetrisS() lipgloss.TerminalColor { return t.tetrisS }
func (t *BaseTheme) TetrisZ() lipgloss.TerminalColor { return t.tetrisZ }
func (t *BaseTheme) TetrisJ() lipgloss.TerminalColor { return t.tetrisJ }
func (t *BaseTheme) TetrisL() lipgloss.TerminalColor { return t.tetrisL }
