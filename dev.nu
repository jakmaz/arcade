#!/usr/bin/env nu

# Hot-reload development script for Arcade
# Starts a new tmux session and watches for changes, automatically restarting the game
# Usage: nu dev.nu [game]  # game: snake, chess, tetris, tictactoe

def main [game?: string] {
    let session_name = "arcade-dev"

    # Kill existing session if it exists
    try { ^tmux kill-session -t $session_name }

    # Determine what command to run
    let run_cmd = if ($game == null) {
        "go run ."
    } else {
        $"go run . play ($game)"
    }

    let mode = if ($game == null) { "main menu" } else { $"($game)" }

    print $"Hot-reload: ($mode)"

    # Create new tmux session
    ^tmux new-session -d -s $session_name
    ^tmux send-keys -t $session_name $run_cmd Enter

    print "Use 'tmux attach -t arcade-dev' to view"

    # Watch for changes
    watch . --glob="**/*.go" { |op, path, new_path|
        let timestamp = (date now | format date "%H:%M:%S")
        print $"[($timestamp)] ($op): ($path)"

        # Kill current process and restart
        ^tmux send-keys -t $session_name C-c
        sleep 200ms
        ^tmux send-keys -t $session_name $run_cmd Enter
    }
}
