package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// TestPackage represents a Go package with tests
type TestPackage struct {
	Name                 string
	Path                 string
	TestFiles            []string
	HasIntegrationTests  bool
}

// hasIntegrationBuildTag checks if a file has integration test build tags
func hasIntegrationBuildTag(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	// Only check first 10 lines (build tags must be at top of file)
	for scanner.Scan() && lineCount < 10 {
		line := strings.TrimSpace(scanner.Text())
		lineCount++

		// Check for both old and new build tag formats
		if strings.Contains(line, "//go:build") && strings.Contains(line, "integration") {
			return true
		}
		if strings.Contains(line, "// +build") && strings.Contains(line, "integration") {
			return true
		}

		// Stop at package declaration
		if strings.HasPrefix(line, "package ") {
			break
		}
	}

	return false
}

// ScanForTests recursively scans a directory for *_test.go files
func ScanForTests(rootPath string) ([]TestPackage, error) {
	LogInfo("Scanning for test packages", "root_path", rootPath)
	packages := make(map[string]*TestPackage)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and vendor
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor") {
			return filepath.SkipDir
		}

		// Look for test files
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_test.go") {
			dir := filepath.Dir(path)
			relDir, _ := filepath.Rel(rootPath, dir)

			if packages[dir] == nil {
				packages[dir] = &TestPackage{
					Name:      relDir,
					Path:      dir,
					TestFiles: []string{},
					HasIntegrationTests: false,
				}
			}
			packages[dir].TestFiles = append(packages[dir].TestFiles, info.Name())

			// Check if this test file has integration build tags
			if hasIntegrationBuildTag(path) {
				packages[dir].HasIntegrationTests = true
				LogDebug("Found integration tests", "package", relDir, "file", info.Name())
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert map to sorted slice
	result := make([]TestPackage, 0, len(packages))
	for _, pkg := range packages {
		sort.Strings(pkg.TestFiles)
		result = append(result, *pkg)
	}

	// Sort by package name
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	LogInfo("Test package scan complete",
		"package_count", len(result),
		"root_path", rootPath,
	)

	return result, nil
}

// RenderTestTree creates a visual tree representation of test packages
func RenderTestTree(packages []TestPackage, selectedIndex int, theme Theme) string {
	if len(packages) == 0 {
		return "No test files found.\n\nRun from a Go project directory."
	}

	var sb strings.Builder

	sb.WriteString(packageTitleStyle(theme).Render("Test Packages") + "\n\n")

	treeStyle := treeSymbolStyle(theme)
	normalStyle := packageNormalStyle(theme)
	selectedStyle := packageSelectedStyle(theme)

	for i, pkg := range packages {
		prefix := "├─ "
		if i == len(packages)-1 {
			prefix = "└─ "
		}

		// Package name
		pkgName := pkg.Name
		if pkgName == "" || pkgName == "." {
			pkgName = "."
		}

		// Add [i] indicator if package has integration tests
		if pkg.HasIntegrationTests {
			pkgName += " [i]"
		}

		sb.WriteString(treeStyle.Render(prefix))

		if i == selectedIndex {
			// Full-width highlight for selected item
			sb.WriteString(selectedStyle.Render(" " + pkgName + " ") + "\n")
		} else {
			sb.WriteString(normalStyle.Render("  " + pkgName) + "\n")
		}

		// Show test file count
		filePrefix := "│  "
		if i == len(packages)-1 {
			filePrefix = "   "
		}

		sb.WriteString(treeStyle.Render(filePrefix) +
			testCountStyle(theme).Render(fmt.Sprintf("  (%d tests)", len(pkg.TestFiles))) + "\n")
	}

	return sb.String()
}
