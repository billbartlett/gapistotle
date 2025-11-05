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

var themes = map[string]Theme{
	"gapistotle": {
		Name:             "gapistotle",
		SelectedBg:       lipgloss.Color("#003746"),
		SelectedFg:       lipgloss.Color("#b2b2b2"),
		NormalFg:         lipgloss.Color("#84804a"),
		SeparatorColor:   lipgloss.Color("#f5ed8b"),
		HelpColor:        lipgloss.Color("#4c90a4"),
		TestCountColor:   lipgloss.Color("#4c90a4"),
		BorderColor:      lipgloss.Color("#f5ed8b"),
		TreeSymbolColor:  lipgloss.Color("#8f6dc0"),
		CoverageGoodFg:   lipgloss.Color("#a3be8c"),
		CoverageMediumFg: lipgloss.Color("#ebcb8b"),
		CoveragePoorFg:   lipgloss.Color("#bf616a"),
		MenuNormalFg:     lipgloss.Color("#f53900"),
		MenuActiveFg:     lipgloss.Color("#d8dee9"),
		MenuSelectedBg:   lipgloss.Color("#7e5f00"),
		MenuSelectedFg:   lipgloss.Color("#ffffff"),
	},
	"dracula": {
		Name:             "dracula",
		SelectedBg:       lipgloss.Color("#bd93f9"), // Purple
		SelectedFg:       lipgloss.Color("#f8f8f2"), // Foreground
		NormalFg:         lipgloss.Color("#f8f8f2"), // Foreground
		SeparatorColor:   lipgloss.Color("#44475a"), // Current line
		HelpColor:        lipgloss.Color("#6272a4"), // Comment
		TestCountColor:   lipgloss.Color("#6272a4"), // Comment
		BorderColor:      lipgloss.Color("#6272a4"), // Comment
		TreeSymbolColor:  lipgloss.Color("#6272a4"), // Comment
		CoverageGoodFg:   lipgloss.Color("#50fa7b"), // Green
		CoverageMediumFg: lipgloss.Color("#f1fa8c"), // Yellow
		CoveragePoorFg:   lipgloss.Color("#ff5555"), // Red
		MenuNormalFg:     lipgloss.Color("#6272a4"), // Comment
		MenuActiveFg:     lipgloss.Color("#f8f8f2"), // Foreground
		MenuSelectedBg:   lipgloss.Color("#bd93f9"), // Purple
		MenuSelectedFg:   lipgloss.Color("#f8f8f2"), // Foreground
	},
	"nord": {
		Name:             "nord",
		SelectedBg:       lipgloss.Color("#88c0d0"), // Frost
		SelectedFg:       lipgloss.Color("#2e3440"), // Polar night (dark)
		NormalFg:         lipgloss.Color("#d8dee9"), // Snow storm
		SeparatorColor:   lipgloss.Color("#4c566a"), // Polar night
		HelpColor:        lipgloss.Color("#4c566a"), // Polar night
		TestCountColor:   lipgloss.Color("#4c566a"), // Polar night
		BorderColor:      lipgloss.Color("#4c566a"), // Polar night
		TreeSymbolColor:  lipgloss.Color("#4c566a"), // Polar night
		CoverageGoodFg:   lipgloss.Color("#a3be8c"), // Green
		CoverageMediumFg: lipgloss.Color("#ebcb8b"), // Yellow
		CoveragePoorFg:   lipgloss.Color("#bf616a"), // Red
		MenuNormalFg:     lipgloss.Color("#4c566a"), // Polar night
		MenuActiveFg:     lipgloss.Color("#d8dee9"), // Snow storm
		MenuSelectedBg:   lipgloss.Color("#88c0d0"), // Frost
		MenuSelectedFg:   lipgloss.Color("#2e3440"), // Polar night
	},
	"monokai": {
		Name:             "monokai",
		SelectedBg:       lipgloss.Color("#66d9ef"), // Cyan
		SelectedFg:       lipgloss.Color("#272822"), // Background
		NormalFg:         lipgloss.Color("#f8f8f2"), // Foreground
		SeparatorColor:   lipgloss.Color("#75715e"), // Comment
		HelpColor:        lipgloss.Color("#75715e"), // Comment
		TestCountColor:   lipgloss.Color("#75715e"), // Comment
		BorderColor:      lipgloss.Color("#75715e"), // Comment
		TreeSymbolColor:  lipgloss.Color("#75715e"), // Comment
		CoverageGoodFg:   lipgloss.Color("#a6e22e"), // Green
		CoverageMediumFg: lipgloss.Color("#e6db74"), // Yellow
		CoveragePoorFg:   lipgloss.Color("#f92672"), // Pink/Red
		MenuNormalFg:     lipgloss.Color("#75715e"), // Comment
		MenuActiveFg:     lipgloss.Color("#f8f8f2"), // Foreground
		MenuSelectedBg:   lipgloss.Color("#66d9ef"), // Cyan
		MenuSelectedFg:   lipgloss.Color("#272822"), // Background
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
