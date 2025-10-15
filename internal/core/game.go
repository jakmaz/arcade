package core

type Game struct {
	Name        string
	Description string
}

func AvailableGames() []Game {
	return []Game{
		{
			Name:        "Snake",
			Description: "A classic game of snake",
		},
		{
			Name:        "Tetris",
			Description: "A classic game of tetris",
		},
		{
			Name:        "Chess",
			Description: "A classic game of chess",
		},
	}
}
