package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunTests executes tests for a specific package
// packageDir is the directory containing the test files
func RunTests(packageDir string, packageName string) (*PackageTestResult, error) {
	LogInfo("Running tests",
		"package", packageName,
		"directory", packageDir,
	)

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

	// Run go test in the directory containing the tests
	// -json provides structured output for easier parsing
	// -count=1 disables test caching to ensure fresh results
	// -coverprofile generates per-file coverage data (use absolute path)
	cmd := exec.Command("go", "test", "-json", "-cover", "-coverprofile="+coverageFile, "-count=1", ".")
	cmd.Dir = packageDir

	// Capture combined output
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	result.FullOutput = outputStr

	// Parse the output
	parseTestOutput(result, outputStr)

	// Parse coverage profile if it exists
	if _, statErr := os.Stat(coverageFile); statErr == nil {
		parseCoverageProfile(result, coverageFile)
		// Pass the full temp file path
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
