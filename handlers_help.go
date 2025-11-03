package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// calculateHelpMaxScroll calculates the max scroll for help screen
// Help content has approximately 50 lines
func calculateHelpMaxScroll(screenHeight int) int {
	helpContentLines := 50 // Approximate number of lines in help content
	visibleLines := HelpScreenPageSize(screenHeight)
	maxScroll := helpContentLines - visibleLines
	if maxScroll < 0 {
		return 0
	}
	return maxScroll
}

// handleHelpKeys handles keys for the help screen
func handleHelpKeys(m *model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if m.currentScreen != screenHelp {
		return false, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.helpScroll > 0 {
			m.helpScroll--
		}
		return true, nil

	case "down", "j":
		if m.helpScroll < m.helpMaxScroll {
			m.helpScroll++
		}
		return true, nil

	case "g":
		// Jump to top
		m.helpScroll = 0
		return true, nil

	case "G":
		// Jump to bottom
		m.helpScroll = m.helpMaxScroll
		return true, nil

	case "pgup":
		// Page up
		pageSize := HelpScreenPageSize(m.height)
		m.helpScroll -= pageSize
		if m.helpScroll < 0 {
			m.helpScroll = 0
		}
		return true, nil

	case "pgdown":
		// Page down
		pageSize := HelpScreenPageSize(m.height)
		m.helpScroll += pageSize
		if m.helpScroll > m.helpMaxScroll {
			m.helpScroll = m.helpMaxScroll
		}
		return true, nil
	}

	return false, nil
}
