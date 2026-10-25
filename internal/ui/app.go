package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	app := NewApp()
	if game := core.CreateGame(gameID); game != nil {
		app.state = GameState
		app.currentGame = game
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
		// Keep the menu sized while a game is open, ready for Esc.
		var menuCmd tea.Cmd
		a.menu, menuCmd = a.menu.Update(msg)
		if a.state == MenuState {
			return a, menuCmd
		}
		var gameCmd tea.Cmd
		a.currentGame, gameCmd = a.currentGame.Update(msg)
		return a, tea.Batch(menuCmd, gameCmd)

	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return a, tea.Quit
		}

	case StartGameMsg:
		game := core.CreateGame(msg.GameID)
		if game == nil {
			return a, nil
		}
		a.currentGame = game
		a.state = GameState

		sizeMsg := tea.WindowSizeMsg{Width: a.width, Height: a.height}
		var sizeCmd tea.Cmd
		a.currentGame, sizeCmd = a.currentGame.Update(sizeMsg)
		return a, tea.Batch(sizeCmd, a.currentGame.Init())

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

func (a *App) View() tea.View {
	var v tea.View

	// Get child view
	if a.state == MenuState {
		v = a.menu.View()
	} else {
		v = a.currentGame.View()
	}

	v.AltScreen = true

	bg := styles.GetTerminalBackgroundStyle().GetBackground()
	if _, isNoColor := bg.(lipgloss.NoColor); !isNoColor {
		v.BackgroundColor = bg
	}

	return v
}

func (a *App) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.menu, cmd = a.menu.Update(msg)

	return a, cmd
}

func (a *App) updateGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for ESC to return to menu
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok && keyMsg.String() == "esc" {
		return a.Update(ReturnToMenuMsg{})
	}

	var cmd tea.Cmd
	a.currentGame, cmd = a.currentGame.Update(msg)
	return a, cmd
}
