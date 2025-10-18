package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jakmaz/arcade/internal/theme"
)

var (
	// Get the themed styles (will be initialized lazily)
	styles *theme.Styles

	// Legacy adaptive colors for backwards compatibility
	Primary   = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}
	Secondary = lipgloss.AdaptiveColor{Light: "#585858", Dark: "#a8a8a8"}
	Accent    = lipgloss.AdaptiveColor{Light: "#0066cc", Dark: "#66b3ff"}
	Success   = lipgloss.AdaptiveColor{Light: "#22c55e", Dark: "#4ade80"}
)

// Lazy initialization variables
var (
	TitleStyle        lipgloss.Style
	MenuItemStyle     lipgloss.Style
	SelectedItemStyle lipgloss.Style
	HelpStyle         lipgloss.Style
	CellStyle         lipgloss.Style
	SelectedCellStyle lipgloss.Style
	BorderStyle       lipgloss.Style
	GameOverStyle     lipgloss.Style
	SidebarStyle      lipgloss.Style
	WhitePieceStyle   lipgloss.Style
	BlackPieceStyle   lipgloss.Style
	SnakeStyle        lipgloss.Style
	SnakeHeadStyle    lipgloss.Style
)

// ensureInitialized makes sure themes are loaded and styles are set
func ensureInitialized() {
	if styles == nil {
		// Initialize themes first
		theme.Initialize()
		styles = theme.NewStyles()

		// Set all style variables
		TitleStyle = styles.TitleStyle()
		MenuItemStyle = styles.MenuItemStyle()
		SelectedItemStyle = styles.SelectedItemStyle()
		HelpStyle = styles.HelpStyle()

		CellStyle = styles.CellStyle()
		SelectedCellStyle = styles.SelectedCellStyle()
		BorderStyle = styles.BorderStyle()
		GameOverStyle = styles.GameOverStyle()
		SidebarStyle = styles.SidebarStyle()

		WhitePieceStyle = styles.WhitePieceStyle()
		BlackPieceStyle = styles.BlackPieceStyle()
		SnakeStyle = styles.SnakeStyle()
		SnakeHeadStyle = styles.SnakeHeadStyle()
	}
}

// GetTitleStyle returns the title style, initializing if needed
func GetTitleStyle() lipgloss.Style {
	ensureInitialized()
	return TitleStyle
}

// GetMenuItemStyle returns the menu item style, initializing if needed
func GetMenuItemStyle() lipgloss.Style {
	ensureInitialized()
	return MenuItemStyle
}

// GetSelectedItemStyle returns the selected item style, initializing if needed
func GetSelectedItemStyle() lipgloss.Style {
	ensureInitialized()
	return SelectedItemStyle
}

// GetHelpStyle returns the help style, initializing if needed
func GetHelpStyle() lipgloss.Style {
	ensureInitialized()
	return HelpStyle
}

// GetCellStyle returns the cell style, initializing if needed
func GetCellStyle() lipgloss.Style {
	ensureInitialized()
	return CellStyle
}

// GetSelectedCellStyle returns the selected cell style, initializing if needed
func GetSelectedCellStyle() lipgloss.Style {
	ensureInitialized()
	return SelectedCellStyle
}

// GetBorderStyle returns the border style, initializing if needed
func GetBorderStyle() lipgloss.Style {
	ensureInitialized()
	return BorderStyle
}

// GetGameOverStyle returns the game over style, initializing if needed
func GetGameOverStyle() lipgloss.Style {
	ensureInitialized()
	return GameOverStyle
}

// GetSidebarStyle returns the sidebar style, initializing if needed
func GetSidebarStyle() lipgloss.Style {
	ensureInitialized()
	return SidebarStyle
}

// GetWhitePieceStyle returns the white piece style, initializing if needed
func GetWhitePieceStyle() lipgloss.Style {
	ensureInitialized()
	return WhitePieceStyle
}

// GetBlackPieceStyle returns the black piece style, initializing if needed
func GetBlackPieceStyle() lipgloss.Style {
	ensureInitialized()
	return BlackPieceStyle
}

// GetSnakeStyle returns the snake style, initializing if needed
func GetSnakeStyle() lipgloss.Style {
	ensureInitialized()
	return SnakeStyle
}

// GetSnakeHeadStyle returns the snake head style, initializing if needed
func GetSnakeHeadStyle() lipgloss.Style {
	ensureInitialized()
	return SnakeHeadStyle
}

// RefreshStyles updates all styles with the current theme
func RefreshStyles() {
	styles = theme.NewStyles()

	// Update all style variables
	TitleStyle = styles.TitleStyle()
	MenuItemStyle = styles.MenuItemStyle()
	SelectedItemStyle = styles.SelectedItemStyle()
	HelpStyle = styles.HelpStyle()

	CellStyle = styles.CellStyle()
	SelectedCellStyle = styles.SelectedCellStyle()
	BorderStyle = styles.BorderStyle()
	GameOverStyle = styles.GameOverStyle()
	SidebarStyle = styles.SidebarStyle()

	WhitePieceStyle = styles.WhitePieceStyle()
	BlackPieceStyle = styles.BlackPieceStyle()
	SnakeStyle = styles.SnakeStyle()
	SnakeHeadStyle = styles.SnakeHeadStyle()
}

// GetStyles returns the current themed styles instance
func GetStyles() *theme.Styles {
	ensureInitialized()
	return styles
}
