package main

import "github.com/charmbracelet/lipgloss"

// Common style helpers to reduce lipgloss.NewStyle() duplication

// headerStyle returns a centered header style using theme colors
func (m model) headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(m.currentTheme.SelectedFg).
		Width(m.width).
		Align(lipgloss.Center).
		Padding(0, 1)
}

// selectedItemStyle returns the style for selected menu items/list items
func (m model) selectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(m.currentTheme.SelectedBg).
		Foreground(m.currentTheme.SelectedFg).
		Bold(true).
		Padding(0, 1)
}

// normalItemStyle returns the style for normal (unselected) items
func (m model) normalItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(m.currentTheme.NormalFg).
		Padding(0, 2)
}

// helpBarStyle returns the style for help text at the bottom
func (m model) helpBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(m.currentTheme.HelpColor).
		Width(m.width).
		Align(lipgloss.Center)
}

// contentAreaStyle returns a style for main content areas
func (m model) contentAreaStyle(height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(m.width).
		Height(height).
		Padding(2)
}

// menuItemSelectedStyle returns the style for selected menu items in the menu bar
func (m model) menuItemSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(m.currentTheme.MenuSelectedBg).
		Foreground(m.currentTheme.MenuSelectedFg).
		Bold(true).
		Padding(0, 1)
}

// menuItemActiveStyle returns the style for menu items when menu is active
func (m model) menuItemActiveStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(m.currentTheme.MenuActiveFg).
		Padding(0, 1)
}

// menuItemNormalStyle returns the style for menu items when menu is inactive
func (m model) menuItemNormalStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(m.currentTheme.MenuNormalFg).
		Padding(0, 1)
}

// borderedContentStyle returns a style for bordered content areas (settings, help)
func (m model) borderedContentStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 3).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(m.currentTheme.SelectedBg).
		Padding(1)
}

// boldSectionHeader returns a bold style for section headers
func boldSectionHeader() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true)
}

// currentThemeMarkerStyle returns style for "(current)" marker in theme list
func (m model) currentThemeMarkerStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(m.currentTheme.CoverageGoodFg)
}

// Standalone helper functions that don't need model access

// treeSymbolStyle returns a style for tree drawing symbols (takes theme as parameter)
func treeSymbolStyle(theme Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.TreeSymbolColor)
}

// packageNormalStyle returns a style for normal (unselected) package names
func packageNormalStyle(theme Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.NormalFg)
}

// packageSelectedStyle returns a style for selected package names
func packageSelectedStyle(theme Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.SelectedBg).
		Foreground(theme.SelectedFg).
		Bold(true)
}

// testCountStyle returns a style for test count display
func testCountStyle(theme Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.TestCountColor)
}

// packageTitleStyle returns a style for package list title
func packageTitleStyle(theme Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.NormalFg)
}
