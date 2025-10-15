# Arcade

**Arcade** is a collection of classic games for your terminal, written in Go.
Play Snake, Tetris, and Chess directly from the command line - perfect for quick breaks between coding sessions.

## Features

* Runs entirely in your terminal
* Works smoothly with your favorite color schemes
* Easy to extend with your own games

## Installation

```bash
go install github.com/jakmaz/arcade/cmd/arcade@latest
```

## Usage

Launch the menu:

```bash
arcade
```

Play a specific game:

```bash
arcade play snake
```

## Contributing

Arcade is open to contributions - new games, bug fixes, or creative hacks are all welcome.
If you’d like to add a new game, implement the core.Game interface and open a pull request.
