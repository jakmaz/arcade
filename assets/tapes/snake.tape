Require bun
Require tmux
Source assets/tapes/oxide.tape
Hide
Type "bun assets/tapes/play.ts snake"
Enter
Wait+Screen /PAUSED/
Space
Wait+Screen /Find food/
Show
Screenshot "assets/screenshots/snake.png"
Sleep 100ms
Hide
Ctrl+C
