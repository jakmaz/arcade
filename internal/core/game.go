package core

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jakmaz/arcade/internal/games/chess"
	"github.com/jakmaz/arcade/internal/games/snake"
	"github.com/jakmaz/arcade/internal/games/tetris"
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
}

func AvailableGames() []GameInfo {
	var games []GameInfo
	for _, game := range Games {
		games = append(games, game)
	}
	return games
}

func CreateGame(id string) tea.Model {
	if game, exists := Games[id]; exists {
		return game.NewModel()
	}
	return nil
}
