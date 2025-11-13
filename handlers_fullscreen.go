package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleFullScreenKeys handles keys for full-screen views (test results and coverage gaps)
func handleFullScreenKeys(m *model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if m.currentScreen != screenFullTestResults && m.currentScreen != screenFullCoverageGaps {
		return false, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.fullScreenScroll > 0 {
			m.fullScreenScroll--
		}
		return true, nil

	case "down", "j":
		// Max scroll will be calculated in render
		m.fullScreenScroll++
		return true, nil

	case "g":
		// Go to top
		m.fullScreenScroll = 0
		return true, nil

	case "G":
		// Go to bottom (will be capped in render)
		m.fullScreenScroll = 999999
		return true, nil

	case "pgup", "ctrl+u":
		pageSize := m.height - 4 // Account for header/footer
		if pageSize < 1 {
			pageSize = 10
		}
		m.fullScreenScroll -= pageSize
		if m.fullScreenScroll < 0 {
			m.fullScreenScroll = 0
		}
		return true, nil

	case "pgdown", "ctrl+d":
		pageSize := m.height - 4
		if pageSize < 1 {
			pageSize = 10
		}
		m.fullScreenScroll += pageSize
		return true, nil
	}

	return false, nil
}
