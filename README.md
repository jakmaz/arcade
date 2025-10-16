# Arcade

**Arcade** is a collection of classic games for your terminal, written in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Cobra CLI](https://github.com/spf13/cobra).
Play Snake, Tetris, Chess, and Tic-Tac-Toe directly from the command line - perfect for quick breaks between coding sessions.

## Features

* **Multiple Games**: Snake, Tetris, Chess, and Tic-Tac-Toe
* **Terminal Native**: Runs entirely in your terminal with rich TUI
* **Adaptive Colors**: Works smoothly with your favorite color schemes
* **Direct Launch**: Jump straight into games or browse the interactive menu
* **Extensible**: Easy to add new games with consistent UI patterns

## Installation

```bash
go install github.com/jakmaz/arcade/cmd/arcade@latest
```

## Usage

### Interactive Menu
Launch the interactive game menu:
```bash
arcade
```

### Direct Game Launch
Jump directly into any game:
```bash
arcade snake
arcade tetris
arcade chess
arcade tictactoe
```

### Get Help
View all available commands and options:
```bash
arcade --help 
arcade -h
```

### List Games
See all available games:
```bash
arcade list
```

### Game Information
Get detailed info about a specific game:
```bash
arcade info snake
```

### Other Commands
```bash
arcade --version 
arcade random
```

## Available Games

| Game | Command | Description |
|------|---------|-------------|
| **Snake** | `arcade snake` | Classic snake game - eat food, grow longer, avoid walls |
| **Tetris** | `arcade tetris` | Block puzzle game with falling pieces |
| **Chess** | `arcade chess` | Strategic board game with full piece set |
| **Tic-Tac-Toe** | `arcade tictactoe` | Classic X's and O's game |

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal User Interface framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling and layout
- [Cobra](https://github.com/spf13/cobra) - CLI framework

## Contributing

Arcade welcomes contributions! Whether you want to add new games, fix bugs, or improve the UI, your help is appreciated.

### Adding a New Game

1. Create a new package in `internal/games/yourgame/`
2. Implement the Bubble Tea model interface:
   ```go
   type Model struct { /* your game state */ }
   func (m Model) Init() tea.Cmd { /* initialization */ }
   func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { /* handle input */ }
   func (m Model) View() string { /* render UI */ }
   ```
3. Register your game in `internal/core/games.go`
4. Follow existing UI patterns from other games
5. Use the shared styles from `internal/ui/styles/`

### Development

```bash
git clone https://github.com/jakmaz/arcade
cd arcade
go run cmd/arcade/main.go
```
