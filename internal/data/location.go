package data

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetDataPath returns the platform-specific directory for storing arcade data
// Follows XDG Base Directory specification on Linux and platform conventions elsewhere
func GetDataPath() (string, error) {
	// Try XDG_DATA_HOME first (Linux standard)
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		return filepath.Join(dataHome, "arcade"), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Use platform-specific data directories
	switch runtime.GOOS {
	case "darwin":
		// macOS: ~/Library/Application Support/arcade
		return filepath.Join(homeDir, "Library", "Application Support", "arcade"), nil
	case "windows":
		// Windows: %APPDATA%/arcade or ~/.arcade as fallback
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "arcade"), nil
		}
		return filepath.Join(homeDir, ".arcade"), nil
	default:
		// Linux: ~/.local/share/arcade (XDG fallback)
		return filepath.Join(homeDir, ".local", "share", "arcade"), nil
	}
}
