package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Theme editor mode constants
const (
	themeEditorModeProperty = "property"
	themeEditorModeColor    = "color"
	themeEditorModeSave     = "save"
)

// Color palette - Curated visually distinct colors
var colorPalette = []lipgloss.Color{
	lipgloss.Color("#232323"),
	lipgloss.Color("#333333"),
	lipgloss.Color("#424242"),
	lipgloss.Color("#b2b2b2"),
	lipgloss.Color("#b27070"),
	lipgloss.Color("#ffffff"),
	lipgloss.Color("#ffa0a0"),
	lipgloss.Color("#ffd0d0"),
	lipgloss.Color("#841e00"),
	lipgloss.Color("#f53900"),
	lipgloss.Color("#bd2c00"),
	lipgloss.Color("#7e5f00"),
	lipgloss.Color("#b58900"),
	lipgloss.Color("#eee8d5"),
	lipgloss.Color("#a6a295"),
	lipgloss.Color("#a8a162"),
	lipgloss.Color("#f0e68c"),
	lipgloss.Color("#f5ed8b"),
	lipgloss.Color("#bdb76b"),
	lipgloss.Color("#84804a"),
	lipgloss.Color("#ffffb6"),
	lipgloss.Color("#8bb82e"),
	lipgloss.Color("#6b8e24"),
	lipgloss.Color("#4a6319"),
	lipgloss.Color("#8cff58"),
	lipgloss.Color("#4b8a2f"),
	lipgloss.Color("#6cc644"),
	lipgloss.Color("#89fb98"),
	lipgloss.Color("#5faf6a"),
	lipgloss.Color("#1f8236"),
	lipgloss.Color("#2dba4e"),
	lipgloss.Color("#3af165"),
	lipgloss.Color("#b2ffc5"),
	lipgloss.Color("#8dffff"),
	lipgloss.Color("#839496"),
	lipgloss.Color("#aac0c3"),
	lipgloss.Color("#5b6769"),
	lipgloss.Color("#001e25"),
	lipgloss.Color("#002b36"),
	lipgloss.Color("#003746"),
	lipgloss.Color("#4c90a4"),
	lipgloss.Color("#6dceeb"),
	lipgloss.Color("#2c5486"),
	lipgloss.Color("#539cf9"),
	lipgloss.Color("#4078c0"),
	lipgloss.Color("#6e5494"),
	lipgloss.Color("#8f6dc0"),
	lipgloss.Color("#4d3a67"),
}

// ThemeProperty represents a configurable theme property
type ThemeProperty struct {
	Name        string
	Description string
	FieldName   string // Which field in the theme this represents
	GetValue    func(*Theme) lipgloss.Color
	SetValue    func(*Theme, lipgloss.Color)
}

// ThemeEditorState holds the state for the theme editor
type ThemeEditorState struct {
	editingTheme     Theme
	properties       []ThemeProperty
	selectedProperty int
	selectedColor    int
	mode             string // "property", "color", or "save"
	saveThemeName    string // For text input when saving
	saveError        string // Error message if save fails
}

func newThemeEditorState(currentTheme Theme) ThemeEditorState {
	// Define properties that can be edited with function pointers for get/set
	properties := []ThemeProperty{
		{
			Name:        "Package Name (Normal)",
			Description: "Text color for unselected packages",
			FieldName:   "NormalFg",
			GetValue:    func(t *Theme) lipgloss.Color { return t.NormalFg },
			SetValue:    func(t *Theme, c lipgloss.Color) { t.NormalFg = c },
		},
		{
			Name:        "Package Name (Selected Text)",
			Description: "Text color when package is selected",
			FieldName:   "SelectedFg",
			GetValue:    func(t *Theme) lipgloss.Color { return t.SelectedFg },
			SetValue:    func(t *Theme, c lipgloss.Color) { t.SelectedFg = c },
		},
		{
			Name:        "Package Name (Selected BG)",
			Description: "Background color when package is selected",
			FieldName:   "SelectedBg",
			GetValue:    func(t *Theme) lipgloss.Color { return t.SelectedBg },
			SetValue:    func(t *Theme, c lipgloss.Color) { t.SelectedBg = c },
		},
		{
			Name:        "Test Count",
			Description: "Text color for test count '(N tests)'",
			FieldName:   "TestCountColor",
			GetValue:    func(t *Theme) lipgloss.Color { return t.TestCountColor },
			SetValue:    func(t *Theme, c lipgloss.Color) { t.TestCountColor = c },
		},
		{
			Name:        "Help Bar Text",
			Description: "Text color for help bar at bottom",
			FieldName:   "HelpColor",
			GetValue:    func(t *Theme) lipgloss.Color { return t.HelpColor },
			SetValue:    func(t *Theme, c lipgloss.Color) { t.HelpColor = c },
		},
		{
			Name:        "Tree Symbols",
			Description: "Color for tree symbols (├─ │)",
			FieldName:   "TreeSymbolColor",
			GetValue:    func(t *Theme) lipgloss.Color { return t.TreeSymbolColor },
			SetValue:    func(t *Theme, c lipgloss.Color) { t.TreeSymbolColor = c },
		},
		{
			Name:        "Border/Separator",
			Description: "Color for borders and separators",
			FieldName:   "BorderColor",
			GetValue:    func(t *Theme) lipgloss.Color { return t.BorderColor },
			SetValue: func(t *Theme, c lipgloss.Color) {
				// Update both border and separator to the same color
				t.BorderColor = c
				t.SeparatorColor = c
			},
		},
		{
			Name:        "Menu (Normal)",
			Description: "Menu text when inactive",
			FieldName:   "MenuNormalFg",
			GetValue:    func(t *Theme) lipgloss.Color { return t.MenuNormalFg },
			SetValue:    func(t *Theme, c lipgloss.Color) { t.MenuNormalFg = c },
		},
		{
			Name:        "Menu (Selected Text)",
			Description: "Menu text when selected",
			FieldName:   "MenuSelectedFg",
			GetValue:    func(t *Theme) lipgloss.Color { return t.MenuSelectedFg },
			SetValue:    func(t *Theme, c lipgloss.Color) { t.MenuSelectedFg = c },
		},
		{
			Name:        "Menu (Selected BG)",
			Description: "Menu background when selected",
			FieldName:   "MenuSelectedBg",
			GetValue:    func(t *Theme) lipgloss.Color { return t.MenuSelectedBg },
			SetValue:    func(t *Theme, c lipgloss.Color) { t.MenuSelectedBg = c },
		},
	}

	return ThemeEditorState{
		editingTheme:     currentTheme,
		properties:       properties,
		selectedProperty: 0,
		selectedColor:    0,
		mode:             themeEditorModeProperty,
	}
}

// GetCurrentPropertyColor returns the current color of the selected property
func (te *ThemeEditorState) GetCurrentPropertyColor() lipgloss.Color {
	if te.selectedProperty >= len(te.properties) {
		return lipgloss.Color("#ffffff")
	}

	prop := te.properties[te.selectedProperty]
	if prop.GetValue != nil {
		return prop.GetValue(&te.editingTheme)
	}

	return lipgloss.Color("#ffffff")
}

// SetCurrentPropertyColor sets the color for the selected property
func (te *ThemeEditorState) SetCurrentPropertyColor(color lipgloss.Color) {
	if te.selectedProperty >= len(te.properties) {
		return
	}

	prop := te.properties[te.selectedProperty]
	if prop.SetValue != nil {
		prop.SetValue(&te.editingTheme, color)
	}
}

// RenderColorPalette renders the color palette with selection
func (te *ThemeEditorState) RenderColorPalette() string {
	var result string
	result += lipgloss.NewStyle().Bold(true).Render("Color Palette") + "\n\n"

	totalColors := len(colorPalette)
	colsPerRow := 8
	rows := (totalColors + colsPerRow - 1) / colsPerRow // Ceiling division

	// Show all colors in rows of 8
	for row := 0; row < rows; row++ {
		// First render the selection indicator row if in color mode
		if te.mode == themeEditorModeColor {
			indicatorLine := ""
			for col := 0; col < colsPerRow; col++ {
				idx := row*colsPerRow + col
				if idx >= totalColors {
					break
				}
				if idx == te.selectedColor {
					indicatorLine += "  ▼  "
				} else {
					indicatorLine += "     "
				}
			}
			result += lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00")).Render(indicatorLine) + "\n"
		}

		// Render the color boxes
		for col := 0; col < colsPerRow; col++ {
			idx := row*colsPerRow + col
			if idx >= totalColors {
				break
			}
			color := colorPalette[idx]

			style := lipgloss.NewStyle().
				Background(color).
				Foreground(color) // Make text same as background

			result += style.Render(" ███ ")
		}
		result += "\n"
	}

	return result
}

// RenderProperties renders the list of editable properties
func (te *ThemeEditorState) RenderProperties(theme Theme) string {
	var result string
	result += lipgloss.NewStyle().Bold(true).Render("Theme Properties") + "\n\n"

	for i, prop := range te.properties {
		// Get the current color for this property
		var currentColor lipgloss.Color
		switch prop.FieldName {
		case "NormalFg":
			currentColor = te.editingTheme.NormalFg
		case "SelectedFg":
			currentColor = te.editingTheme.SelectedFg
		case "SelectedBg":
			currentColor = te.editingTheme.SelectedBg
		case "HelpColor":
			currentColor = te.editingTheme.HelpColor
		case "TestCountColor":
			currentColor = te.editingTheme.TestCountColor
		case "TreeSymbolColor":
			currentColor = te.editingTheme.TreeSymbolColor
		case "BorderColor":
			currentColor = te.editingTheme.BorderColor
		case "MenuNormalFg":
			currentColor = te.editingTheme.MenuNormalFg
		case "MenuSelectedFg":
			currentColor = te.editingTheme.MenuSelectedFg
		case "MenuSelectedBg":
			currentColor = te.editingTheme.MenuSelectedBg
		default:
			currentColor = lipgloss.Color("#ffffff")
		}

		// Color preview box
		colorBox := lipgloss.NewStyle().
			Background(currentColor).
			Width(3).
			Render("   ")

		line := fmt.Sprintf("%s  %s", colorBox, prop.Name)

		if te.mode == themeEditorModeProperty && i == te.selectedProperty {
			// Highlight selected property
			line = lipgloss.NewStyle().
				Background(theme.SelectedBg).
				Foreground(theme.SelectedFg).
				Bold(true).
				Render(" " + line + " ")
		} else {
			line = "  " + line
		}

		result += line + "\n"
	}

	result += "\n"
	if te.mode == themeEditorModeProperty {
		result += "Press Enter to select color, ↑↓ to navigate\n"
	} else {
		result += "Press Enter to apply color, ESC to cancel\n"
	}

	return result
}

// RenderPreview renders a preview of the theme
func (te *ThemeEditorState) RenderPreview() string {
	var result string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(te.editingTheme.NormalFg)
	result += titleStyle.Render("Preview") + "\n\n"

	// Preview the tree with current colors from editingTheme
	treeStyle := lipgloss.NewStyle().Foreground(te.editingTheme.TreeSymbolColor)
	normalStyle := lipgloss.NewStyle().Foreground(te.editingTheme.NormalFg)
	selectedStyle := lipgloss.NewStyle().
		Background(te.editingTheme.SelectedBg).
		Foreground(te.editingTheme.SelectedFg).
		Bold(true)
	countStyle := lipgloss.NewStyle().Foreground(te.editingTheme.TestCountColor)

	result += lipgloss.NewStyle().Bold(true).Render("Test List:") + "\n"
	result += treeStyle.Render("├─ ") + normalStyle.Render("  internal/listener") + "\n"
	result += treeStyle.Render("│  ") + countStyle.Render("  (23 tests)") + "\n"
	result += treeStyle.Render("├─ ") + selectedStyle.Render(" internal/database ") + "\n"
	result += treeStyle.Render("│  ") + countStyle.Render("  (15 tests)") + "\n\n"

	// Preview menu colors
	result += lipgloss.NewStyle().Bold(true).Render("Menu Bar:") + "\n"
	menuNormalStyle := lipgloss.NewStyle().
		Foreground(te.editingTheme.MenuNormalFg).
		Padding(0, 1)
	menuSelectedStyle := lipgloss.NewStyle().
		Background(te.editingTheme.MenuSelectedBg).
		Foreground(te.editingTheme.MenuSelectedFg).
		Bold(true).
		Padding(0, 1)

	result += menuNormalStyle.Render("Settings") + " " +
		menuSelectedStyle.Render("Tests") + " " +
		menuNormalStyle.Render("Help") + "\n\n"

	// Show border/separator color
	borderStyle := lipgloss.NewStyle().
		Foreground(te.editingTheme.BorderColor)
	separatorStyle := lipgloss.NewStyle().
		Foreground(te.editingTheme.SeparatorColor)

	result += lipgloss.NewStyle().Bold(true).Render("Border/Separator:") + "\n"
	result += borderStyle.Render("─────────") + " " + separatorStyle.Render("│") + "\n\n"

	// Show help bar text
	helpStyle := lipgloss.NewStyle().Foreground(te.editingTheme.HelpColor)
	result += lipgloss.NewStyle().Bold(true).Render("Help Bar:") + "\n"
	result += helpStyle.Render("F1: menu | ↑↓/jk: scroll | Tab: switch panel") + "\n"

	return result
}
