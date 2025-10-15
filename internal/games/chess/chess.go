package chess

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	board  [][]int
	width  int
	height int
}

type Position struct {
	x, y int
}

func New() *Model {
	return &Model{}
}

func (m *Model) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	return "Chess Game, Coming soon"
}

type tickMsg struct{}
