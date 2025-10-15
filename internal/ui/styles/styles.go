package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Adaptive colors work with light/dark terminals
	Primary   = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}
	Secondary = lipgloss.AdaptiveColor{Light: "#585858", Dark: "#a8a8a8"}
	Accent    = lipgloss.AdaptiveColor{Light: "#0066cc", Dark: "#66b3ff"}
	Success   = lipgloss.AdaptiveColor{Light: "#22c55e", Dark: "#4ade80"}
)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			MarginBottom(1)

	MenuItemStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Accent).
				Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			MarginTop(1)
)

// Game board styles
var (
	CellStyle = lipgloss.NewStyle().
			Width(3).
			Height(1).
			Align(lipgloss.Center, lipgloss.Center).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary)

	SelectedCellStyle = CellStyle.
				BorderForeground(Accent).
				BorderStyle(lipgloss.ThickBorder())

	BorderStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	GameOverStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6b6b")).
			Bold(true)

	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(1).
			Width(15)
)

// Piece/game element styles
var (
	WhitePieceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Bold(true)

	BlackPieceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444")).
			Bold(true)

	SnakeStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	SnakeHeadStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)
)
