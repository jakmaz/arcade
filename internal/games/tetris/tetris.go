package tetris

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jakmaz/arcade/internal/theme"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

const (
	boardWidth   = 10
	boardHeight  = 20
	cellWidth    = 2
	sidebarWidth = 23
	viewWidth    = boardWidth*cellWidth + 2 + 3 + sidebarWidth
	viewHeight   = boardHeight + 2
)

type pieceKind uint8

const (
	empty pieceKind = iota
	pieceI
	pieceO
	pieceT
	pieceS
	pieceZ
	pieceJ
	pieceL
)

var pieceNames = [...]string{"", "I", "O", "T", "S", "Z", "J", "L"}

var shapes = [...][4]Position{
	empty:  {},
	pieceI: {{0, 1}, {1, 1}, {2, 1}, {3, 1}},
	pieceO: {{1, 0}, {2, 0}, {1, 1}, {2, 1}},
	pieceT: {{1, 0}, {0, 1}, {1, 1}, {2, 1}},
	pieceS: {{1, 0}, {2, 0}, {0, 1}, {1, 1}},
	pieceZ: {{0, 0}, {1, 0}, {1, 1}, {2, 1}},
	pieceJ: {{0, 0}, {0, 1}, {1, 1}, {2, 1}},
	pieceL: {{2, 0}, {0, 1}, {1, 1}, {2, 1}},
}

type Position struct {
	x, y int
}

type Piece struct {
	kind     pieceKind
	rotation int
	x, y     int
}

func newPiece(kind pieceKind) Piece {
	return Piece{kind: kind, x: 3}
}

func (p Piece) cells() [4]Position {
	cells := shapes[p.kind]
	for i, cell := range cells {
		if p.kind != pieceO {
			pivot := 2
			if p.kind == pieceI {
				pivot = 3
			}
			for range p.rotation {
				cell.x, cell.y = pivot-cell.y, cell.x
			}
		}
		cells[i] = Position{p.x + cell.x, p.y + cell.y}
	}
	return cells
}

type Model struct {
	board         [boardHeight][boardWidth]pieceKind
	currentPiece  Piece
	nextPiece     pieceKind
	bag           [7]pieceKind
	bagIndex      int
	score         int
	lines         int
	level         int
	paused        bool
	gameOver      bool
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
	m.board = [boardHeight][boardWidth]pieceKind{}
	m.bagIndex = len(m.bag)
	m.currentPiece = newPiece(m.drawPiece())
	m.nextPiece = m.drawPiece()
	m.score, m.lines, m.level = 0, 0, 1
	m.paused, m.gameOver = false, false
	m.tickID++
}

func (m *Model) drawPiece() pieceKind {
	if m.bagIndex == len(m.bag) {
		m.bag = [7]pieceKind{pieceI, pieceO, pieceT, pieceS, pieceZ, pieceJ, pieceL}
		rand.Shuffle(len(m.bag), func(i, j int) { m.bag[i], m.bag[j] = m.bag[j], m.bag[i] })
		m.bagIndex = 0
	}
	piece := m.bag[m.bagIndex]
	m.bagIndex++
	return piece
}

func (m *Model) playing() bool {
	return !m.paused && !m.gameOver && m.width >= viewWidth && m.height >= viewHeight
}

func (m *Model) gravityInterval() time.Duration {
	return max(80*time.Millisecond, 800*time.Millisecond-time.Duration(m.level-1)*60*time.Millisecond)
}

// Each one-shot owns a generation. Suspending or replacing it invalidates any
// queued message, including messages left over from another game instance.
func (m *Model) tick() tea.Cmd {
	m.tickID++
	if !m.playing() {
		return nil
	}
	id := m.tickID
	return tea.Tick(m.gravityInterval(), func(time.Time) tea.Msg {
		return tickMsg{game: m, id: id}
	})
}

func (m *Model) Init() tea.Cmd {
	return m.tick()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		wasPlaying := m.playing()
		m.width, m.height = max(0, msg.Width), max(0, msg.Height)
		if wasPlaying != m.playing() {
			return m, m.tick()
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "r", "R":
			m.reset()
			return m, m.tick()
		case "p", "P":
			if !m.gameOver {
				m.paused = !m.paused
				return m, m.tick()
			}
		}
		if !m.playing() {
			return m, nil
		}
		switch msg.String() {
		case "left":
			m.move(-1, 0)
		case "right":
			m.move(1, 0)
		case "up":
			m.rotate()
		case "down":
			if m.move(0, 1) {
				m.score++
			} else {
				m.lockPiece()
				return m, m.tick()
			}
		case "space":
			landing := m.landingPiece()
			m.score += (landing.y - m.currentPiece.y) * 2
			m.currentPiece = landing
			m.lockPiece()
			return m, m.tick()
		}
	case tickMsg:
		if msg.game != m || msg.id != m.tickID || !m.playing() {
			return m, nil
		}
		if !m.move(0, 1) {
			m.lockPiece()
		}
		return m, m.tick()
	}
	return m, nil
}

func (m *Model) fits(piece Piece) bool {
	for _, cell := range piece.cells() {
		if cell.x < 0 || cell.x >= boardWidth || cell.y >= boardHeight {
			return false
		}
		if cell.y >= 0 && m.board[cell.y][cell.x] != empty {
			return false
		}
	}
	return true
}

func (m *Model) move(dx, dy int) bool {
	piece := m.currentPiece
	piece.x += dx
	piece.y += dy
	if !m.fits(piece) {
		return false
	}
	m.currentPiece = piece
	return true
}

func (m *Model) rotate() bool {
	if m.currentPiece.kind == pieceO {
		return false
	}
	piece := m.currentPiece
	piece.rotation = (piece.rotation + 1) % 4
	// Short sideways kicks, then upward floor kicks; intentionally simpler than SRS.
	for _, kick := range [...]Position{{0, 0}, {-1, 0}, {1, 0}, {-2, 0}, {2, 0}, {0, -1}, {-1, -1}, {1, -1}, {0, -2}} {
		candidate := piece
		candidate.x += kick.x
		candidate.y += kick.y
		if m.fits(candidate) {
			m.currentPiece = candidate
			return true
		}
	}
	return false
}

func (m *Model) landingPiece() Piece {
	piece := m.currentPiece
	for {
		next := piece
		next.y++
		if !m.fits(next) {
			return piece
		}
		piece = next
	}
}

func (m *Model) lockPiece() {
	cells := m.currentPiece.cells()
	for _, cell := range cells {
		if cell.y < 0 {
			m.gameOver = true
			return
		}
	}
	for _, cell := range cells {
		m.board[cell.y][cell.x] = m.currentPiece.kind
	}
	cleared := m.clearRows()
	m.score += [...]int{0, 100, 300, 500, 800}[cleared] * m.level
	m.lines += cleared
	m.level = 1 + m.lines/10
	m.currentPiece = newPiece(m.nextPiece)
	m.nextPiece = m.drawPiece()
	m.gameOver = !m.fits(m.currentPiece)
}

func (m *Model) clearRows() int {
	write := boardHeight - 1
	for read := boardHeight - 1; read >= 0; read-- {
		full := true
		for _, cell := range m.board[read] {
			if cell == empty {
				full = false
				break
			}
		}
		if !full {
			m.board[write] = m.board[read]
			write--
		}
	}
	cleared := write + 1
	for ; write >= 0; write-- {
		m.board[write] = [boardWidth]pieceKind{}
	}
	return cleared
}

func (m *Model) View() tea.View {
	styles.GetStyles() // Direct launches also need the theme initialized.
	if m.width <= 0 || m.height <= 0 {
		return tea.NewView("")
	}
	if m.width < viewWidth || m.height < viewHeight {
		message := fmt.Sprintf("Terminal too small\nNeed %d × %d cells\nResize to resume · Esc back", viewWidth, viewHeight)
		content := styles.GetHelpStyle().MarginTop(0).
			Width(m.width).MaxWidth(m.width).MaxHeight(m.height).Render(message)
		return tea.NewView(styles.CenterContent(content, m.width, m.height))
	}
	content := lipgloss.JoinHorizontal(lipgloss.Top, m.renderBoard(), "   ", m.renderSidebar())
	return tea.NewView(styles.CenterContent(content, m.width, m.height))
}

func (m *Model) renderBoard() string {
	border := styles.GetBorderStyle()
	themed := theme.NewStylesWithTheme(theme.CurrentTheme())
	var blocks [8]string
	for kind := pieceI; kind <= pieceL; kind++ {
		blocks[kind] = themed.TetrisPieceStyle(pieceNames[kind]).Render("██")
	}
	ghost := themed.TetrisPieceStyle(pieceNames[m.currentPiece.kind]).Bold(false).Render("░░")
	board := m.board
	var ghostCells [boardHeight][boardWidth]bool
	if !m.gameOver {
		for _, cell := range m.landingPiece().cells() {
			if cell.y >= 0 && cell.y < boardHeight && cell.x >= 0 && cell.x < boardWidth {
				ghostCells[cell.y][cell.x] = true
			}
		}
	}
	for _, cell := range m.currentPiece.cells() {
		if cell.y >= 0 && cell.y < boardHeight && cell.x >= 0 && cell.x < boardWidth && board[cell.y][cell.x] == empty {
			board[cell.y][cell.x] = m.currentPiece.kind
		}
	}
	vertical := border.Render("│")
	var out strings.Builder
	out.WriteString(border.Render("┌" + strings.Repeat("─", boardWidth*cellWidth) + "┐"))
	for y := range boardHeight {
		out.WriteByte('\n')
		out.WriteString(vertical)
		for x := range boardWidth {
			switch {
			case board[y][x] != empty:
				out.WriteString(blocks[board[y][x]])
			case ghostCells[y][x]:
				out.WriteString(ghost)
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

func (m *Model) renderSidebar() string {
	status, hint := "PLAYING", "░░ landing"
	statusStyle := styles.GetMenuItemStyle()
	switch {
	case m.gameOver:
		status, hint = "GAME OVER", "R to retry"
		statusStyle = styles.GetGameOverStyle()
	case m.paused:
		status, hint = "PAUSED", "P to resume"
		statusStyle = styles.GetSelectedItemStyle()
	}
	help := styles.GetHelpStyle().MarginTop(0)
	content := strings.Join([]string{
		styles.GetTitleStyle().MarginBottom(0).Render("T E T R I S"),
		"",
		styles.GetSelectedItemStyle().Render(fmt.Sprintf("SCORE %06d", m.score)),
		help.Render(fmt.Sprintf("LEVEL %d · LINES %d", m.level, m.lines)),
		"",
		styles.GetSelectedItemStyle().Render("NEXT"),
		m.renderNextPiece(),
		"",
		statusStyle.Render(status),
		help.Render(hint),
		"",
		help.Render("← → move"),
		help.Render("↑ rotate · ↓ soft drop"),
		help.Render("Space hard drop"),
		"",
		help.Render("P pause · R restart"),
		help.Render("Esc back"),
	}, "\n")
	return lipgloss.NewStyle().Width(sidebarWidth).MaxWidth(sidebarWidth).MaxHeight(viewHeight).Render(content)
}

func (m *Model) renderNextPiece() string {
	styles.GetStyles()
	block := theme.NewStylesWithTheme(theme.CurrentTheme()).TetrisPieceStyle(pieceNames[m.nextPiece]).Render("██")
	cells := newPiece(m.nextPiece).cells()
	minX, maxX, minY, maxY := boardWidth, 0, boardHeight, 0
	for _, cell := range cells {
		minX, maxX = min(minX, cell.x), max(maxX, cell.x)
		minY, maxY = min(minY, cell.y), max(maxY, cell.y)
	}
	var rows []string
	for y := minY; y <= maxY; y++ {
		var row strings.Builder
		for x := minX; x <= maxX; x++ {
			filled := false
			for _, cell := range cells {
				filled = filled || cell == (Position{x, y})
			}
			if filled {
				row.WriteString(block)
			} else {
				row.WriteString("  ")
			}
		}
		rows = append(rows, row.String())
	}
	if len(rows) == 1 {
		rows = append(rows, "")
	}
	return strings.Join(rows, "\n")
}
