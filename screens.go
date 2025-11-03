package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	noTestResultsMessage = "No tests have been run yet.\n\nSelect a package from the left and\npress Enter to run tests and view coverage."
)

func (m model) renderMenuBar() string {
	var menuItems []string

	for i, item := range m.menuItems {
		if m.menuActive && i == m.menuIndex {
			// Highlight selected menu item
			menuItems = append(menuItems, m.menuItemSelectedStyle().Render(item))
		} else if m.menuActive {
			// Normal menu item when menu is active
			menuItems = append(menuItems, m.menuItemActiveStyle().Render(item))
		} else {
			// Dimmed when menu is not active
			menuItems = append(menuItems, m.menuItemNormalStyle().Render(item))
		}
	}

	menuBar := lipgloss.JoinHorizontal(lipgloss.Left, menuItems...)

	// Create menu content without border
	menuStyle := lipgloss.NewStyle().
		Width(m.width).
		Foreground(m.currentTheme.NormalFg)

	menuContent := menuStyle.Render(menuBar)

	// Create custom border with different colors for left and right sections
	leftBorderWidth := m.leftPanelWidth
	rightBorderWidth := m.width - m.leftPanelWidth

	leftBorderColor := m.currentTheme.BorderColor
	rightBorderColor := m.currentTheme.BorderColor

	if m.currentScreen == screenMain {
		if m.currentFocus == focusLeftPanel {
			leftBorderColor = m.currentTheme.SelectedBg
		} else if m.currentFocus == focusRightPanel {
			rightBorderColor = m.currentTheme.SelectedBg
		}
	}

	leftBorderStyle := lipgloss.NewStyle().Foreground(leftBorderColor)
	rightBorderStyle := lipgloss.NewStyle().Foreground(rightBorderColor)

	leftBorder := HR(leftBorderWidth, leftBorderStyle, "─")
	rightBorder := HR(rightBorderWidth, rightBorderStyle, "─")

	border := leftBorder + rightBorder

	return menuContent + "\n" + border
}

func (m *model) renderMainScreen() string {
	// Compute layout dimensions
	frame := Frame{Width: m.width, Height: m.height}
	split := NewSplit(frame, m.leftPanelWidth, m.config.MinPanelWidth, (m.width*m.config.MaxPanelWidthPercent)/100)

	// Define styles
	leftStyle := lipgloss.NewStyle().
		Width(split.LeftW - 1).
		Height(split.ContentH).
		Padding(0, 1)

	rightStyle := lipgloss.NewStyle().
		Width(split.RightW).
		Height(split.ContentH).
		Padding(0, 2)

	separatorStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.SeparatorColor)

	helpStyle := m.helpBarStyle()

	// Left panel content - show scan error or render test tree
	var leftContent string
	if m.scanError != nil {
		leftContent = fmt.Sprintf("Scan Error\n\nFailed to scan for test packages.\n\nPath: %s\n\nError:\n%v\n\nPlease check the path and try again.", m.scanPath, m.scanError)
	} else {
		leftContent = RenderTestTree(m.testPackages, m.selectedIndex, m.currentTheme)
	}

	// Right panel content - show test results if available
	var rightContent string
	if m.selectedIndex < len(m.testPackages) {
		selectedPkg := m.testPackages[m.selectedIndex]
		// Check if test is currently running
		if running, exists := m.testsRunning[selectedPkg.Name]; exists && running {
			rightContent = fmt.Sprintf("Running tests...\n\nPackage: %s\n\nPlease wait while tests execute.\nThis may take a few moments for larger test suites.", selectedPkg.Name)
		} else if err, exists := m.testErrors[selectedPkg.Name]; exists {
			// Check for errors
			rightContent = fmt.Sprintf("Test Error\n\nFailed to run tests for package: %s\n\nError:\n%v", selectedPkg.Name, err)
		} else if result, exists := m.testResults[selectedPkg.Name]; exists {
			// Show test results based on current view mode
			switch m.rightPanelView {
			case viewSummary:
				rightContent = FormatTestResultSummary(result, m.currentTheme, m.summaryButtonIndex)
			case viewDetails:
				rightContent = FormatTestResult(result, m.currentTheme)
			case viewCoverageGaps:
				rightContent = FormatCoverageGaps(result, m.currentTheme)
			default:
				rightContent = FormatTestResultSummary(result, m.currentTheme, m.summaryButtonIndex)
			}
		} else {
			// No results yet
			rightContent = noTestResultsMessage
		}
	} else {
		rightContent = noTestResultsMessage
	}

	// Handle scrolling for right panel
	rightContentLines := strings.Split(rightContent, "\n")
	totalLines := len(rightContentLines)
	visibleLines := split.ContentH - 2 // Account for padding

	// Update max scroll
	m.rightPanelMaxLines = totalLines - visibleLines
	if m.rightPanelMaxLines < 0 {
		m.rightPanelMaxLines = 0
	}

	// Ensure scroll is within bounds
	if m.rightPanelScroll > m.rightPanelMaxLines {
		m.rightPanelScroll = m.rightPanelMaxLines
	}

	// Extract visible portion of content
	var visibleContent string
	if totalLines > visibleLines {
		endLine := m.rightPanelScroll + visibleLines
		if endLine > totalLines {
			endLine = totalLines
		}
		visibleContent = strings.Join(rightContentLines[m.rightPanelScroll:endLine], "\n")
	} else {
		visibleContent = rightContent
	}

	// Render panels with separator
	leftPanel := leftStyle.Render(leftContent)

	// Create vertical separator using layout helper
	separator := VSep(split.ContentH, separatorStyle)

	rightPanel := rightStyle.Render(visibleContent)

	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, separator, rightPanel)

	// Help text - show different text based on focus and view
	var helpText string
	if m.currentFocus == focusLeftPanel {
		helpText = fmt.Sprintf("F1: menu | ↑↓/jk: navigate | Tab: switch panel | ]/[: resize | t: theme (%s) | q: quit", m.currentTheme.Name)
	} else {
		// In right panel - show context-specific help
		if m.rightPanelView == viewCoverageGaps {
			helpText = fmt.Sprintf("F1: menu | ESC: return to summary | Tab: switch panel | t: theme (%s) | q: quit", m.currentTheme.Name)
		} else if m.rightPanelView == viewSummary {
			helpText = fmt.Sprintf("F1: menu | Enter: select | Tab: switch panel | ]/[: resize | t: theme (%s) | q: quit", m.currentTheme.Name)
		} else {
			helpText = fmt.Sprintf("F1: menu | ↑↓/jk: scroll | Tab: switch panel | ]/[: resize | t: theme (%s) | q: quit", m.currentTheme.Name)
		}
	}
	help := helpStyle.Render(helpText)

	// Combine everything
	return lipgloss.JoinVertical(lipgloss.Left, panels, help)
}

func (m model) renderThemeMenu() string {
	contentHeight := m.height - MenuBarH

	header := m.headerStyle().Render("═══ THEME MENU ═══")

	// Content
	contentStyle := m.contentAreaStyle(contentHeight - 2)

	var content string
	content += boldSectionHeader().Render("Theme Options") + "\n\n"

	for i, item := range m.themeMenuItems {
		if i == m.themeMenuIndex {
			content += m.selectedItemStyle().Render(" > "+item+" ") + "\n"
		} else {
			content += m.normalItemStyle().Render("   "+item) + "\n"
		}
	}

	footer := m.helpBarStyle().Render("↑↓/jk: navigate | Enter: select | ESC: back")

	return lipgloss.JoinVertical(lipgloss.Left, header, contentStyle.Render(content), footer)
}

func (m model) renderThemeSelection() string {
	contentHeight := m.height - MenuBarH

	var title string
	if m.themeSelectionMode == themeSelectModeApply {
		title = "═══ SELECT THEME ═══"
	} else {
		title = "═══ SELECT THEME TO EDIT ═══"
	}
	header := m.headerStyle().Render(title)

	// Content
	contentStyle := m.contentAreaStyle(contentHeight - 2)

	var content string
	content += boldSectionHeader().Render("Available Themes") + "\n\n"

	for i, themeName := range m.themeNames {
		var line string
		if i == m.themeSelectionIndex {
			line = m.selectedItemStyle().Render(" > " + themeName + " ")
		} else {
			line = m.normalItemStyle().Render("   " + themeName)
		}

		// Mark current theme
		if i == m.themeIndex {
			line += " " + m.currentThemeMarkerStyle().Render("(current)")
		}

		content += line + "\n"
	}

	// Add "Create New Theme" option when in edit mode
	if m.themeSelectionMode == themeSelectModeEdit {
		createNewLabel := "Create New Theme"
		createNewIndex := len(m.themeNames)
		var line string
		if createNewIndex == m.themeSelectionIndex {
			line = m.selectedItemStyle().Render(" > " + createNewLabel + " ")
		} else {
			line = m.normalItemStyle().Render("   " + createNewLabel)
		}
		content += "\n" + line + "\n"
	}

	var footer string
	if m.themeSelectionMode == themeSelectModeApply {
		footer = m.helpBarStyle().Render("↑↓/jk: navigate | Enter: apply theme | ESC: back")
	} else {
		footer = m.helpBarStyle().Render("↑↓/jk: navigate | Enter: edit theme | ESC: back")
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, contentStyle.Render(content), footer)
}

func (m model) renderThemeEditor() string {
	contentHeight := m.height - MenuBarH

	header := m.headerStyle().Render("═══ THEME EDITOR ═══")

	// If in save mode, show save dialog
	if m.themeEditor.mode == themeEditorModeSave {
		dialogStyle := lipgloss.NewStyle().
			Width(m.width).
			Height(contentHeight - 2).
			Padding(2).
			Align(lipgloss.Center, lipgloss.Center)

		var content string
		content += boldSectionHeader().Render("Save Theme") + "\n\n"
		content += "Enter theme name:\n\n"

		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(m.currentTheme.BorderColor).
			Padding(0, 1).
			Width(40)

		content += inputStyle.Render(m.themeEditor.saveThemeName+"_") + "\n\n"

		if m.themeEditor.saveError != "" {
			errorStyle := lipgloss.NewStyle().
				Foreground(m.currentTheme.CoveragePoorFg).
				Background(lipgloss.Color("#330000")).
				Bold(true).
				Padding(0, 1)
			content += "\n"
			content += errorStyle.Render("[ FAIL ] "+m.themeEditor.saveError) + "\n\n"
		}

		content += lipgloss.NewStyle().Foreground(m.currentTheme.HelpColor).Render("Only alphanumeric, hyphen, and underscore allowed")

		footer := m.helpBarStyle().Render("Enter: save | Ctrl+U: clear | ESC: cancel")

		return lipgloss.JoinVertical(lipgloss.Left, header, dialogStyle.Render(content), footer)
	}

	// Left column: Color palette and properties
	leftColumn := m.themeEditor.RenderColorPalette() + "\n\n" + m.themeEditor.RenderProperties(m.currentTheme)

	leftStyle := lipgloss.NewStyle().
		Width(m.width / 2).
		Height(contentHeight - 3).
		Padding(1)

	// Right column: Preview
	rightColumn := m.themeEditor.RenderPreview()

	rightStyle := lipgloss.NewStyle().
		Width(m.width/2 - 2).
		Height(contentHeight - 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.currentTheme.BorderColor).
		Padding(1)

	// Combine columns
	columns := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftColumn),
		rightStyle.Render(rightColumn),
	)

	var help string
	if m.themeEditor.mode == themeEditorModeProperty {
		help = "↑↓: navigate | Enter: edit color | s: save theme | ESC: cancel"
	} else {
		help = "←→↑↓/hjkl: select color | Enter: apply | ESC: cancel"
	}

	footer := m.helpBarStyle().Render(help)

	return lipgloss.JoinVertical(lipgloss.Left, header, columns, footer)
}

func (m model) renderSettings() string {
	content := boldSectionHeader().Render("SETTINGS") + "\n\n"
	content += "Coming soon!\n\n"
	content += "Press ESC to return to main screen"

	return m.borderedContentStyle().Render(content)
}

func (m model) renderHelp() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.MenuActiveFg)
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(m.currentTheme.SelectedFg)
	keyStyle := lipgloss.NewStyle().Foreground(m.currentTheme.CoverageGoodFg)

	content := headerStyle.Render("GAPISTOTLE HELP") + "\n"
	content += lipgloss.NewStyle().Foreground(m.currentTheme.HelpColor).Render("Know Thy Code") + "\n\n"

	// General Navigation
	content += sectionStyle.Render("═══ GENERAL ═══") + "\n"
	content += keyStyle.Render("  ?         ") + " - Show this help screen\n"
	content += keyStyle.Render("  F1        ") + " - Toggle menu\n"
	content += keyStyle.Render("  ←→ / h l  ") + " - Navigate menu items\n"
	content += keyStyle.Render("  Enter     ") + " - Select menu item / Run tests\n"
	content += keyStyle.Render("  Tab       ") + " - Switch focus (left/right panel)\n"
	content += keyStyle.Render("  ESC       ") + " - Cancel / Go back / Deactivate menu\n"
	content += keyStyle.Render("  q         ") + " - Quit (from main screen)\n"
	content += keyStyle.Render("  Ctrl+C    ") + " - Force quit\n\n"

	// Test Package Navigation
	content += sectionStyle.Render("═══ TEST PACKAGES ═══") + "\n"
	content += keyStyle.Render("  ↑↓ / j k  ") + " - Navigate test package list\n"
	content += keyStyle.Render("  Enter     ") + " - Run tests for selected package\n"
	content += keyStyle.Render("  ] / [     ") + " - Widen / narrow left panel\n\n"

	// Test Results Navigation
	content += sectionStyle.Render("═══ TEST RESULTS ═══") + "\n"
	content += keyStyle.Render("  Tab       ") + " - Focus right panel\n"
	content += keyStyle.Render("  ↑↓ / j k  ") + " - Navigate buttons / Scroll output\n"
	content += keyStyle.Render("  g         ") + " - Jump to top of output\n"
	content += keyStyle.Render("  G         ") + " - Jump to bottom of output\n"
	content += keyStyle.Render("  PgUp      ") + " - Scroll up one page\n"
	content += keyStyle.Render("  PgDn      ") + " - Scroll down one page\n"
	content += keyStyle.Render("  Enter     ") + " - Select button (TEST DETAILS / COVERAGE GAPS)\n"
	content += keyStyle.Render("  ESC       ") + " - Return to summary view\n\n"

	// Themes
	content += sectionStyle.Render("═══ THEMES ═══") + "\n"
	content += keyStyle.Render("  t         ") + " - Cycle through themes quickly\n"
	content += keyStyle.Render("  F1 → Theme") + " - Access theme menu\n"
	content += "    - Select Theme: Choose from all available themes\n"
	content += "    - Edit Theme: Customize colors with live preview\n"
	content += "    - Reload Themes: Refresh theme list\n\n"

	// Tips
	content += sectionStyle.Render("═══ TIPS ═══") + "\n"
	content += "  • Press " + keyStyle.Render("Enter") + " on any package to run its tests\n"
	content += "  • Use " + keyStyle.Render("Tab") + " to focus right panel and navigate results\n"
	content += "  • Coverage gaps show functions with biggest impact on coverage\n"
	content += "  • Theme changes save automatically to " + keyStyle.Render("~/.gapistotle.conf") + "\n"
	content += "  • Custom themes go in " + keyStyle.Render("~/.config/gapistotle/themes/") + "\n\n"

	content += lipgloss.NewStyle().Foreground(m.currentTheme.HelpColor).Render("Press ESC to return | ↑↓/jk: scroll | g/G: top/bottom | PgUp/PgDn: page")

	// Handle scrolling
	contentLines := strings.Split(content, "\n")
	totalLines := len(contentLines)
	// borderedContentStyle: Height(m.height-3) with DoubleBorder(2 lines) and Padding(1 = 2 lines)
	// Menu bar takes 4 lines, so available = m.height - 4
	// Border and padding = 4 lines, so actual content area = m.height - 4 - 4 = m.height - 8
	visibleLines := HelpScreenPageSize(m.height)

	// Update max scroll
	m.helpMaxScroll = totalLines - visibleLines
	if m.helpMaxScroll < 0 {
		m.helpMaxScroll = 0
	}

	// Ensure scroll is within bounds
	if m.helpScroll > m.helpMaxScroll {
		m.helpScroll = m.helpMaxScroll
	}

	// Extract visible portion of content
	var visibleContent string
	if totalLines > visibleLines {
		endLine := m.helpScroll + visibleLines
		if endLine > totalLines {
			endLine = totalLines
		}
		visibleContent = strings.Join(contentLines[m.helpScroll:endLine], "\n")
	} else {
		visibleContent = content
	}

	return m.borderedContentStyle().Render(visibleContent)
}
