package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme cache to avoid repeated filesystem scans
var (
	customThemesCache map[string]Theme
	cacheInitialized  bool
	themesDir         string = "themes" // Default, can be overridden with SetThemesDir
)

type Theme struct {
	Name string

	// UI Colors
	SelectedBg       lipgloss.Color
	SelectedFg       lipgloss.Color
	NormalFg         lipgloss.Color
	SeparatorColor   lipgloss.Color
	HelpColor        lipgloss.Color // Help bar text at bottom of screen
	TestCountColor   lipgloss.Color // Test count "(N tests)" in package list
	BorderColor      lipgloss.Color
	TreeSymbolColor  lipgloss.Color
	CoverageGoodFg   lipgloss.Color
	CoverageMediumFg lipgloss.Color
	CoveragePoorFg   lipgloss.Color

	// Menu Colors
	MenuNormalFg     lipgloss.Color
	MenuActiveFg     lipgloss.Color
	MenuSelectedBg   lipgloss.Color
	MenuSelectedFg   lipgloss.Color
}

// Built-in themes (hardcoded fallback if theme files aren't available)
// Only gapistotle theme is hardcoded - all other themes are loaded from .conf files
var themes = map[string]Theme{
	"gapistotle": {
		Name:             "gapistotle",
		SelectedBg:       lipgloss.Color("#ffffff"),
		SelectedFg:       lipgloss.Color("#f53900"),
		NormalFg:         lipgloss.Color("#ffffb6"),
		SeparatorColor:   lipgloss.Color("#bd2c00"),
		HelpColor:        lipgloss.Color("#f53900"),
		TestCountColor:   lipgloss.Color("#4c90a4"),
		BorderColor:      lipgloss.Color("#bd2c00"),
		TreeSymbolColor:  lipgloss.Color("#8f6dc0"),
		CoverageGoodFg:   lipgloss.Color("#a3be8c"),
		CoverageMediumFg: lipgloss.Color("#ebcb8b"),
		CoveragePoorFg:   lipgloss.Color("#bf616a"),
		MenuNormalFg:     lipgloss.Color("#f53900"),
		MenuActiveFg:     lipgloss.Color("#d8dee9"),
		MenuSelectedBg:   lipgloss.Color("#7e5f00"),
		MenuSelectedFg:   lipgloss.Color("#ffffff"),
	},
}

// GetTheme returns a theme by name, or gapistotle if not found
func GetTheme(name string) Theme {
	allThemes := GetAllThemes()
	if theme, ok := allThemes[name]; ok {
		return theme
	}
	return themes["gapistotle"]
}

// ListThemes returns all available theme names (built-in + custom)
func ListThemes() []string {
	allThemes := GetAllThemes()
	names := make([]string, 0, len(allThemes))
	for name := range allThemes {
		names = append(names, name)
	}
	return names
}

// GetAllThemes returns both built-in and custom themes
// Uses a cache to avoid repeated filesystem scans
func GetAllThemes() map[string]Theme {
	allThemes := make(map[string]Theme)

	// Start with built-in themes
	for name, theme := range themes {
		allThemes[name] = theme
	}

	// Add custom themes using cache
	if !cacheInitialized {
		customThemesCache = LoadCustomThemes()
		cacheInitialized = true
	}

	for name, theme := range customThemesCache {
		allThemes[name] = theme
	}

	return allThemes
}

// InvalidateThemeCache forces a reload of custom themes on the next GetAllThemes call
func InvalidateThemeCache() {
	cacheInitialized = false
	customThemesCache = nil
}

// SetThemesDir sets the directory for custom themes
// Should be called at startup after loading config
func SetThemesDir(dir string) {
	themesDir = dir
	// Invalidate cache so themes will be reloaded from new directory
	InvalidateThemeCache()
}
