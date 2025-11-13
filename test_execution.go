package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunTests executes tests for a specific package
// packageDir is the directory containing the test files
// mode specifies which tests to run (unit, integration, or all)
func RunTests(packageDir string, packageName string, mode testMode) (*PackageTestResult, error) {
	LogInfo("Running tests",
		"package", packageName,
		"directory", packageDir,
		"mode", mode,
	)

	// For "All" mode, run tests twice and compare to identify integration tests
	if mode == testModeAll {
		return runAllTests(packageDir, packageName)
	}

	// Single run for unit or integration mode
	return runSingleTestMode(packageDir, packageName, mode)
}

// runSingleTestMode runs tests once with the specified mode
func runSingleTestMode(packageDir string, packageName string, mode testMode) (*PackageTestResult, error) {
	result := &PackageTestResult{
		PackagePath:       packageName,
		Status:            "RUNNING",
		Tests:             []TestResult{},
		FileCoverages:     []FileCoverage{},
		FunctionCoverages: []FunctionCoverage{},
	}

	// Create temp file for coverage output
	tempFile, err := os.CreateTemp("", "gapistotle-coverage-*.out")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp coverage file: %w", err)
	}
	coverageFile := tempFile.Name()
	tempFile.Close()
	defer os.Remove(coverageFile) // Clean up after parsing

	// Build command args based on test mode
	args := []string{"test", "-json", "-cover", "-coverprofile=" + coverageFile, "-count=1"}

	// Add tags based on mode
	testType := "unit"
	if mode == testModeIntegration {
		args = append(args, "-tags=integration")
		testType = "integration"
	}

	args = append(args, ".")

	// Run go test in the directory containing the tests
	cmd := exec.Command("go", args...)
	cmd.Dir = packageDir

	// Capture combined output
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	result.FullOutput = outputStr

	// Parse the output and tag tests with type
	parseTestOutput(result, outputStr)
	for i := range result.Tests {
		result.Tests[i].TestType = testType
	}

	// Parse coverage profile if it exists
	if _, statErr := os.Stat(coverageFile); statErr == nil {
		parseCoverageProfile(result, coverageFile)
		parseFunctionCoverage(result, coverageFile, packageDir)
		LogDebug("Parsed function coverage",
			"package", packageDir,
			"function_count", len(result.FunctionCoverages),
		)
	} else {
		LogWarn("Coverage file not found",
			"coverage_file", coverageFile,
			"error", statErr,
		)
	}

	// Determine overall status
	if err != nil {
		result.Status = "FAIL"
	} else {
		result.Status = "PASS"
	}

	LogInfo("Test execution complete",
		"package", packageName,
		"status", result.Status,
		"coverage", result.Coverage,
		"total_tests", result.TotalTests,
		"passed", result.PassedTests,
		"failed", result.FailedTests,
		"duration", result.Duration.String(),
	)

	return result, nil
}

// runAllTests runs both unit and integration tests and combines results
func runAllTests(packageDir string, packageName string) (*PackageTestResult, error) {
	// Run unit tests first (no tags)
	unitResult, unitErr := runSingleTestMode(packageDir, packageName, testModeUnit)
	if unitErr != nil {
		return nil, unitErr
	}

	// Build map of unit test names for comparison
	unitTestNames := make(map[string]bool)
	for _, test := range unitResult.Tests {
		unitTestNames[test.Name] = true
	}

	// Run with integration tags to get all tests
	allResult, allErr := runSingleTestMode(packageDir, packageName, testModeIntegration)
	if allErr != nil {
		return nil, allErr
	}

	// Tag tests: if it's in unitTestNames, it's a unit test; otherwise it's integration
	for i := range allResult.Tests {
		if unitTestNames[allResult.Tests[i].Name] {
			allResult.Tests[i].TestType = "unit"
		} else {
			allResult.Tests[i].TestType = "integration"
		}
	}

	LogInfo("All tests execution complete (unit + integration)",
		"package", packageName,
		"total_tests", allResult.TotalTests,
		"unit_tests", len(unitTestNames),
		"integration_tests", allResult.TotalTests-len(unitTestNames),
	)

	return allResult, nil
}
