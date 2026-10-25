# The capture player reads the live board and uses normal arrow/drop controls.
Require bun
Require tmux
Output "assets/gameplay.gif"
Set PlaybackSpeed 1.5
Source assets/tapes/oxide.tape
Hide
Type "bun assets/tapes/play.ts tetris"
Enter
Wait+Screen /Theme:/
Show
Wait+Screen /T E T R I S/
Wait+Screen /LINES [1-9]/
Sleep 1500ms
Screenshot "assets/screenshots/tetris.png"
Sleep 1s
Hide
Ctrl+C
