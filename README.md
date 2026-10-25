**Arcade** is a collection of classic games for your terminal, written in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Cobra CLI](https://github.com/spf13/cobra).
Play Snake, Tetris, 2048, and Tic-Tac-Toe directly from the command line - perfect for quick breaks between coding sessions.

<p align="center">
  <img src="assets/gameplay.gif" width="800" alt="Arcade's menu and Tetris gameplay with line clears">
  <br>
  <sub>System theme with the Oxide terminal palette.</sub>
</p>

## Features

* **Terminal Native**: Runs entirely in your terminal with rich TUI
* **Multiple Games**: Snake, Tetris, 2048, and Tic-Tac-Toe
* **Theming System**: 6 built-in themes with support for system theme detection and custom themes
* **Extensible**: Easy to add new games with consistent UI patterns

## Installation

```bash
go install github.com/jakmaz/arcade@latest
```

## Usage

### Interactive Menu
Launch the interactive game menu:
```bash
arcade
```

### Other Commands
```bash
arcade list                # List all available games
arcade play [game]         # Start a game directly
arcade --help              # View all available commands and options
arcade --version           # Show version information
```

## Games

<table>
  <tr>
    <td align="center" width="50%">
      <strong>Snake</strong><br>
      <a href="assets/screenshots/snake.png"><img src="assets/screenshots/snake.png" width="420" alt="Snake after eating food, with a growing body and score"></a>
    </td>
    <td align="center" width="50%">
      <strong>Tetris</strong><br>
      <a href="assets/screenshots/tetris.png"><img src="assets/screenshots/tetris.png" width="420" alt="Tetris with placed blocks, a landing preview, and a cleared line"></a>
    </td>
  </tr>
  <tr>
    <td align="center" width="50%">
      <strong>2048</strong><br>
      <a href="assets/screenshots/2048.png"><img src="assets/screenshots/2048.png" width="420" alt="2048 during play, with merged numbered tiles and a score"></a>
    </td>
    <td align="center" width="50%">
      <strong>Tic-Tac-Toe</strong><br>
      <a href="assets/screenshots/tictactoe.png"><img src="assets/screenshots/tictactoe.png" width="420" alt="A local Tic-Tac-Toe match with X and O marks and the selected cell"></a>
    </td>
  </tr>
</table>

## Themes

Arcade supports multiple built-in themes with custom theme support:

- **Dracula** - Dark theme with purple accents
- **GitHub** - GitHub's dark color scheme
- **Gruvbox Dark** - Warm, retro color scheme
- **Nord** - Arctic, north-bluish theme
- **Tokyo Night** - Dark theme with neon highlights

### Custom Themes
Create custom themes by adding YAML files to `~/.config/arcade`:
```yaml
name: mytheme
palette:
  bg: "#1a1a1a"
  fg: "#ffffff"
  # ... more colors
ui:
  primary: fg
  accent: purple
  # ... UI mappings
```
You can find the structure of the theme file in [internal/theme/themes/dracula.yaml](internal/theme/themes/dracula.yaml).
If you are happy with your theme, please consider contributing it back to the project!

## Contributing

Arcade welcomes contributions! Whether you want to add new games, new themes, fix bugs, or improve the UI, your help is highly appreciated.
The ultimate goal is to make `arcade` an amazing community-driven game collection for your terminal.

### Adding a New Game

1. Create a new package in `internal/games/yourgame/`
2. Implement the Bubble Tea model interface:
   ```go
   type Model struct { /* your game state */ }
   func (m Model) Init() tea.Cmd { /* initialization */ }
   func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { /* handle input */ }
   func (m Model) View() tea.View { /* render UI */ }
   ```
3. Register your game in `internal/core/games.go`
4. Follow existing UI patterns from other games
5. Use the shared styles from `internal/ui/styles/`

## Acknowledgments

- **Theming System**: Inspired by [OpenCode](https://github.com/sst/opencode)
- **Tetris**: Inspired by [tetrigo](https://github.com/Broderick-Westrope/tetrigo)

## License

This project is licensed under the [MIT License](./LICENSE).
