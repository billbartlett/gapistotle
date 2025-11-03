package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleThemeMenuKeys handles keys for the theme menu screen
func handleThemeMenuKeys(m *model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if m.currentScreen != screenThemeMenu {
		return false, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.themeMenuIndex > 0 {
			m.themeMenuIndex--
		}
		return true, nil

	case "down", "j":
		if m.themeMenuIndex < len(m.themeMenuItems)-1 {
			m.themeMenuIndex++
		}
		return true, nil

	case "enter":
		switch m.themeMenuIndex {
		case 0: // Select Theme
			m.currentScreen = screenThemeSelection
			m.themeSelectionMode = themeSelectModeApply
			m.themeSelectionIndex = m.themeIndex
		case 1: // Edit Theme
			m.currentScreen = screenThemeSelection
			m.themeSelectionMode = themeSelectModeEdit
			m.themeSelectionIndex = m.themeIndex
		case 2: // Reload Themes
			// Invalidate cache and reload themes from disk
			InvalidateThemeCache()
			m.themeNames = ListThemes()
			// Update theme index to maintain current theme selection
			m.themeIndex = findThemeIndex(m.themeNames, m.currentTheme.Name)
			// Stay on theme menu screen
		}
		return true, nil
	}

	return false, nil
}

// handleThemeSelectionKeys handles keys for the theme selection screen
func handleThemeSelectionKeys(m *model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if m.currentScreen != screenThemeSelection {
		return false, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.themeSelectionIndex > 0 {
			m.themeSelectionIndex--
		}
		return true, nil

	case "down", "j":
		maxIndex := len(m.themeNames) - 1
		// Allow one extra item in edit mode (Create New Theme)
		if m.themeSelectionMode == themeSelectModeEdit {
			maxIndex = len(m.themeNames)
		}
		if m.themeSelectionIndex < maxIndex {
			m.themeSelectionIndex++
		}
		return true, nil

	case "enter":
		if m.themeSelectionMode == themeSelectModeApply {
			// Apply the selected theme
			m.currentTheme = GetTheme(m.themeNames[m.themeSelectionIndex])
			m.themeIndex = m.themeSelectionIndex
			m.config.CurrentTheme = m.themeNames[m.themeSelectionIndex]
			SaveConfig(m.config, m.configPath)
			m.currentScreen = screenMain
		} else {
			// Check if "Create New Theme" was selected
			if m.themeSelectionIndex == len(m.themeNames) {
				// Create a new theme based on current theme
				newTheme := m.currentTheme
				newTheme.Name = ""
				m.themeEditor = newThemeEditorState(newTheme)
				m.currentScreen = screenThemeEditor
			} else {
				// Go to theme editor with selected theme
				selectedTheme := GetTheme(m.themeNames[m.themeSelectionIndex])
				m.themeEditor = newThemeEditorState(selectedTheme)
				m.currentScreen = screenThemeEditor
			}
		}
		return true, nil
	}

	return false, nil
}

// handleThemeEditorKeys handles keys for the theme editor screen
func handleThemeEditorKeys(m *model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if m.currentScreen != screenThemeEditor {
		return false, nil
	}

	// Save mode is handled by handleTextInput
	if m.themeEditor.mode == themeEditorModeSave {
		return false, nil
	}

	switch msg.String() {
	case "s":
		if m.themeEditor.mode == themeEditorModeProperty {
			// Enter save mode with pre-populated theme name
			m.themeEditor.mode = themeEditorModeSave
			m.themeEditor.saveThemeName = m.themeEditor.editingTheme.Name
			m.themeEditor.saveError = ""
			return true, nil
		}

	case "up", "k":
		if m.themeEditor.mode == themeEditorModeProperty {
			if m.themeEditor.selectedProperty > 0 {
				m.themeEditor.selectedProperty--
			}
		} else {
			// Navigate color palette (up = -8)
			newPos := m.themeEditor.selectedColor - colorPaletteColumns
			if newPos >= 0 {
				m.themeEditor.selectedColor = newPos
			}
		}
		return true, nil

	case "down", "j":
		if m.themeEditor.mode == themeEditorModeProperty {
			if m.themeEditor.selectedProperty < len(m.themeEditor.properties)-1 {
				m.themeEditor.selectedProperty++
			}
		} else {
			// Navigate color palette (down = +8)
			newPos := m.themeEditor.selectedColor + colorPaletteColumns
			if newPos < len(colorPalette) {
				m.themeEditor.selectedColor = newPos
			}
		}
		return true, nil

	case "left", "h":
		if m.themeEditor.mode == themeEditorModeColor {
			if m.themeEditor.selectedColor%colorPaletteColumns > 0 {
				m.themeEditor.selectedColor--
			}
			return true, nil
		}

	case "right", "l":
		if m.themeEditor.mode == themeEditorModeColor {
			if m.themeEditor.selectedColor%colorPaletteColumns < colorPaletteColumns-1 {
				m.themeEditor.selectedColor++
			}
			return true, nil
		}

	case "enter":
		if m.themeEditor.mode == themeEditorModeProperty {
			// Switch to color selection mode
			m.themeEditor.mode = themeEditorModeColor
			// Find the current color in the palette
			currentColor := m.themeEditor.GetCurrentPropertyColor()
			for i, c := range colorPalette {
				if c == currentColor {
					m.themeEditor.selectedColor = i
					break
				}
			}
		} else {
			// Apply selected color and return to property mode
			m.themeEditor.SetCurrentPropertyColor(colorPalette[m.themeEditor.selectedColor])
			m.themeEditor.mode = themeEditorModeProperty
		}
		return true, nil
	}

	return false, nil
}
