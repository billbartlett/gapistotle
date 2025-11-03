package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// FormatTestResultSummary formats a compact summary of test results
func FormatTestResultSummary(result *PackageTestResult, theme Theme, selectedButton int) string {
	var output strings.Builder

	// Styles
	separatorStyle := lipgloss.NewStyle().Foreground(theme.TreeSymbolColor)
	normalStyle := lipgloss.NewStyle().Foreground(theme.NormalFg)
	metricStyle := lipgloss.NewStyle().Foreground(theme.MenuActiveFg)
	passStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

	separator := separatorStyle.Render("========================================")
	output.WriteString(separator + "\n")
	output.WriteString(normalStyle.Render(fmt.Sprintf("Package: %s", result.PackagePath)) + "\n")
	output.WriteString(separator + "\n")

	// Status summary with color based on pass/fail
	var statusStyled string
	if result.Status == "PASS" {
		statusStyled = passStyle.Render("PASS")
	} else {
		statusStyled = failStyle.Render("FAIL")
	}
	output.WriteString(normalStyle.Render("Status: ") + statusStyled +
		normalStyle.Render(fmt.Sprintf(" (%d/%d passed)", result.PassedTests, result.TotalTests)) + "\n")

	output.WriteString(normalStyle.Render("Coverage: ") +
		metricStyle.Render(fmt.Sprintf("%.1f%%", result.Coverage)) + "\n")

	// Show time or "(cached)" if no duration
	if result.Duration > 0 {
		output.WriteString(normalStyle.Render("Total Time: ") +
			metricStyle.Render(formatDuration(result.Duration)) + "\n\n")
	} else {
		output.WriteString(normalStyle.Render("Total Time: ") +
			metricStyle.Render("(cached)") + "\n\n")
	}

	// Add buttons - use theme selection colors for visibility
	selectedStyle := lipgloss.NewStyle().
		Foreground(theme.SelectedFg).
		Background(theme.SelectedBg).
		Bold(true)
	normalButtonStyle := lipgloss.NewStyle().
		Foreground(theme.NormalFg)

	// TEST DETAILS button
	if selectedButton == 0 {
		output.WriteString(selectedStyle.Render("[ TEST DETAILS ]") + "\n")
	} else {
		output.WriteString(normalButtonStyle.Render("[ TEST DETAILS ]") + "\n")
	}

	// COVERAGE GAPS button
	if selectedButton == 1 {
		output.WriteString(selectedStyle.Render("[ COVERAGE GAPS ]") + "\n")
	} else {
		output.WriteString(normalButtonStyle.Render("[ COVERAGE GAPS ]") + "\n")
	}

	return output.String()
}

// FormatCoverageGaps formats coverage gaps analysis with ASCII progress bars
func FormatCoverageGaps(result *PackageTestResult, theme Theme) string {
	var output strings.Builder

	// Styles
	separatorStyle := lipgloss.NewStyle().Foreground(theme.TreeSymbolColor)
	normalStyle := lipgloss.NewStyle().Foreground(theme.NormalFg)

	separator := separatorStyle.Render("========================================")
	output.WriteString(separator + "\n")
	output.WriteString(normalStyle.Render("COVERAGE ANALYSIS") + "\n")
	output.WriteString(separator + "\n\n")

	// Current coverage with progress bar
	output.WriteString(normalStyle.Render("Current Coverage:") + "\n")
	currentBar := renderProgressBar(result.Coverage, theme)
	output.WriteString(currentBar + "\n\n")

	// Count functions by coverage level
	var uncoveredCount, partialCount, coveredCount int
	for _, fc := range result.FunctionCoverages {
		if fc.CoveragePercent == 0.0 {
			uncoveredCount++
		} else if fc.CoveragePercent < 100.0 {
			partialCount++
		} else {
			coveredCount++
		}
	}

	// Summary of function coverage
	totalFuncs := len(result.FunctionCoverages)
	output.WriteString(normalStyle.Render(fmt.Sprintf("Functions: %d total", totalFuncs)) + "\n")
	output.WriteString(normalStyle.Render(fmt.Sprintf("  %d fully covered", coveredCount)) + "\n")
	output.WriteString(normalStyle.Render(fmt.Sprintf("  %d partially covered", partialCount)) + "\n")
	output.WriteString(normalStyle.Render(fmt.Sprintf("  %d not covered", uncoveredCount)) + "\n\n")

	// Calculate max function name length for alignment
	maxNameLen := 0
	for _, fc := range result.FunctionCoverages {
		if len(fc.FunctionName) > maxNameLen {
			maxNameLen = len(fc.FunctionName)
		}
	}
	nameColWidth := maxNameLen + 4 // Add 4 characters padding

	// Show fully covered functions (first - good news first!)
	if coveredCount > 0 {
		output.WriteString(normalStyle.Render("Fully Covered:") + "\n")
		output.WriteString(separatorStyle.Render("-------------------------------------------") + "\n")

		for _, fc := range result.FunctionCoverages {
			if fc.CoveragePercent == 100.0 {
				output.WriteString(fmt.Sprintf("  %s%s%s:%d (%d stmts)\n",
					normalStyle.Render(fc.FunctionName),
					strings.Repeat(" ", nameColWidth-len(fc.FunctionName)),
					normalStyle.Render(fc.FileName),
					fc.Line,
					fc.TotalStmts))
			}
		}
	}

	// Show partially covered functions
	if partialCount > 0 {
		if coveredCount > 0 {
			output.WriteString("\n")
		}
		output.WriteString(normalStyle.Render("Partially Covered - Full Testing Would Add:") + "\n")
		output.WriteString(separatorStyle.Render("-------------------------------------------") + "\n")

		for _, fc := range result.FunctionCoverages {
			if fc.CoveragePercent > 0.0 && fc.CoveragePercent < 100.0 {
				output.WriteString(fmt.Sprintf("  %s%s%+6.1f%%  (%5.1f%% covered)  %s:%d (%d stmts)\n",
					normalStyle.Render(fc.FunctionName),
					strings.Repeat(" ", nameColWidth-len(fc.FunctionName)),
					fc.ImpactPercent,
					fc.CoveragePercent,
					normalStyle.Render(fc.FileName),
					fc.Line,
					fc.UncoveredStmts))
			}
		}
	}

	// Show untested functions (0% coverage) - filtered by impact
	if uncoveredCount > 0 {
		if coveredCount > 0 || partialCount > 0 {
			output.WriteString("\n")
		}
		output.WriteString(normalStyle.Render("Not Covered - Testing Would Add:") + "\n")
		output.WriteString(separatorStyle.Render("-------------------------------------------") + "\n")

		const maxUncoveredToShow = 15
		const minImpactThreshold = 3.0 // Show functions with at least 3% impact

		shown := 0
		skippedLowImpact := 0
		var lowestShownImpact float64 = 999.0  // Track the lowest impact we actually showed
		var highestSkippedImpact float64 = 0.0 // Track the highest impact we skipped

		for _, fc := range result.FunctionCoverages {
			if fc.CoveragePercent == 0.0 {
				// Show top 15 OR functions with >= 3% impact
				if shown < maxUncoveredToShow || fc.ImpactPercent >= minImpactThreshold {
					output.WriteString(fmt.Sprintf("  %s%s%+6.1f%%  %s:%d (%d stmts)\n",
						normalStyle.Render(fc.FunctionName),
						strings.Repeat(" ", nameColWidth-len(fc.FunctionName)),
						fc.ImpactPercent,
						normalStyle.Render(fc.FileName),
						fc.Line,
						fc.UncoveredStmts))
					if fc.ImpactPercent < lowestShownImpact {
						lowestShownImpact = fc.ImpactPercent
					}
					shown++
				} else {
					if fc.ImpactPercent > highestSkippedImpact {
						highestSkippedImpact = fc.ImpactPercent
					}
					skippedLowImpact++
				}
			}
		}

		if skippedLowImpact > 0 {
			dimStyle := lipgloss.NewStyle().Foreground(theme.HelpColor)
			output.WriteString("\n" + dimStyle.Render(fmt.Sprintf("  ... and %d more with <= %.1f%% impact", skippedLowImpact, highestSkippedImpact)) + "\n")
		}
	}

	output.WriteString("\n")
	// Make ESC message highly visible with red asterisks and white text
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	whiteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	output.WriteString(redStyle.Render("** ") + whiteStyle.Render("Press ESC to return") + redStyle.Render(" **") + "\n")

	return output.String()
}

// renderProgressBar creates an ASCII progress bar
func renderProgressBar(percentage float64, theme Theme) string {
	barWidth := 30
	filled := int(percentage / 100.0 * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	passStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	normalStyle := lipgloss.NewStyle().Foreground(theme.NormalFg)
	metricStyle := lipgloss.NewStyle().Foreground(theme.MenuActiveFg)

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += "]"

	styledBar := passStyle.Render(bar[:filled+1]) + normalStyle.Render(bar[filled+1:])
	return styledBar + " " + metricStyle.Render(fmt.Sprintf("%.1f%%", percentage))
}

// FormatTestResult formats a test result for display with theme styling
func FormatTestResult(result *PackageTestResult, theme Theme) string {
	var output strings.Builder

	// Styles
	separatorStyle := lipgloss.NewStyle().Foreground(theme.TreeSymbolColor)
	normalStyle := lipgloss.NewStyle().Foreground(theme.NormalFg)
	metricStyle := lipgloss.NewStyle().Foreground(theme.MenuActiveFg)
	passStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

	separator := separatorStyle.Render("========================================")
	output.WriteString(separator + "\n")
	output.WriteString(normalStyle.Render(fmt.Sprintf("Package: %s", result.PackagePath)) + "\n")
	output.WriteString(separator + "\n")

	if result.Status == "RUNNING" {
		output.WriteString(normalStyle.Render("Status: Running tests...") + "\n\n")
		return output.String()
	}

	if result.Status == "NOT_RUN" {
		output.WriteString(normalStyle.Render("Status: Not run") + "\n\n")
		output.WriteString(normalStyle.Render("Press Enter to run tests for this package") + "\n")
		return output.String()
	}

	// Status summary with color based on pass/fail
	var statusStyled string
	if result.Status == "PASS" {
		statusStyled = passStyle.Render("PASS")
	} else {
		statusStyled = failStyle.Render("FAIL")
	}
	output.WriteString(normalStyle.Render("Status: ") + statusStyled +
		normalStyle.Render(fmt.Sprintf(" (%d/%d passed)", result.PassedTests, result.TotalTests)) + "\n")

	output.WriteString(normalStyle.Render("Coverage: ") +
		metricStyle.Render(fmt.Sprintf("%.1f%%", result.Coverage)) + "\n")

	// Show time or "(cached)" if no duration
	if result.Duration > 0 {
		output.WriteString(normalStyle.Render("Total Time: ") +
			metricStyle.Render(formatDuration(result.Duration)) + "\n\n")
	} else {
		output.WriteString(normalStyle.Render("Total Time: ") +
			metricStyle.Render("(cached)") + "\n\n")
	}

	// FAILURES FIRST - Show failed tests prominently at the top
	if result.FailedTests > 0 {
		failureBox := separatorStyle.Render("┌─────────────────────────────────────────┐\n│            FAILURES                     │\n└─────────────────────────────────────────┘")
		output.WriteString(failureBox + "\n\n")

		for _, test := range result.Tests {
			if test.Status == "FAIL" {
				output.WriteString(failStyle.Render("[ FAIL ] ") +
					normalStyle.Render(test.Name) +
					metricStyle.Render(fmt.Sprintf(" (%s)", formatDuration(test.Duration))) + "\n")
				output.WriteString(separatorStyle.Render("-------------------------------------------") + "\n")

				if test.Output != "" {
					// Format failure output more clearly
					lines := strings.Split(strings.TrimSpace(test.Output), "\n")
					for _, line := range lines {
						trimmed := strings.TrimSpace(line)
						if trimmed == "" {
							continue
						}

						// Highlight common assertion patterns
						if strings.Contains(trimmed, "Error:") || strings.Contains(trimmed, "error:") {
							output.WriteString(failStyle.Render("  >>> ") + normalStyle.Render(trimmed) + "\n")
						} else if strings.Contains(trimmed, "want") || strings.Contains(trimmed, "got") {
							output.WriteString(normalStyle.Render("      "+trimmed) + "\n")
						} else if strings.Contains(trimmed, ".go:") {
							// Likely a file location
							output.WriteString(normalStyle.Render("  at: "+trimmed) + "\n")
						} else {
							output.WriteString(normalStyle.Render("      "+trimmed) + "\n")
						}
					}
				} else {
					output.WriteString(normalStyle.Render("  (No failure details captured)") + "\n")
				}
				output.WriteString("\n")
			}
		}
	}

	// Summary of all tests
	if len(result.Tests) > 0 {
		output.WriteString(normalStyle.Render("All Tests:") + "\n")
		output.WriteString(separatorStyle.Render("-------------------------------------------") + "\n")

		// Group by status: FAIL, PASS, SKIP
		var failed, passed, skipped []TestResult
		for _, test := range result.Tests {
			switch test.Status {
			case "FAIL":
				failed = append(failed, test)
			case "PASS":
				passed = append(passed, test)
			case "SKIP":
				skipped = append(skipped, test)
			}
		}

		// Show failures (just names, details shown above)
		for _, test := range failed {
			output.WriteString(failStyle.Render("  [FAIL] ") +
				normalStyle.Render(fmt.Sprintf("%-45s", test.Name)) +
				metricStyle.Render(fmt.Sprintf(" %8s", formatDuration(test.Duration))) + "\n")
		}

		// Show passes
		for _, test := range passed {
			output.WriteString(passStyle.Render("  [PASS] ") +
				normalStyle.Render(fmt.Sprintf("%-45s", test.Name)) +
				metricStyle.Render(fmt.Sprintf(" %8s", formatDuration(test.Duration))) + "\n")
		}

		// Show skipped
		for _, test := range skipped {
			output.WriteString(normalStyle.Render(fmt.Sprintf("  [SKIP] %-45s (skipped)", test.Name)) + "\n")
		}

	} else {
		// No tests found - show the actual output to help debug
		output.WriteString("\n")
		noTestsBox := separatorStyle.Render("┌─────────────────────────────────────────┐\n│       NO TESTS PARSED                   │\n└─────────────────────────────────────────┘")
		output.WriteString(noTestsBox + "\n\n")

		if result.Status == "FAIL" {
			output.WriteString(normalStyle.Render("The test command failed. Raw output:") + "\n")
			output.WriteString(separatorStyle.Render("-------------------------------------------") + "\n")
			// Show first 30 lines of output to help debug
			lines := strings.Split(result.FullOutput, "\n")
			maxLines := 30
			if len(lines) > maxLines {
				for i := 0; i < maxLines; i++ {
					output.WriteString(normalStyle.Render(lines[i]) + "\n")
				}
				output.WriteString(normalStyle.Render(fmt.Sprintf("\n... (%d more lines, scroll down to see all)", len(lines)-maxLines)) + "\n")
			} else {
				output.WriteString(normalStyle.Render(result.FullOutput))
			}
		} else {
			output.WriteString(normalStyle.Render("No test functions found in this package,\nor the test output format was not recognized.") + "\n")
		}
	}

	// Per-file coverage breakdown
	if len(result.FileCoverages) > 0 {
		output.WriteString("\n")
		output.WriteString(normalStyle.Render("Per-File Coverage:") + "\n")
		output.WriteString(separatorStyle.Render("-------------------------------------------") + "\n")

		for _, fc := range result.FileCoverages {
			// Color-code based on coverage
			// [!!] = poor (<50%), [**] = medium (50-80%), [++] = good (>=80%)
			var indicatorStyled string
			var coverageColor lipgloss.Style
			if fc.CoveragePercent >= 80.0 {
				coverageColor = passStyle
				indicatorStyled = passStyle.Render("[++]") // Good coverage
			} else if fc.CoveragePercent >= 50.0 {
				coverageColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00")) // Yellow
				indicatorStyled = coverageColor.Render("[**]")                            // Medium coverage
			} else {
				coverageColor = failStyle
				indicatorStyled = failStyle.Render("[!!]") // Poor coverage
			}

			output.WriteString(fmt.Sprintf("  %s %s %s %s\n",
				indicatorStyled,
				normalStyle.Render(fmt.Sprintf("%-40s", fc.FileName)),
				metricStyle.Render(fmt.Sprintf("%6.1f%%", fc.CoveragePercent)),
				normalStyle.Render(fmt.Sprintf("(%d/%d stmts)", fc.CoveredLines, fc.TotalLines))))
		}
		output.WriteString("\n")
		output.WriteString(normalStyle.Render("Legend: ") +
			passStyle.Render("[++]") + normalStyle.Render(" >=80%  ") +
			lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff00")).Render("[**]") + normalStyle.Render(" 50-80%  ") +
			failStyle.Render("[!!]") + normalStyle.Render(" <50%") + "\n")
	}

	output.WriteString("\n")
	// Make ESC message highly visible with red asterisks and white text
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	whiteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	output.WriteString(redStyle.Render("** ") + whiteStyle.Render("Press ESC to return") + redStyle.Render(" **") + "\n")

	return output.String()
}
