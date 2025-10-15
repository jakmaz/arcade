package core

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jakmaz/arcade/internal/games/chess"
	"github.com/jakmaz/arcade/internal/games/snake"
	"github.com/jakmaz/arcade/internal/games/tetris"
	"github.com/jakmaz/arcade/internal/games/tictactoe"
)

type GameInfo struct {
	ID          string
	Name        string
	Description string
	NewModel    func() tea.Model
}

// Simple, explicit game registry
var Games = map[string]GameInfo{
	"snake": {
		ID:          "snake",
		Name:        "Snake",
		Description: "Classic Snake game",
		NewModel:    func() tea.Model { return snake.New() },
	},
	"tetris": {
		ID:          "tetris",
		Name:        "Tetris",
		Description: "Block puzzle game",
		NewModel:    func() tea.Model { return tetris.New() },
	},
	"chess": {
		ID:          "chess",
		Name:        "Chess",
		Description: "Strategic board game",
		NewModel:    func() tea.Model { return chess.New() },
	},
	"tictactoe": {
		ID:          "tictactoe",
		Name:        "Tic-Tac-Toe",
		Description: "We all know it",
		NewModel:    func() tea.Model { return tictactoe.New() },
	},
}

func AvailableGames() []GameInfo {
	var games []GameInfo
	for _, game := range Games {
		games = append(games, game)
	}
	sort.Slice(games, func(i, j int) bool {
		return games[i].ID < games[j].ID
	})
	return games
}

func CreateGame(id string) tea.Model {
	if game, exists := Games[id]; exists {
		return game.NewModel()
	}
	return nil
}
