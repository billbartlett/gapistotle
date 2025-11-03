package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SaveThemeToFile saves a theme to a config file in themes/ directory
func SaveThemeToFile(theme Theme, filename string) error {
	LogInfo("Saving theme",
		"theme_name", theme.Name,
		"filename", filename,
		"themes_dir", themesDir,
	)

	// Ensure themes directory exists (uses package-level themesDir)
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		LogError("Failed to create themes directory", "error", err)
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	// Create the file path
	filePath := filepath.Join(themesDir, filename+".conf")

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create theme file: %w", err)
	}
	defer file.Close()

	// Write theme properties
	writer := bufio.NewWriter(file)
	fmt.Fprintf(writer, "name=%s\n", theme.Name)
	fmt.Fprintf(writer, "selectedBg=%s\n", theme.SelectedBg)
	fmt.Fprintf(writer, "selectedFg=%s\n", theme.SelectedFg)
	fmt.Fprintf(writer, "normalFg=%s\n", theme.NormalFg)
	fmt.Fprintf(writer, "separatorColor=%s\n", theme.SeparatorColor)
	fmt.Fprintf(writer, "helpColor=%s\n", theme.HelpColor)
	fmt.Fprintf(writer, "testCountColor=%s\n", theme.TestCountColor)
	fmt.Fprintf(writer, "borderColor=%s\n", theme.BorderColor)
	fmt.Fprintf(writer, "treeSymbolColor=%s\n", theme.TreeSymbolColor)
	fmt.Fprintf(writer, "coverageGoodFg=%s\n", theme.CoverageGoodFg)
	fmt.Fprintf(writer, "coverageMediumFg=%s\n", theme.CoverageMediumFg)
	fmt.Fprintf(writer, "coveragePoorFg=%s\n", theme.CoveragePoorFg)
	fmt.Fprintf(writer, "menuNormalFg=%s\n", theme.MenuNormalFg)
	fmt.Fprintf(writer, "menuActiveFg=%s\n", theme.MenuActiveFg)
	fmt.Fprintf(writer, "menuSelectedBg=%s\n", theme.MenuSelectedBg)
	fmt.Fprintf(writer, "menuSelectedFg=%s\n", theme.MenuSelectedFg)

	if err := writer.Flush(); err != nil {
		LogError("Failed to write theme file", "error", err)
		return err
	}

	LogInfo("Theme saved successfully", "theme_name", theme.Name, "file_path", filePath)
	return nil
}

// LoadThemeFromFile loads a theme from a config file
func LoadThemeFromFile(filePath string) (Theme, error) {
	theme := Theme{}

	file, err := os.Open(filePath)
	if err != nil {
		return theme, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "name":
			theme.Name = value
		case "selectedBg":
			theme.SelectedBg = lipgloss.Color(value)
		case "selectedFg":
			theme.SelectedFg = lipgloss.Color(value)
		case "normalFg":
			theme.NormalFg = lipgloss.Color(value)
		case "separatorColor":
			theme.SeparatorColor = lipgloss.Color(value)
		case "helpColor":
			theme.HelpColor = lipgloss.Color(value)
		case "testCountColor":
			theme.TestCountColor = lipgloss.Color(value)
		case "borderColor":
			theme.BorderColor = lipgloss.Color(value)
		case "treeSymbolColor":
			theme.TreeSymbolColor = lipgloss.Color(value)
		case "coverageGoodFg":
			theme.CoverageGoodFg = lipgloss.Color(value)
		case "coverageMediumFg":
			theme.CoverageMediumFg = lipgloss.Color(value)
		case "coveragePoorFg":
			theme.CoveragePoorFg = lipgloss.Color(value)
		case "menuNormalFg":
			theme.MenuNormalFg = lipgloss.Color(value)
		case "menuActiveFg":
			theme.MenuActiveFg = lipgloss.Color(value)
		case "menuSelectedBg":
			theme.MenuSelectedBg = lipgloss.Color(value)
		case "menuSelectedFg":
			theme.MenuSelectedFg = lipgloss.Color(value)
		}
	}

	// Backward compatibility: if testCountColor wasn't set, use helpColor
	if theme.TestCountColor == "" && theme.HelpColor != "" {
		theme.TestCountColor = theme.HelpColor
	}

	return theme, scanner.Err()
}

// LoadCustomThemes loads all custom themes from the themes/ directory
func LoadCustomThemes() map[string]Theme {
	customThemes := make(map[string]Theme)

	// Uses package-level themesDir variable
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		// Directory doesn't exist or can't be read - that's okay
		return customThemes
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".conf") {
			continue
		}

		filePath := filepath.Join(themesDir, entry.Name())
		theme, err := LoadThemeFromFile(filePath)
		if err != nil {
			continue // Skip invalid files
		}

		if theme.Name != "" {
			customThemes[theme.Name] = theme
		}
	}

	return customThemes
}
