package twentyfortyeight

import (
	"fmt"
	"image/color"
	"math"
	"math/bits"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/jakmaz/arcade/internal/theme"
	"github.com/jakmaz/arcade/internal/ui/styles"
)

const (
	// Six terminal columns are approximately as wide as three rows are tall.
	cellWidth   = 6
	cellHeight  = 3
	boardHeight = boardSize*cellHeight + boardSize + 1
	viewHeight  = boardHeight + 5
)

func (m *Model) tileWidth() int {
	width := cellWidth
	for _, row := range m.board {
		for _, value := range row {
			width = max(width, len(strconv.Itoa(value)))
		}
	}
	return width
}

func (m *Model) viewWidth() int {
	return max(boardSize*m.tileWidth()+boardSize+1, len("2048 SCORE ")+len(strconv.Itoa(m.score)))
}

func (m *Model) View() tea.View {
	// Direct game launches need the same lazy theme initialization as the menu.
	styles.GetStyles()
	if m.width <= 0 || m.height <= 0 {
		return tea.NewView("")
	}
	width := m.viewWidth()
	if m.width < width || m.height < viewHeight {
		message := fmt.Sprintf("Terminal too small\nNeed %d × %d cells", width, viewHeight)
		message = styles.GetHelpStyle().MarginTop(0).Render(message)
		message = lipgloss.NewStyle().MaxWidth(m.width).MaxHeight(m.height).Render(message)
		return tea.NewView(styles.CenterContent(message, m.width, m.height))
	}

	title := styles.GetTitleStyle().MarginBottom(0).Render("2048")
	score := styles.GetSelectedItemStyle().Render(fmt.Sprintf("SCORE %d", m.score))
	header := title + strings.Repeat(" ", width-lipgloss.Width(title)-lipgloss.Width(score)) + score
	status := "Merge tiles. Reach 2048."
	statusStyle := styles.GetMenuItemStyle()
	switch {
	case m.won && !m.continued:
		status = "2048! Enter/Space to continue"
		statusStyle = styles.GetStyles().SuccessStyle()
	case m.gameOver:
		status = "NO MOVES · R TO RETRY"
		statusStyle = styles.GetGameOverStyle()
	case m.won:
		status = "2048 reached. Keep going."
		statusStyle = styles.GetStyles().SuccessStyle()
	}
	help := styles.GetHelpStyle().MarginTop(0)
	content := strings.Join([]string{
		header,
		lipgloss.PlaceHorizontal(width, lipgloss.Center, m.renderBoard()),
		"",
		statusStyle.Render(status),
		help.Render("← ↑ ↓ → move · R restart"),
		help.Render("Esc back"),
	}, "\n")
	return tea.NewView(styles.CenterContent(content, m.width, m.height))
}

func (m *Model) renderBoard() string {
	styles.GetStyles()
	t := theme.CurrentTheme()
	width := m.tileWidth()
	border := styles.GetBorderStyle()
	vertical := border.Render("│")
	horizontal := strings.Repeat("─", width)
	separator := func(start, joint, end string) string {
		return border.Render(start + strings.Repeat(horizontal+joint, boardSize-1) + horizontal + end)
	}

	var out strings.Builder
	out.WriteString(separator("┌", "┬", "┐"))
	for row := range boardSize {
		var tiles [boardSize][]string
		for column, value := range m.board[row] {
			label := ""
			style := lipgloss.NewStyle().Width(width).Height(cellHeight).Align(lipgloss.Center, lipgloss.Center)
			if value != 0 {
				label = strconv.Itoa(value)
				style = style.Inherit(m.tileStyle(t, value))
			}
			tiles[column] = strings.Split(style.Render(label), "\n")
		}
		for line := range cellHeight {
			out.WriteByte('\n')
			out.WriteString(vertical)
			for column := range boardSize {
				out.WriteString(tiles[column][line])
				out.WriteString(vertical)
			}
		}
		out.WriteByte('\n')
		if row == boardSize-1 {
			out.WriteString(separator("└", "┴", "┘"))
		} else {
			out.WriteString(separator("├", "┼", "┤"))
		}
	}
	return out.String()
}

func (m *Model) tileStyle(t theme.Theme, value int) lipgloss.Style {
	palette := [...]color.Color{
		t.Secondary(), t.TetrisI(), t.TetrisJ(), t.TetrisT(),
		t.TetrisS(), t.TetrisO(), t.TetrisL(), t.TetrisZ(),
		t.Accent(), t.Warning(), t.Success(),
	}
	index := min(max(bits.Len(uint(value))-2, 0), len(palette)-1)
	fill := palette[index]
	if basic, ok := fill.(ansi.BasicColor); ok && basic < lipgloss.BrightBlack {
		// Keep terminal palette references, using their brighter variants for fills.
		fill = basic + lipgloss.BrightBlack
	}
	style := lipgloss.NewStyle().Bold(true)
	if _, system := t.Primary().(lipgloss.NoColor); system {
		if m.darkBackground {
			// Reverse uses the terminal's own background as the label color.
			return style.Foreground(fill).Reverse(true)
		}
		return style.Foreground(t.Primary()).Background(fill)
	}

	text := t.Primary()
	if background := t.Background(); background != (lipgloss.NoColor{}) && contrast(background, fill) > contrast(text, fill) {
		text = background
	}
	if contrast(text, fill) < 4.5 {
		// Some theme midtones need a stronger label than their normal UI text.
		for _, candidate := range []color.Color{lipgloss.Darken(t.Primary(), 0.98), lipgloss.Lighten(t.Primary(), 0.98)} {
			if contrast(candidate, fill) > contrast(text, fill) {
				text = candidate
			}
		}
	}
	return style.Foreground(text).Background(fill)
}

func contrast(a, b color.Color) float64 {
	luminance := func(c color.Color) float64 {
		r, g, b, _ := c.RGBA()
		linear := func(component uint32) float64 {
			v := float64(component) / 65535
			if v <= 0.04045 {
				return v / 12.92
			}
			return math.Pow((v+0.055)/1.055, 2.4)
		}
		return 0.2126*linear(r) + 0.7152*linear(g) + 0.0722*linear(b)
	}
	x, y := luminance(a), luminance(b)
	return (max(x, y) + 0.05) / (min(x, y) + 0.05)
}
