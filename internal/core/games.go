package core

import (
	"sort"

	tea "charm.land/bubbletea/v2"
	"github.com/jakmaz/arcade/internal/games/snake"
	"github.com/jakmaz/arcade/internal/games/tetris"
	"github.com/jakmaz/arcade/internal/games/tictactoe"
	"github.com/jakmaz/arcade/internal/games/twentyfortyeight"
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
		Description: "Eat, grow, avoid the walls",
		NewModel:    func() tea.Model { return snake.New() },
	},
	"tetris": {
		ID:          "tetris",
		Name:        "Tetris",
		Description: "Stack blocks and clear lines",
		NewModel:    func() tea.Model { return tetris.New() },
	},
	"2048": {
		ID:          "2048",
		Name:        "2048",
		Description: "Slide, merge, reach 2048",
		NewModel:    func() tea.Model { return twentyfortyeight.New() },
	},
	"tictactoe": {
		ID:          "tictactoe",
		Name:        "Tic-Tac-Toe",
		Description: "Get three in a row",
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
