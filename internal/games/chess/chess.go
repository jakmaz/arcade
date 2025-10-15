package chess

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *Model) View() string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "Chess Game, coming Soon")
}

type tickMsg struct{}
