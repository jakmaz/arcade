package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jakmaz/arcade/internal/core"
	"github.com/spf13/cobra"
)

// GameWrapper wraps a game model to handle universal exit controls for direct launches
type GameWrapper struct {
	game tea.Model
}

// NewGameWrapper creates a new wrapper around a game
func NewGameWrapper(game tea.Model) *GameWrapper {
	return &GameWrapper{
		game: game,
	}
}

func (gw *GameWrapper) Init() tea.Cmd {
	return gw.game.Init()
}

func (gw *GameWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle universal exit controls
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return gw, tea.Quit
		}
	}

	// Delegate to the wrapped game
	var cmd tea.Cmd
	gw.game, cmd = gw.game.Update(msg)
	return gw, cmd
}

func (gw *GameWrapper) View() string {
	return gw.game.View()
}

func init() {
	rootCmd.AddCommand(playCmd)
}

var playCmd = &cobra.Command{
	Use:   "play <game>",
	Short: "Play a specific game",
	Long:  "Launch a game directly without going through the menu",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		// Validate game ID
		if _, exists := core.Games[gameID]; !exists {
			fmt.Printf("Game %s does not exist\n", gameID)
			os.Exit(1)
		}
		playGame(gameID)
	},
}

func playGame(gameID string) {
	game := core.CreateGame(gameID)
	if game == nil {
		fmt.Printf("Game '%s' not found\n", gameID)
		return
	}

	// Wrap the game to handle exit controls
	wrappedGame := NewGameWrapper(game)

	p := tea.NewProgram(wrappedGame, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
