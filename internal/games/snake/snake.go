package snake

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

const (
	boardWidth  = 28
	boardHeight = 16
	// A tile spans two columns so horizontal and vertical moves look similar.
	cellWidth    = 2
	viewWidth    = boardWidth*cellWidth + 2
	viewHeight   = boardHeight + 7
	stepInterval = 125 * time.Millisecond
)

type Position struct {
	x, y int
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Model struct {
	snake         []Position
	food          Position
	direction     Direction
	nextDirection Direction
	score         int
	paused        bool
	gameOver      bool
	won           bool
	width, height int
	tickID        uint64
}

type tickMsg struct {
	game *Model
	id   uint64
}

func New() *Model {
	m := &Model{}
	m.reset()
	return m
}

func (m *Model) reset() {
	m.snake = []Position{{14, 8}, {13, 8}, {12, 8}}
	m.food = Position{20, 8}
	m.direction = Right
	m.nextDirection = Right
	m.score = 0
	m.paused = false
	m.gameOver = false
	m.won = false
	m.tickID++
}

func (m *Model) tick() tea.Cmd {
	id := m.tickID
	return tea.Tick(stepInterval, func(time.Time) tea.Msg {
		return tickMsg{game: m, id: id}
	})
}

func (m *Model) Init() tea.Cmd {
	return m.tick()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.turn(Up)
		case "down":
			m.turn(Down)
		case "left":
			m.turn(Left)
		case "right":
			m.turn(Right)
		case "space":
			if !m.gameOver {
				m.paused = !m.paused
				m.tickID++
				if !m.paused {
					return m, m.tick()
				}
			}
		case "r":
			m.reset()
			return m, m.tick()
		}
	case tickMsg:
		if msg.game != m || msg.id != m.tickID || m.paused || m.gameOver {
			return m, nil
		}
		if m.width >= viewWidth && m.height >= viewHeight {
			m.advance()
		}
		if !m.gameOver {
			return m, m.tick()
		}
	}
	return m, nil
}

func (m *Model) turn(direction Direction) {
	if m.paused || m.gameOver {
		return
	}
	if (m.direction == Up && direction == Down) ||
		(m.direction == Down && direction == Up) ||
		(m.direction == Left && direction == Right) ||
		(m.direction == Right && direction == Left) {
		return
	}
	m.nextDirection = direction
}

func (m *Model) advance() {
	m.direction = m.nextDirection
	head := m.snake[0]
	switch m.direction {
	case Up:
		head.y--
	case Down:
		head.y++
	case Left:
		head.x--
	case Right:
		head.x++
	}

	if head.x < 0 || head.x >= boardWidth || head.y < 0 || head.y >= boardHeight {
		m.gameOver = true
		return
	}

	eating := head == m.food
	body := m.snake
	if !eating {
		body = body[:len(body)-1] // The tail moves away on this step.
	}
	for _, segment := range body {
		if head == segment {
			m.gameOver = true
			return
		}
	}

	if eating {
		m.snake = append(m.snake, Position{})
	}
	copy(m.snake[1:], m.snake[:len(m.snake)-1])
	m.snake[0] = head

	if eating {
		m.score++
		if len(m.snake) == boardWidth*boardHeight {
			m.gameOver = true
			m.won = true
		} else {
			m.placeFood()
		}
	}
}

func (m *Model) placeFood() {
	var occupied [boardHeight][boardWidth]bool
	for _, segment := range m.snake {
		occupied[segment.y][segment.x] = true
	}

	free := make([]Position, 0, boardWidth*boardHeight-len(m.snake))
	for y := range boardHeight {
		for x := range boardWidth {
			if !occupied[y][x] {
				free = append(free, Position{x, y})
			}
		}
	}
	m.food = free[rand.IntN(len(free))]
}

func (m *Model) View() tea.View {
	// Ensure themed styles are ready when Snake is launched directly.
	styles.GetStyles()
	if m.width < viewWidth || m.height < viewHeight {
		message := fmt.Sprintf("Terminal too small\nNeed at least %d × %d cells", viewWidth, viewHeight)
		return tea.NewView(styles.CenterContent(styles.GetHelpStyle().Render(message), m.width, m.height))
	}

	title := styles.GetTitleStyle().MarginBottom(0).Render("S N A K E")
	score := styles.GetSelectedItemStyle().Render(fmt.Sprintf("SCORE %03d", m.score))
	gap := viewWidth - lipgloss.Width(title) - lipgloss.Width(score)
	header := title + strings.Repeat(" ", gap) + score

	status := "Find food. Keep moving."
	statusStyle := styles.GetMenuItemStyle()
	switch {
	case m.won:
		status = "BOARD CLEARED · YOU WIN"
		statusStyle = styles.GetSelectedItemStyle()
	case m.gameOver:
		status = "GAME OVER · R TO RETRY"
		statusStyle = styles.GetGameOverStyle()
	case m.paused:
		status = "PAUSED · SPACE TO RESUME"
		statusStyle = styles.GetSelectedItemStyle()
	}

	help := styles.GetHelpStyle().MarginTop(0)
	content := strings.Join([]string{
		header,
		m.renderBoard(),
		"",
		statusStyle.Render(status),
		help.Render("← ↑ ↓ → move · Space pause"),
		help.Render("R restart · Esc back"),
	}, "\n")
	return tea.NewView(styles.CenterContent(content, m.width, m.height))
}

func (m *Model) renderBoard() string {
	var board [boardHeight][boardWidth]byte
	for _, segment := range m.snake {
		board[segment.y][segment.x] = 1
	}
	board[m.snake[0].y][m.snake[0].x] = 2
	if !m.won {
		board[m.food.y][m.food.x] = 3
	}

	border := styles.GetBorderStyle()
	vertical := border.Render("│")
	body := styles.GetSnakeStyle().Render("██")
	head := styles.GetSnakeHeadStyle().Render("▓▓")
	food := styles.GetStyles().FoodStyle().Render("◆◆")

	var out strings.Builder
	out.Grow((viewWidth + 1) * (boardHeight + 2))
	out.WriteString(border.Render("┌" + strings.Repeat("─", boardWidth*cellWidth) + "┐"))
	for y := range boardHeight {
		out.WriteByte('\n')
		out.WriteString(vertical)
		for x := range boardWidth {
			switch board[y][x] {
			case 1:
				out.WriteString(body)
			case 2:
				out.WriteString(head)
			case 3:
				out.WriteString(food)
			default:
				out.WriteString("  ")
			}
		}
		out.WriteString(vertical)
	}
	out.WriteByte('\n')
	out.WriteString(border.Render("└" + strings.Repeat("─", boardWidth*cellWidth) + "┘"))
	return out.String()
}
