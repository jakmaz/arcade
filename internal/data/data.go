package data

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Used for migrating breaking changes
const CurrentDataVersion = 1

type AppData struct {
	Version  int
	Settings Settings
	Stats    map[string]GameStats
}

type Settings struct {
	CurrentTheme string
}

type GameStats struct {
	HighScore   int `yaml:"highScore"`
	GamesPlayed int `yaml:"gamesPlayed"`

	// Game-specific stats (extensible)
	Extra map[string]any `yaml:"extra,omitempty"`
}

// Main API functions
func LoadData() (*AppData, error) {
	dataDir, err := GetDataPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get data path: %w", err)
	}

	data, err := os.ReadFile(filepath.Join(dataDir, "data.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	var appData AppData
	err = yaml.Unmarshal(data, &appData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &appData, nil
}

func SaveData(data *AppData) error {
	dataDir, err := GetDataPath()
	if err != nil {
		return fmt.Errorf("failed to get data path: %w", err)
	}

	dataToSave, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = os.MkdirAll(dataDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	err = os.WriteFile(filepath.Join(dataDir, "data.yaml"), dataToSave, 0o644)
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	return nil
}

// Convenience functions
func LoadCurrentTheme() string {
	data, err := LoadData()
	if err != nil {
		// File doesn't exist or corrupted - return system theme
		return "system"
	}

	if data.Settings.CurrentTheme == "" {
		return "system"
	}

	return data.Settings.CurrentTheme
}

func SaveCurrentTheme(theme string) error {
	data, err := LoadData()
	if err != nil {
		// File doesn't exist - create default data
		data = &AppData{
			Version:  CurrentDataVersion,
			Settings: Settings{CurrentTheme: theme},
			Stats:    make(map[string]GameStats),
		}
	} else {
		// Update existing data
		data.Settings.CurrentTheme = theme
	}

	return SaveData(data)
}

func GetGameStats(game string) *GameStats {
	data, err := LoadData()
	if err != nil {
		return nil // No data file
	}

	stats, exists := data.Stats[game]
	if !exists {
		return nil // Game not found
	}

	return &stats
}

func UpdateGameStats(game string, stats *GameStats) error {
	data, err := LoadData()
	if err != nil {
		// Create default data if file doesn't exist
		data = &AppData{
			Version:  CurrentDataVersion,
			Settings: Settings{CurrentTheme: "system"},
			Stats:    make(map[string]GameStats),
		}
	}

	data.Stats[game] = *stats
	return SaveData(data)
}
