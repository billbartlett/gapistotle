package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleMainScreenKeys handles keys specific to the main screen
func handleMainScreenKeys(m *model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if m.currentScreen != screenMain || m.menuActive {
		return false, nil
	}

	switch msg.String() {
	case "tab":
		// Switch panel focus
		if m.currentFocus == focusLeftPanel {
			m.currentFocus = focusRightPanel
		} else {
			m.currentFocus = focusLeftPanel
		}
		return true, nil

	case "enter":
		// Handle Enter based on focus
		if m.currentFocus == focusLeftPanel {
			// Run tests for the selected package
			if m.selectedIndex < len(m.testPackages) {
				pkg := m.testPackages[m.selectedIndex]
				// Clear old results and errors for this package
				delete(m.testResults, pkg.Name)
				delete(m.testErrors, pkg.Name)
				// Mark test as running
				m.testsRunning[pkg.Name] = true
				// Reset view state when running a new test
				m.rightPanelView = viewSummary
				m.summaryButtonIndex = 0
				m.rightPanelScroll = 0
				return true, runTestsCmd(pkg.Path, pkg.Name)
			}
		} else if m.currentFocus == focusRightPanel && m.rightPanelView == viewSummary {
			// User pressed Enter on a button - navigate based on which button
			if m.summaryButtonIndex == 0 {
				// TEST DETAILS
				m.rightPanelView = viewDetails
			} else {
				// COVERAGE GAPS
				m.rightPanelView = viewCoverageGaps
			}
			m.rightPanelScroll = 0
			return true, nil
		}
		return true, nil

	case "up", "k":
		if m.currentFocus == focusLeftPanel {
			// Navigate packages
			if m.selectedIndex > 0 {
				m.selectedIndex--
				// Reset view state when changing package
				m.rightPanelView = viewSummary
				m.summaryButtonIndex = 0
				m.rightPanelScroll = 0
			}
		} else {
			// Right panel focused
			if m.rightPanelView == viewSummary {
				// Navigate between buttons
				if m.summaryButtonIndex > 0 {
					m.summaryButtonIndex--
				}
			} else {
				// Scroll right panel up
				if m.rightPanelScroll > 0 {
					m.rightPanelScroll--
				}
			}
		}
		return true, nil

	case "down", "j":
		if m.currentFocus == focusLeftPanel {
			// Navigate packages
			if m.selectedIndex < len(m.testPackages)-1 {
				m.selectedIndex++
				// Reset view state when changing package
				m.rightPanelView = viewSummary
				m.summaryButtonIndex = 0
				m.rightPanelScroll = 0
			}
		} else {
			// Right panel focused
			if m.rightPanelView == viewSummary {
				// Navigate between buttons (0 = TEST DETAILS, 1 = COVERAGE GAPS)
				if m.summaryButtonIndex < 1 {
					m.summaryButtonIndex++
				}
			} else {
				// Scroll right panel down
				if m.rightPanelScroll < m.rightPanelMaxLines {
					m.rightPanelScroll++
				}
			}
		}
		return true, nil

	case "g":
		// Jump to top of right panel
		if m.currentFocus == focusRightPanel && m.rightPanelView != viewSummary {
			m.rightPanelScroll = 0
		}
		return true, nil

	case "G":
		// Jump to bottom of right panel
		if m.currentFocus == focusRightPanel && m.rightPanelView != viewSummary {
			m.rightPanelScroll = m.rightPanelMaxLines
		}
		return true, nil

	case "pgup":
		// Page up in right panel
		if m.currentFocus == focusRightPanel && m.rightPanelView != viewSummary {
			frame := Frame{Width: m.width, Height: m.height}
			split := NewSplit(frame, m.leftPanelWidth, m.config.MinPanelWidth, (m.width*m.config.MaxPanelWidthPercent)/100)
			pageSize := RightPanelPageSize(split)
			m.rightPanelScroll -= pageSize
			if m.rightPanelScroll < 0 {
				m.rightPanelScroll = 0
			}
		}
		return true, nil

	case "pgdown":
		// Page down in right panel
		if m.currentFocus == focusRightPanel && m.rightPanelView != viewSummary {
			frame := Frame{Width: m.width, Height: m.height}
			split := NewSplit(frame, m.leftPanelWidth, m.config.MinPanelWidth, (m.width*m.config.MaxPanelWidthPercent)/100)
			pageSize := RightPanelPageSize(split)
			m.rightPanelScroll += pageSize
			if m.rightPanelScroll > m.rightPanelMaxLines {
				m.rightPanelScroll = m.rightPanelMaxLines
			}
		}
		return true, nil

	case "]":
		// Widen left panel
		maxWidth := m.width * m.config.MaxPanelWidthPercent / 100
		if m.leftPanelWidth < maxWidth {
			m.leftPanelWidth += m.config.PanelResizeIncrement
			m.config.LeftPanelWidth = m.leftPanelWidth
			SaveConfig(m.config, m.configPath)
		}
		return true, nil

	case "[":
		// Narrow left panel
		if m.leftPanelWidth > m.config.MinPanelWidth {
			m.leftPanelWidth -= m.config.PanelResizeIncrement
			m.config.LeftPanelWidth = m.leftPanelWidth
			SaveConfig(m.config, m.configPath)
		}
		return true, nil

	case "t":
		// Cycle through themes
		m.themeIndex = (m.themeIndex + 1) % len(m.themeNames)
		m.currentTheme = GetTheme(m.themeNames[m.themeIndex])
		m.config.CurrentTheme = m.themeNames[m.themeIndex]
		SaveConfig(m.config, m.configPath)
		return true, nil

	case "?":
		// Show help modal
		m.currentScreen = screenHelp
		m.menuActive = false
		m.helpScroll = 0 // Reset scroll when opening help
		// Pre-calculate help max scroll (help content is ~50 lines)
		m.helpMaxScroll = calculateHelpMaxScroll(m.height)
		return true, nil

	case "q":
		return true, tea.Quit
	}

	return false, nil
}
