package main

import (
	"fmt"
	"time"
)

// TestEvent represents a single event from go test -json
type TestEvent struct {
	Time    time.Time
	Action  string  // "start", "run", "pass", "fail", "skip", "output"
	Package string
	Test    string  // Present for test-specific events
	Elapsed float64 // Duration in seconds
	Output  string  // Present for "output" events
}

// formatDuration formats a duration with adaptive units for better precision display
// 0: Below measurement threshold (< 10ms)
// < 1s: milliseconds (ms)
// >= 1s: seconds (s)
func formatDuration(d time.Duration) string {
	if d == 0 {
		// go test -v only reports to 0.01s precision, so 0.00s means < 5ms (rounds to 0)
		return "< 10ms"
	} else if d < time.Second {
		// Show in milliseconds
		ms := float64(d.Microseconds()) / 1000.0
		if ms < 10 {
			return fmt.Sprintf("%.2fms", ms)
		} else if ms < 100 {
			return fmt.Sprintf("%.1fms", ms)
		}
		return fmt.Sprintf("%.0fms", ms)
	}
	// Show in seconds
	s := d.Seconds()
	if s < 10 {
		return fmt.Sprintf("%.3fs", s)
	}
	return fmt.Sprintf("%.2fs", s)
}

// TestResult represents the result of a single test
type TestResult struct {
	Name     string
	Status   string // "PASS", "FAIL", "SKIP"
	Duration time.Duration
	Output   string // Detailed output for failed tests
}

// FileCoverage represents coverage for a single file
type FileCoverage struct {
	FileName       string
	CoveredLines   int
	TotalLines     int
	CoveragePercent float64
}

// FunctionCoverage represents coverage for a single function
type FunctionCoverage struct {
	FunctionName    string
	FileName        string
	Line            int
	CoveragePercent float64
	TotalStmts      int     // Total statements in this function
	UncoveredStmts  int     // Number of uncovered statements in this function
	ImpactPercent   float64 // How much package coverage would increase if this function was fully tested
}

// PackageTestResult represents test results for an entire package
type PackageTestResult struct {
	PackagePath       string
	Status            string // "PASS", "FAIL", "RUNNING", "NOT_RUN"
	Coverage          float64
	TotalTests        int
	PassedTests       int
	FailedTests       int
	SkippedTests      int
	Duration          time.Duration
	Tests             []TestResult
	FileCoverages     []FileCoverage     // Per-file coverage details
	FunctionCoverages []FunctionCoverage // Per-function coverage details
	FullOutput        string
}

// coverageBlock represents a coverage block from the coverage profile
type coverageBlock struct {
	startLine int
	endLine   int
	numStmt   int
	count     int
}
