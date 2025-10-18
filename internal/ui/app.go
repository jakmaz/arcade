package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jakmaz/arcade/internal/core"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

type AppState int

const (
	MenuState AppState = iota
	GameState
)

type App struct {
	state       AppState
	menu        tea.Model
	currentGame tea.Model
	width       int
	height      int
}

func NewApp() *App {
	return &App{
		state: MenuState,
		menu:  NewMenu(),
	}
}

func NewAppWithGame(gameID string) *App {
	app := &App{
		state:       GameState,
		currentGame: core.CreateGame(gameID),
		width:       80,
		height:      25,
	}
	return app
}

// Custom messages for navigation
type StartGameMsg struct {
	GameID string
}

type ReturnToMenuMsg struct{}

func (a *App) Init() tea.Cmd {
	if a.state == GameState && a.currentGame != nil {
		return a.currentGame.Init()
	}
	return a.menu.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		// Forward to current model
		if a.state == MenuState {
			var cmd tea.Cmd
			a.menu, cmd = a.menu.Update(msg)
			return a, cmd
		} else {
			var cmd tea.Cmd
			a.currentGame, cmd = a.currentGame.Update(msg)
			return a, cmd
		}

	case tea.KeyMsg:
		// Global Ctrl+C handling
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}

	case StartGameMsg:
		// Transition to game
		a.currentGame = core.CreateGame(msg.GameID)
		a.state = GameState

		sizeMsg := tea.WindowSizeMsg{Width: a.width, Height: a.height}
		a.currentGame, _ = a.currentGame.Update(sizeMsg)
		return a, a.currentGame.Init()

	case ReturnToMenuMsg:
		// Transition back to menu
		a.state = MenuState
		a.currentGame = nil
		return a, nil

	case ThemeChangedMsg:
		styles.RefreshStyles()
		return a, nil
	}

	// Delegate to current state
	if a.state == MenuState {
		return a.updateMenu(msg)
	} else {
		return a.updateGame(msg)
	}
}

func (a *App) View() string {
	var content string
	if a.state == MenuState {
		content = a.menu.View()
	} else {
		content = a.currentGame.View()
	}

	// Apply terminal background with full viewport dimensions
	if a.width > 0 && a.height > 0 {
		return styles.GetTerminalBackgroundStyle().
			Width(a.width).
			Height(a.height).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(content)
	}

	return content
}

func (a *App) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.menu, cmd = a.menu.Update(msg)

	// Check if menu wants to start a game
	if cmd != nil {
		if msg := cmd(); msg != nil {
			if startMsg, ok := msg.(StartGameMsg); ok {
				return a.Update(startMsg)
			}
			if themeMsg, ok := msg.(ThemeChangedMsg); ok {
				return a.Update(themeMsg)
			}
		}
	}

	return a, cmd
}

func (a *App) updateGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for ESC to return to menu
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return a.Update(ReturnToMenuMsg{})
	}

	var cmd tea.Cmd
	a.currentGame, cmd = a.currentGame.Update(msg)
	return a, cmd
}
