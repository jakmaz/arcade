package tictactoe

import tea "github.com/charmbracelet/bubbletea"

func New() *Model {
	return &Model{}
}

type Model struct {
	board            [3][3]rune
	turn             rune
	cursorX, cursorY int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "Test"
}
