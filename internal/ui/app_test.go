package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

type menuCommandModel struct {
	calls *int
}

type menuCommandMsg struct{}

func (m menuCommandModel) Init() tea.Cmd  { return nil }
func (m menuCommandModel) View() tea.View { return tea.NewView("") }
func (m menuCommandModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		(*m.calls)++
		return menuCommandMsg{}
	}
}

func TestMenuCommandsAreDeliveredByBubbleTea(t *testing.T) {
	calls := 0
	app := &App{menu: menuCommandModel{calls: &calls}}
	_, cmd := app.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if calls != 0 {
		t.Fatal("Update executed the menu command synchronously")
	}
	if cmd == nil {
		t.Fatal("menu command was discarded")
	}
	app.Update(cmd())
	if calls != 1 {
		t.Fatalf("menu command ran %d times, want once", calls)
	}
}

func TestReturningToMenuAfterResize(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	app := NewAppWithGame("snake")
	app.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	app.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if app.state != MenuState || app.currentGame != nil {
		t.Fatal("Esc did not return to menu")
	}
	menu, ok := app.menu.(model)
	if !ok || menu.width != 100 || menu.height != 30 {
		t.Fatalf("returned menu has stale dimensions: %#v", app.menu)
	}
	app.View()
}

func TestUnknownGameKeepsMenu(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	app := NewAppWithGame("missing")
	app.Update(StartGameMsg{GameID: "missing"})
	if app.state != MenuState || app.currentGame != nil {
		t.Fatal("unknown game left the app in an invalid game state")
	}
	app.Init()
	app.View()
}

func TestQuitFromGame(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	app := NewAppWithGame("snake")
	_, cmd := app.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	if cmd == nil {
		t.Fatal("q did not quit from game")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("q returned a non-quit command")
	}
}
