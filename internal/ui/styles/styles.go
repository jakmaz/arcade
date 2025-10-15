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
