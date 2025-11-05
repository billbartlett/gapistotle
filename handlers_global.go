package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// handleTextInput handles text input for theme save mode and hex input mode
// Returns true if the input was consumed
func handleTextInput(m *model, msg tea.KeyMsg) bool {
	if m.currentScreen != screenThemeEditor {
		return false
	}

	// Handle hex input mode
	if m.themeEditor.mode == themeEditorModeHexInput {
		return handleHexInput(m, msg)
	}

	// Handle save mode
	if m.themeEditor.mode != themeEditorModeSave {
		return false
	}

	switch msg.String() {
	case "backspace":
		if len(m.themeEditor.saveThemeName) > 0 {
			m.themeEditor.saveThemeName = m.themeEditor.saveThemeName[:len(m.themeEditor.saveThemeName)-1]
		}
		return true

	case "ctrl+u":
		// Clear entire input
		m.themeEditor.saveThemeName = ""
		return true

	case "enter":
		// Save the theme
		if m.themeEditor.saveThemeName == "" {
			m.themeEditor.saveError = "Theme name cannot be empty"
		} else {
			// Validate theme name (alphanumeric, hyphen, underscore only)
			validName := true
			for _, ch := range m.themeEditor.saveThemeName {
				if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
					(ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
					validName = false
					break
				}
			}

			if !validName {
				m.themeEditor.saveError = "Invalid characters in theme name"
			} else {
				// Set the theme name
				m.themeEditor.editingTheme.Name = m.themeEditor.saveThemeName
				// Save to file
				err := SaveThemeToFile(m.themeEditor.editingTheme, m.themeEditor.saveThemeName)
				if err != nil {
					m.themeEditor.saveError = fmt.Sprintf("Save failed: %v", err)
				} else {
					// Success! Invalidate cache and reload theme list
					InvalidateThemeCache()
					m.themeNames = ListThemes()
					m.currentTheme = GetTheme(m.themeEditor.saveThemeName)
					// Find the index of the new theme
					m.themeIndex = findThemeIndex(m.themeNames, m.themeEditor.saveThemeName)
					// Save to config
					m.config.CurrentTheme = m.themeEditor.saveThemeName
					SaveConfig(m.config, m.configPath)
					// Return to main screen
					m.currentScreen = screenMain
				}
			}
		}
		return true

	case "esc":
		// Cancel save mode
		m.themeEditor.mode = themeEditorModeProperty
		m.themeEditor.saveThemeName = ""
		m.themeEditor.saveError = ""
		return true

	default:
		// Check if it's a navigation key that should be treated as text input
		key := msg.String()
		if key == "h" || key == "j" || key == "k" || key == "l" || key == "s" || key == "t" {
			m.themeEditor.saveThemeName += key
			return true
		}

		// Handle other alphanumeric characters
		if len(key) == 1 {
			char := key[0]
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') || char == '-' || char == '_' {
				m.themeEditor.saveThemeName += key
				return true
			}
		}
	}

	return false
}

// handleHexInput handles text input for hex color input mode
// Returns true if the input was consumed
func handleHexInput(m *model, msg tea.KeyMsg) bool {
	switch msg.String() {
	case "tab":
		// Switch back to color palette mode
		m.themeEditor.mode = themeEditorModeColor
		m.themeEditor.hexInput = ""
		m.themeEditor.hexError = ""
		return true

	case "backspace":
		if len(m.themeEditor.hexInput) > 0 {
			m.themeEditor.hexInput = m.themeEditor.hexInput[:len(m.themeEditor.hexInput)-1]
			m.themeEditor.hexError = "" // Clear error on edit
		}
		return true

	case "ctrl+u":
		// Clear entire input
		m.themeEditor.hexInput = ""
		m.themeEditor.hexError = ""
		return true

	case "enter":
		// Validate and apply the hex color
		hexColor := m.themeEditor.hexInput

		// Add # prefix if not present
		if len(hexColor) > 0 && hexColor[0] != '#' {
			hexColor = "#" + hexColor
		}

		// Validate hex format (should be #RGB, #RRGGBB, or #RRGGBBAA)
		validHex := false
		if len(hexColor) == 4 || len(hexColor) == 7 || len(hexColor) == 9 {
			validHex = true
			for i := 1; i < len(hexColor); i++ {
				ch := hexColor[i]
				if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
					validHex = false
					break
				}
			}
		}

		if !validHex {
			m.themeEditor.hexError = "Invalid hex color (use #RGB or #RRGGBB)"
		} else {
			// Apply the color
			m.themeEditor.SetCurrentPropertyColor(lipgloss.Color(hexColor))
			// Return to color mode
			m.themeEditor.mode = themeEditorModeColor
			m.themeEditor.hexInput = ""
			m.themeEditor.hexError = ""
		}
		return true

	case "esc":
		// Cancel hex input mode
		m.themeEditor.mode = themeEditorModeColor
		m.themeEditor.hexInput = ""
		m.themeEditor.hexError = ""
		return true

	default:
		// Handle hex characters (0-9, a-f, A-F, #)
		key := msg.String()
		if len(key) == 1 {
			char := key[0]
			// Allow # only as first character
			if char == '#' && len(m.themeEditor.hexInput) == 0 {
				m.themeEditor.hexInput += key
				return true
			}
			// Allow hex digits
			if (char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F') {
				// Limit length to 9 characters (#RRGGBBAA)
				if len(m.themeEditor.hexInput) < 9 {
					m.themeEditor.hexInput += key
					m.themeEditor.hexError = "" // Clear error on edit
				}
				return true
			}
		}
	}

	return false
}

// handleGlobalKeys handles keys that work from any screen
func handleGlobalKeys(m *model, msg tea.KeyMsg) (bool, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return true, tea.Quit

	case "f1":
		m.menuActive = !m.menuActive
		return true, nil

	case "esc":
		// ESC behavior depends on context
		if m.currentScreen == screenThemeEditor {
			if m.themeEditor.mode == themeEditorModeSave {
				// Handled by handleTextInput
				return false, nil
			} else if m.themeEditor.mode == themeEditorModeColor {
				m.themeEditor.mode = themeEditorModeProperty
			} else {
				m.currentScreen = screenThemeSelection
			}
		} else if m.currentScreen == screenThemeSelection {
			m.currentScreen = screenThemeMenu
		} else if m.currentScreen == screenThemeMenu {
			m.currentScreen = screenMain
			m.menuActive = false
		} else if m.currentScreen != screenMain {
			// Reset screen-specific state when returning to main
			if m.currentScreen == screenHelp {
				m.helpScroll = 0
			}
			m.currentScreen = screenMain
			m.menuActive = false
		} else if m.currentScreen == screenMain && (m.rightPanelView == viewCoverageGaps || m.rightPanelView == viewDetails) {
			// Return to summary view from details or coverage gaps
			m.rightPanelView = viewSummary
			m.rightPanelScroll = 0
			m.summaryButtonIndex = 0
		} else if m.menuActive {
			m.menuActive = false
		}
		return true, nil
	}

	return false, nil
}

// handleMenuKeys handles menu navigation when menu is active
func handleMenuKeys(m *model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if !m.menuActive {
		return false, nil
	}

	switch msg.String() {
	case "left", "h":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
		return true, nil

	case "right", "l":
		if m.menuIndex < len(m.menuItems)-1 {
			m.menuIndex++
		}
		return true, nil

	case "enter":
		// Handle menu selection
		switch m.menuIndex {
		case 0: // Settings
			m.currentScreen = screenSettings
			m.menuActive = false
		case 1: // Tests
			m.menuActive = false
			// Run tests for all packages
			// Reset view state when running tests
			m.rightPanelView = viewSummary
			m.summaryButtonIndex = 0
			m.rightPanelScroll = 0
			var cmds []tea.Cmd
			for _, pkg := range m.testPackages {
				// Clear old results and errors
				delete(m.testResults, pkg.Name)
				delete(m.testErrors, pkg.Name)
				// Mark test as running
				m.testsRunning[pkg.Name] = true
				cmds = append(cmds, runTestsCmd(pkg.Path, pkg.Name))
			}
			return true, tea.Batch(cmds...)
		case 2: // Theme
			m.currentScreen = screenThemeMenu
			m.themeMenuIndex = 0
			m.menuActive = false
		case 3: // Help
			m.currentScreen = screenHelp
			m.menuActive = false
			m.helpScroll = 0 // Reset scroll when opening help
			m.helpMaxScroll = calculateHelpMaxScroll(m.height)
		case 4: // Quit
			return true, tea.Quit
		}
		return true, nil
	}

	return false, nil
}
