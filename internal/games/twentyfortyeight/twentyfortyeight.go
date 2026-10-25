package twentyfortyeight

import (
	"math/rand/v2"

	tea "charm.land/bubbletea/v2"
)

const boardSize = 4

type board [boardSize][boardSize]int

type direction uint8

const (
	left direction = iota
	right
	up
	down
)

type Model struct {
	board          board
	score          int
	won            bool
	continued      bool
	gameOver       bool
	width, height  int
	darkBackground bool
	random         func(int) int
}

func New() *Model {
	m := &Model{random: rand.IntN, darkBackground: true}
	m.reset()
	return m
}

func (m *Model) Init() tea.Cmd {
	return tea.RequestBackgroundColor
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.BackgroundColorMsg:
		m.darkBackground = msg.IsDark()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "r", "R":
			m.reset()
		case "enter", "space":
			if m.won {
				m.continued = true
			}
		case "left", "right", "up", "down":
			if m.width < m.viewWidth() || m.height < viewHeight {
				return m, nil
			}
			switch msg.String() {
			case "left":
				m.move(left)
			case "right":
				m.move(right)
			case "up":
				m.move(up)
			case "down":
				m.move(down)
			}
		}
	}
	return m, nil
}

func (m *Model) reset() {
	m.board = board{}
	m.score = 0
	m.won, m.continued, m.gameOver = false, false, false
	m.spawnTile()
	m.spawnTile()
}

// mergeLine packs tiles toward the start, merging each source tile once.
func mergeLine(line [boardSize]int) ([boardSize]int, int) {
	var packed, merged [boardSize]int
	count := 0
	for _, value := range line {
		if value != 0 {
			packed[count] = value
			count++
		}
	}

	score, destination := 0, 0
	for source := 0; source < count; source++ {
		value := packed[source]
		if source+1 < count && value == packed[source+1] {
			value *= 2
			score += value
			source++
		}
		merged[destination] = value
		destination++
	}
	return merged, score
}

// linePosition visits a row or column from the edge tiles are moving toward.
func linePosition(d direction, line, offset int) (row, column int) {
	switch d {
	case right:
		return line, boardSize - 1 - offset
	case up:
		return offset, line
	case down:
		return boardSize - 1 - offset, line
	default:
		return line, offset
	}
}

func slide(before board, d direction) (after board, score int, changed bool) {
	for line := range boardSize {
		var values [boardSize]int
		for offset := range boardSize {
			row, column := linePosition(d, line, offset)
			values[offset] = before[row][column]
		}
		merged, gain := mergeLine(values)
		score += gain
		for offset, value := range merged {
			row, column := linePosition(d, line, offset)
			after[row][column] = value
		}
	}
	return after, score, before != after
}

func (m *Model) move(d direction) bool {
	if m.gameOver || (m.won && !m.continued) {
		return false
	}
	result, gain, changed := slide(m.board, d)
	if !changed {
		m.gameOver = !hasMoves(m.board)
		return false
	}
	m.board = result
	m.score += gain
	m.spawnTile()
	for _, row := range m.board {
		for _, value := range row {
			if value >= 2048 {
				m.won = true
			}
		}
	}
	m.gameOver = !hasMoves(m.board)
	return true
}

func (m *Model) spawnTile() bool {
	var empty [boardSize * boardSize]int
	count := 0
	for row := range boardSize {
		for column := range boardSize {
			if m.board[row][column] == 0 {
				empty[count] = row*boardSize + column
				count++
			}
		}
	}
	if count == 0 {
		return false
	}
	position := empty[m.random(count)]
	value := 2
	if m.random(10) == 0 {
		value = 4
	}
	m.board[position/boardSize][position%boardSize] = value
	return true
}

func hasMoves(b board) bool {
	for row := range boardSize {
		for column := range boardSize {
			value := b[row][column]
			if value == 0 ||
				(row+1 < boardSize && value == b[row+1][column]) ||
				(column+1 < boardSize && value == b[row][column+1]) {
				return true
			}
		}
	}
	return false
}
