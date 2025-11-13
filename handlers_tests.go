package main

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// handleTestsMenuKeys handles keys for the tests menu screen
func handleTestsMenuKeys(m *model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if m.currentScreen != screenTestsMenu {
		return false, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.testsMenuIndex > 0 {
			m.testsMenuIndex--
		}
		return true, nil

	case "down", "j":
		if m.testsMenuIndex < len(m.testsMenuItems)-1 {
			m.testsMenuIndex++
		}
		return true, nil

	case "enter":
		switch m.testsMenuIndex {
		case 0: // Know It All
			m.currentScreen = screenMain
			// Run tests for all packages SEQUENTIALLY
			// Reset view state when running tests
			m.rightPanelView = viewSummary
			m.summaryButtonIndex = 0
			m.rightPanelScroll = 0

			// Initialize test queue with all packages
			m.testQueue = make([]TestPackage, len(m.testPackages))
			copy(m.testQueue, m.testPackages)
			m.runAllInProgress = true

			// Clear old results and start first test
			if len(m.testQueue) > 0 {
				pkg := m.testQueue[0]
				m.testQueue = m.testQueue[1:]
				// Clear old results and errors
				delete(m.testResults, pkg.Name)
				delete(m.testErrors, pkg.Name)
				// Mark test as running
				m.testsRunning[pkg.Name] = true
				// Use test mode for this specific package's directory
				pkgMode := m.getTestModeForPath(pkg.Path)
				return true, runTestsCmd(pkg.Path, pkg.Name, pkgMode)
			}
			return true, nil
		case 1: // Test Mode
			m.currentScreen = screenTestModeSelection
			// Set index to current mode for the selected package
			var pkgMode testMode = testModeUnit // default
			if m.selectedIndex >= 0 && m.selectedIndex < len(m.testPackages) {
				pkg := m.testPackages[m.selectedIndex]
				pkgMode = m.getTestModeForPath(pkg.Path)
			}
			switch pkgMode {
			case testModeUnit:
				m.testModeIndex = 0
			case testModeIntegration:
				m.testModeIndex = 1
			case testModeAll:
				m.testModeIndex = 2
			}
			return true, nil
		}
		return true, nil
	}

	return false, nil
}

// handleTestModeSelectionKeys handles keys for the test mode selection screen
func handleTestModeSelectionKeys(m *model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if m.currentScreen != screenTestModeSelection {
		return false, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.testModeIndex > 0 {
			m.testModeIndex--
		}
		return true, nil

	case "down", "j":
		if m.testModeIndex < len(m.testModeItems)-1 {
			m.testModeIndex++
		}
		return true, nil

	case "enter":
		// Set the test mode based on selection
		var newMode testMode
		switch m.testModeIndex {
		case 0:
			newMode = testModeUnit
		case 1:
			newMode = testModeIntegration
		case 2:
			newMode = testModeAll
		}

		// Save test mode to config for the currently selected package's directory
		if m.selectedIndex >= 0 && m.selectedIndex < len(m.testPackages) {
			pkg := m.testPackages[m.selectedIndex]
			absPath, err := filepath.Abs(pkg.Path)
			if err == nil {
				m.config.TestModeByDir[absPath] = string(newMode)
				if saveErr := SaveConfig(m.config, m.configPath); saveErr != nil {
					LogWarn("Failed to save test mode config", "error", saveErr)
				} else {
					LogInfo("Test mode saved", "directory", absPath, "mode", newMode, "package", pkg.Name)
				}
			}
		}

		// Return to tests menu
		m.currentScreen = screenTestsMenu
		return true, nil
	}

	return false, nil
}
