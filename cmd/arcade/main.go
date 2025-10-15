package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jakmaz/arcade/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.NewMenu())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
