package theme

import (
	"fmt"
	"path/filepath"
	"sort"
	"sync"
)

// Manager handles theme registration and switching
type Manager struct {
	mu           sync.RWMutex
	themes       map[string]Theme
	currentTheme Theme
	defaultTheme Theme
}

var globalManager = &Manager{
	themes: make(map[string]Theme),
}

// GetManager returns the global theme manager
func GetManager() *Manager {
	return globalManager
}

// RegisterTheme adds a theme to the registry
func (m *Manager) RegisterTheme(theme Theme) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.themes[theme.Name()] = theme

	// Set as default if it's the first theme registered
	if m.defaultTheme == nil {
		m.defaultTheme = theme
	}

	// Set as current if no current theme is set
	if m.currentTheme == nil {
		m.currentTheme = theme
	}
}

// SetCurrentTheme sets the active theme by name
func (m *Manager) SetCurrentTheme(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	theme, exists := m.themes[name]
	if !exists {
		return fmt.Errorf("theme '%s' not found", name)
	}

	m.currentTheme = theme
	return nil
}

// GetCurrentTheme returns the currently active theme
func (m *Manager) GetCurrentTheme() Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.currentTheme != nil {
		return m.currentTheme
	}

	return m.defaultTheme
}

// GetTheme returns a specific theme by name
func (m *Manager) GetTheme(name string) (Theme, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	theme, exists := m.themes[name]
	return theme, exists
}

// ListThemes returns all registered theme names
func (m *Manager) ListThemes() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.themes))
	for name := range m.themes {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// LoadThemesFromDirectories loads themes from multiple directories
func (m *Manager) LoadThemesFromDirectories(dirs ...string) error {
	for _, dir := range dirs {
		themes, err := LoadThemesFromDirectory(dir)
		if err != nil {
			return fmt.Errorf("failed to load themes from %s: %w", dir, err)
		}

		for _, theme := range themes {
			m.RegisterTheme(theme)
		}
	}

	return nil
}

// Initialize loads themes from standard locations
func (m *Manager) Initialize() error {
	// Register built-in themes first
	m.RegisterTheme(NewDefaultTheme())
	m.RegisterTheme(NewSystemTheme())

	// Get executable directory for built-in themes
	execPath, err := getCurrentExecutableDir()
	if err == nil {
		builtinDir := filepath.Join(execPath, "themes")
		m.LoadThemesFromDirectories(builtinDir)
	}

	// Load built-in themes from the source code location (for development)
	m.LoadThemesFromDirectories("internal/theme/themes")

	// If no themes were loaded, we already have default theme registered
	return nil
}

// getCurrentExecutableDir returns the directory containing the current executable
func getCurrentExecutableDir() (string, error) {
	// For now, just return an error since we don't need this for development
	// This would be implemented for production builds
	return "", fmt.Errorf("executable directory detection not implemented")
}

// Convenience functions for the global manager

// RegisterTheme registers a theme with the global manager
func RegisterTheme(theme Theme) {
	globalManager.RegisterTheme(theme)
}

// SetCurrentTheme sets the active theme in the global manager
func SetCurrentTheme(name string) error {
	return globalManager.SetCurrentTheme(name)
}

// GetCurrentTheme returns the currently active theme from the global manager
func GetCurrentTheme() Theme {
	return globalManager.GetCurrentTheme()
}

// GetTheme returns a specific theme by name from the global manager
func GetTheme(name string) (Theme, bool) {
	return globalManager.GetTheme(name)
}

// ListThemes returns all registered theme names from the global manager
func ListThemes() []string {
	return globalManager.ListThemes()
}

// Initialize loads themes into the global manager
func Initialize() error {
	return globalManager.Initialize()
}
