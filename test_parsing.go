package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// parseTestOutput parses go test -json output to extract test results
func parseTestOutput(result *PackageTestResult, output string) {
	scanner := bufio.NewScanner(strings.NewReader(output))

	// Track test outputs by test name
	testOutputs := make(map[string]*strings.Builder)

	// Regex for coverage in output lines
	coverageRegex := regexp.MustCompile(`coverage: (\d+\.\d+)% of statements`)

	for scanner.Scan() {
		line := scanner.Text()

		// Parse JSON event
		var event TestEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			// Skip non-JSON lines (shouldn't happen with -json, but be safe)
			continue
		}

		switch event.Action {
		case "run":
			// Test started - initialize output collector
			if event.Test != "" {
				testOutputs[event.Test] = &strings.Builder{}
			}

		case "output":
			// Collect output for test or check for coverage
			if event.Test != "" {
				// Test-specific output
				if builder, ok := testOutputs[event.Test]; ok {
					builder.WriteString(event.Output)
				}
			} else {
				// Package-level output - check for coverage
				if matches := coverageRegex.FindStringSubmatch(event.Output); matches != nil {
					coverage, _ := strconv.ParseFloat(matches[1], 64)
					result.Coverage = coverage
				}
			}

		case "pass", "fail", "skip":
			if event.Test != "" {
				// Individual test completed
				status := strings.ToUpper(event.Action)
				testResult := TestResult{
					Name:     event.Test,
					Status:   status,
					Duration: time.Duration(event.Elapsed * float64(time.Second)),
				}

				// Attach collected output
				if builder, ok := testOutputs[event.Test]; ok {
					testResult.Output = builder.String()
					delete(testOutputs, event.Test) // Clean up
				}

				result.Tests = append(result.Tests, testResult)
				result.TotalTests++

				switch status {
				case "PASS":
					result.PassedTests++
				case "FAIL":
					result.FailedTests++
				case "SKIP":
					result.SkippedTests++
				}
			} else {
				// Package completed - record total duration
				result.Duration = time.Duration(event.Elapsed * float64(time.Second))
			}
		}
	}
}

// parseCoverageProfile parses a go test coverage profile to extract per-file coverage
func parseCoverageProfile(result *PackageTestResult, profilePath string) {
	file, err := os.Open(profilePath)
	if err != nil {
		return
	}
	defer file.Close()

	// Map of filename -> {covered statements, total statements}
	fileCoverage := make(map[string][2]int)

	scanner := bufio.NewScanner(file)
	// Skip first line (mode: set/count/atomic)
	if scanner.Scan() {
		// Skip mode line
	}

	// Parse coverage lines
	// Format: filename:startline.startcol,endline.endcol numstatements covered
	lineRegex := regexp.MustCompile(`^(.+):(\d+\.\d+),(\d+\.\d+)\s+(\d+)\s+(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := lineRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		filename := filepath.Base(matches[1]) // Just get the filename, not full path
		numStatements, _ := strconv.Atoi(matches[4])
		covered, _ := strconv.Atoi(matches[5])

		stats := fileCoverage[filename]
		stats[1] += numStatements // Total statements
		if covered > 0 {
			stats[0] += numStatements // Covered statements
		}
		fileCoverage[filename] = stats
	}

	// Convert map to slice and calculate percentages
	for filename, stats := range fileCoverage {
		coveredLines := stats[0]
		totalLines := stats[1]
		var percentage float64
		if totalLines > 0 {
			percentage = (float64(coveredLines) / float64(totalLines)) * 100.0
		}

		result.FileCoverages = append(result.FileCoverages, FileCoverage{
			FileName:        filename,
			CoveredLines:    coveredLines,
			TotalLines:      totalLines,
			CoveragePercent: percentage,
		})
	}

	// Sort by coverage percentage (worst first for easy identification)
	sort.Slice(result.FileCoverages, func(i, j int) bool {
		return result.FileCoverages[i].CoveragePercent < result.FileCoverages[j].CoveragePercent
	})
}

// parseFunctionCoverage parses function-level coverage using go tool cover
func parseFunctionCoverage(result *PackageTestResult, profilePath string, packageDir string) {
	// Step 1: Parse go tool cover -func to get function definitions
	cmd := exec.Command("go", "tool", "cover", "-func="+profilePath)
	cmd.Dir = packageDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return // Silently fail if we can't get function coverage
	}

	// Parse function definitions
	// Format: filename.go:line:	FunctionName	percentage%
	type funcDef struct {
		name            string
		line            int
		coveragePercent float64
	}

	funcMap := make(map[string][]funcDef) // filename -> functions
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	lineRegex := regexp.MustCompile(`^(.+):(\d+):\s+(\S+)\s+([\d.]+)%`)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "total:") {
			continue
		}

		matches := lineRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		filename := matches[1]
		lineNum, _ := strconv.Atoi(matches[2])
		functionName := matches[3]
		coveragePercent, _ := strconv.ParseFloat(matches[4], 64)

		funcMap[filename] = append(funcMap[filename], funcDef{
			name:            functionName,
			line:            lineNum,
			coveragePercent: coveragePercent,
		})
	}

	// Sort functions by line number for each file
	for filename := range funcMap {
		sort.Slice(funcMap[filename], func(i, j int) bool {
			return funcMap[filename][i].line < funcMap[filename][j].line
		})
	}

	// Step 2: Parse coverage profile to get statement counts per block
	profileData, err := os.ReadFile(profilePath)
	if err != nil {
		return
	}

	blocks := make(map[string][]coverageBlock) // filename -> blocks
	profileScanner := bufio.NewScanner(strings.NewReader(string(profileData)))
	// Format: filename:startLine.col,endLine.col numStmts covered
	blockRegex := regexp.MustCompile(`^(.+):(\d+)\.\d+,(\d+)\.\d+\s+(\d+)\s+(\d+)`)

	for profileScanner.Scan() {
		line := profileScanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}

		matches := blockRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		filename := matches[1]
		startLine, _ := strconv.Atoi(matches[2])
		endLine, _ := strconv.Atoi(matches[3])
		numStmts, _ := strconv.Atoi(matches[4])
		coveredCount, _ := strconv.Atoi(matches[5])

		blocks[filename] = append(blocks[filename], coverageBlock{
			startLine: startLine,
			endLine:   endLine,
			numStmt:   numStmts,
			count:     coveredCount,
		})
	}

	// Calculate total statements across all blocks for impact calculation
	var totalPackageStmts int
	for _, fileBlocks := range blocks {
		for _, block := range fileBlocks {
			totalPackageStmts += block.numStmt
		}
	}

	// Step 3: Map blocks to functions and calculate uncovered statements
	functionStats := make(map[string]*FunctionCoverage)

	for filename, funcs := range funcMap {
		fileBlocks := blocks[filename]

		for _, fn := range funcs {
			// Find the next function's line to determine this function's range
			endLine := math.MaxInt
			for _, otherFn := range funcs {
				if otherFn.line > fn.line && otherFn.line < endLine {
					endLine = otherFn.line
				}
			}

			// Sum up statements for blocks belonging to this function
			var totalStmts, uncoveredStmts int
			for _, block := range fileBlocks {
				if block.startLine >= fn.line && block.startLine < endLine {
					totalStmts += block.numStmt
					if block.count == 0 {
						uncoveredStmts += block.numStmt
					}
				}
			}

			// Calculate impact: how much would overall coverage increase if this function was fully tested
			var impact float64
			if totalPackageStmts > 0 {
				impact = (float64(uncoveredStmts) / float64(totalPackageStmts)) * 100.0
			}

			key := fmt.Sprintf("%s:%s", filename, fn.name)
			functionStats[key] = &FunctionCoverage{
				FunctionName:    fn.name,
				FileName:        filepath.Base(filename),
				Line:            fn.line,
				CoveragePercent: fn.coveragePercent,
				TotalStmts:      totalStmts,
				UncoveredStmts:  uncoveredStmts,
				ImpactPercent:   impact,
			}
		}
	}

	// Convert map to slice
	for _, fc := range functionStats {
		result.FunctionCoverages = append(result.FunctionCoverages, *fc)
	}

	// Sort by most uncovered statements first (biggest impact)
	sort.Slice(result.FunctionCoverages, func(i, j int) bool {
		return result.FunctionCoverages[i].UncoveredStmts > result.FunctionCoverages[j].UncoveredStmts
	})
}
