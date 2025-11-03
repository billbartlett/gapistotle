package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UI constants
const (
	defaultLeftPanelWidth  = 40
	minPanelWidth          = 15
	maxPanelWidthPercent   = 80 // 80% of screen width
	panelResizeIncrement   = 5
	colorPaletteColumns    = 8
)

type appScreen int

const (
	screenMain appScreen = iota
	screenThemeMenu
	screenThemeSelection
	screenThemeEditor
	screenSettings
	screenHelp
)

type panelFocus int

const (
	focusLeftPanel panelFocus = iota
	focusRightPanel
)

type rightPanelView int

const (
	viewSummary rightPanelView = iota
	viewDetails
	viewCoverageGaps
)

// testStartMsg is sent when tests start running
type testStartMsg struct {
	packageName string
}

// testCompleteMsg is sent when a test run completes
type testCompleteMsg struct {
	result *PackageTestResult
}

// testErrorMsg is sent when a test run fails
type testErrorMsg struct {
	packageName string
	err         error
}

type themeSelectionMode int

const (
	themeSelectModeApply themeSelectionMode = iota
	themeSelectModeEdit
)

type model struct {
	width          int
	height         int
	leftPanelWidth int
	ready          bool
	testPackages   []TestPackage
	selectedIndex  int
	scanPath       string
	currentTheme   Theme
	themeNames     []string
	themeIndex     int

	// Menu state
	menuActive     bool
	menuItems      []string
	menuIndex      int
	currentScreen  appScreen

	// Theme menu state
	themeMenuIndex     int
	themeMenuItems     []string

	// Theme selection state
	themeSelectionMode themeSelectionMode
	themeSelectionIndex int

	// Theme editor state
	themeEditor    ThemeEditorState

	// Configuration
	config     Config
	configPath string // Path to config file for saving

	// Test results - maps package path to test results
	testResults map[string]*PackageTestResult
	// Test errors - maps package name to error
	testErrors map[string]error
	// Tests currently running - maps package name to running state
	testsRunning map[string]bool
	// Scan error - error from initial package scan
	scanError error

	// Panel focus and scrolling
	currentFocus       panelFocus
	rightPanelView     rightPanelView
	rightPanelScroll   int
	rightPanelMaxLines int
	summaryButtonIndex int // 0 = TEST DETAILS, 1 = COVERAGE GAPS
	helpScroll         int // Scroll position in help screen
	helpMaxScroll      int // Max scroll lines in help screen
}

func initialModel(scanPath string, flagConfigPath string) model {
	packages, scanErr := ScanForTests(scanPath)

	// Resolve config path from flag, env var, or default
	configPath := ResolveConfigPath(flagConfigPath)

	// Load config to get saved theme and panel width
	config := LoadConfig(configPath)

	// Initialize logger with config settings
	if config.LogPath != "" {
		logLevel := ParseLogLevel(config.LogLevel)
		if err := InitLogger(config.LogPath, logLevel); err != nil {
			// If logger init fails, continue without logging (non-fatal)
			fmt.Fprintf(os.Stderr, "Warning: Failed to initialize logger: %v\n", err)
		} else {
			LogInfo("Gapistotle starting",
				"version", "1.0.0",
				"scan_path", scanPath,
				"config_path", configPath,
			)
		}
	}

	// Set themes directory based on config (XDG or custom path)
	SetThemesDir(GetThemesDir(config))

	themeNames := ListThemes()
	currentTheme := GetTheme(config.CurrentTheme)

	// Find the index of the current theme
	themeIdx := findThemeIndex(themeNames, config.CurrentTheme)

	return model{
		leftPanelWidth:      config.LeftPanelWidth,
		testPackages:        packages,
		selectedIndex:       0,
		scanPath:            scanPath,
		currentTheme:        currentTheme,
		themeNames:          themeNames,
		themeIndex:          themeIdx,
		menuActive:          false,
		menuItems:           []string{"Settings", "Tests", "Theme", "Help", "Quit"},
		menuIndex:           0,
		currentScreen:       screenMain,
		themeMenuIndex:      0,
		themeMenuItems:      []string{"Select Theme", "Edit Theme", "Reload Themes"},
		themeSelectionMode:  themeSelectModeApply,
		themeSelectionIndex: 0,
		config:              config,
		configPath:          configPath,
		testResults:         make(map[string]*PackageTestResult),
		testErrors:          make(map[string]error),
		testsRunning:        make(map[string]bool),
		scanError:           scanErr,
		currentFocus:        focusLeftPanel,
		rightPanelView:      viewSummary,
		rightPanelScroll:    0,
		rightPanelMaxLines:  0,
		summaryButtonIndex:  0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// findThemeIndex returns the index of the theme with the given name
// Returns 0 if the theme is not found
func findThemeIndex(themeNames []string, themeName string) int {
	for i, name := range themeNames {
		if name == themeName {
			return i
		}
	}
	return 0
}

// runTestsCmd runs tests for a package and returns the result
func runTestsCmd(packageDir string, packageName string) tea.Cmd {
	return func() tea.Msg {
		result, err := RunTests(packageDir, packageName)
		if err != nil {
			return testErrorMsg{packageName: packageName, err: err}
		}
		return testCompleteMsg{result: result}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return &m, nil

	case testCompleteMsg:
		// Store test result
		if msg.result != nil {
			m.testResults[msg.result.PackagePath] = msg.result
			// Clear running state
			delete(m.testsRunning, msg.result.PackagePath)
		}
		return &m, nil

	case testErrorMsg:
		// Store test error
		m.testErrors[msg.packageName] = msg.err
		// Clear running state
		delete(m.testsRunning, msg.packageName)
		return &m, nil

	case tea.KeyMsg:
		// Priority 1: Handle text input (highest priority to prevent navigation interference)
		if handleTextInput(&m, msg) {
			return &m, nil
		}

		// Priority 2: Handle global keys (ctrl+c, F1, ESC)
		if handled, cmd := handleGlobalKeys(&m, msg); handled {
			return &m, cmd
		}

		// Priority 3: Handle menu keys (when menu is active)
		if handled, cmd := handleMenuKeys(&m, msg); handled {
			return &m, cmd
		}

		// Priority 4: Handle screen-specific keys
		if handled, cmd := handleMainScreenKeys(&m, msg); handled {
			return &m, cmd
		}
		if handled, cmd := handleThemeMenuKeys(&m, msg); handled {
			return &m, cmd
		}
		if handled, cmd := handleThemeSelectionKeys(&m, msg); handled {
			return &m, cmd
		}
		if handled, cmd := handleHelpKeys(&m, msg); handled {
			return &m, cmd
		}
		if handled, cmd := handleThemeEditorKeys(&m, msg); handled {
			return &m, cmd
		}
	}

	return &m, nil
}

func (m *model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Render menu bar
	menuBar := m.renderMenuBar()

	// Route to appropriate screen
	var content string
	switch m.currentScreen {
	case screenMain:
		content = m.renderMainScreen()
	case screenThemeMenu:
		content = m.renderThemeMenu()
	case screenThemeSelection:
		content = m.renderThemeSelection()
	case screenThemeEditor:
		content = m.renderThemeEditor()
	case screenSettings:
		content = m.renderSettings()
	case screenHelp:
		content = m.renderHelp()
	default:
		content = m.renderMainScreen()
	}

	return lipgloss.JoinVertical(lipgloss.Left, menuBar, content)
}

func main() {
	// Parse command-line flags
	configPath := flag.String("c", "", "path to config file")
	flag.Parse()

	// Get scan path from remaining args, or use current directory
	scanPath := "."
	if flag.NArg() > 0 {
		scanPath = flag.Arg(0)
	}

	m := initialModel(scanPath, *configPath)
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
