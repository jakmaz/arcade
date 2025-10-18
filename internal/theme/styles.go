package theme

import "github.com/charmbracelet/lipgloss"

// Styles provides lipgloss styles based on the current theme
type Styles struct {
	theme Theme
}

// NewStyles creates a new Styles instance using the current theme
func NewStyles() *Styles {
	return &Styles{theme: GetCurrentTheme()}
}

// NewStylesWithTheme creates a new Styles instance using a specific theme
func NewStylesWithTheme(theme Theme) *Styles {
	return &Styles{theme: theme}
}

// UI Styles
func (s *Styles) TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Primary()).
		Bold(true).
		MarginBottom(1)
}

func (s *Styles) MenuItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Secondary())
}

func (s *Styles) SelectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Accent()).
		Bold(true)
}

func (s *Styles) HelpStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Secondary()).
		MarginTop(1)
}

// Game Board Styles
func (s *Styles) CellStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(3).
		Height(1).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.theme.CellBorder())
}

func (s *Styles) SelectedCellStyle() lipgloss.Style {
	return s.CellStyle().
		BorderForeground(s.theme.SelectedCell()).
		BorderStyle(lipgloss.ThickBorder())
}

func (s *Styles) BorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.BoardBorder()).
		Bold(true)
}

func (s *Styles) GameOverStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Error()).
		Bold(true)
}

func (s *Styles) SidebarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.theme.BoardBorder()).
		Padding(1).
		Width(15)
}

// Game Piece Styles
func (s *Styles) WhitePieceStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.WhitePiece()).
		Bold(true)
}

func (s *Styles) BlackPieceStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.BlackPiece()).
		Bold(true)
}

func (s *Styles) SnakeStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.SnakeBody()).
		Bold(true)
}

func (s *Styles) SnakeHeadStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.SnakeHead()).
		Bold(true)
}

func (s *Styles) FoodStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Food()).
		Bold(true)
}

// Tetris Piece Styles
func (s *Styles) TetrisPieceStyle(pieceType string) lipgloss.Style {
	var color lipgloss.TerminalColor

	switch pieceType {
	case "I":
		color = s.theme.TetrisI()
	case "O":
		color = s.theme.TetrisO()
	case "T":
		color = s.theme.TetrisT()
	case "S":
		color = s.theme.TetrisS()
	case "Z":
		color = s.theme.TetrisZ()
	case "J":
		color = s.theme.TetrisJ()
	case "L":
		color = s.theme.TetrisL()
	default:
		color = s.theme.Primary()
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true)
}

// Player Styles
func (s *Styles) Player1Style() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Player1()).
		Bold(true)
}

func (s *Styles) Player2Style() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Player2()).
		Bold(true)
}

// Success/Warning/Error Styles
func (s *Styles) SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Success()).
		Bold(true)
}

func (s *Styles) WarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Warning()).
		Bold(true)
}

func (s *Styles) ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(s.theme.Error()).
		Bold(true)
}
