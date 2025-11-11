package gcovr

import (
	"fmt"
)

// ComputeCoverageIncrease calculates coverage increases from base to new report
// It returns a report containing functions with increased line coverage
func ComputeCoverageIncrease(baseReport, newReport *GcovrReport) (*CoverageIncreaseReport, error) {
	result := &CoverageIncreaseReport{
		Increases: make([]FunctionCoverageIncrease, 0),
	}

	// Create maps for quick lookup
	baseFileMap := make(map[string]*File)
	for i := range baseReport.Files {
		baseFileMap[baseReport.Files[i].FilePath] = &baseReport.Files[i]
	}

	// Iterate through new report files
	for _, newFile := range newReport.Files {
		baseFile, exists := baseFileMap[newFile.FilePath]
		if !exists {
			// Entirely new file, all covered lines are increases
			result.Increases = append(result.Increases, processNewFile(&newFile)...)
			continue
		}

		// Compare functions in the same file
		increases := compareFunctions(baseFile, &newFile)
		result.Increases = append(result.Increases, increases...)
	}

	return result, nil
}

// processNewFile processes a file that exists in new but not in base
func processNewFile(file *File) []FunctionCoverageIncrease {
	increases := make([]FunctionCoverageIncrease, 0)

	// Group lines by function
	funcLines := make(map[string][]int)
	funcDemangledName := make(map[string]string)

	for _, line := range file.Lines {
		if line.Count > 0 {
			funcLines[line.FunctionName] = append(funcLines[line.FunctionName], line.LineNumber)
		}
	}

	// Get demangled names from function definitions
	for _, fn := range file.Functions {
		funcDemangledName[fn.Name] = fn.DemangledName
	}

	// Create increase records
	for funcName, lineNumbers := range funcLines {
		demangledName := funcDemangledName[funcName]
		if demangledName == "" {
			demangledName = funcName
		}

		totalLines := getTotalFunctionLines(file, funcName)

		increases = append(increases, FunctionCoverageIncrease{
			File:                 file.FilePath,
			FunctionName:         funcName,
			DemangledName:        demangledName,
			LinesIncreased:       len(lineNumbers),
			TotalLines:           totalLines,
			IncreasedLineNumbers: lineNumbers,
			OldCoveredLines:      0, // No old coverage for new file
			NewCoveredLines:      len(lineNumbers),
		})
	}

	return increases
}

// compareFunctions compares functions between base and new file
func compareFunctions(baseFile, newFile *File) []FunctionCoverageIncrease {
	increases := make([]FunctionCoverageIncrease, 0)

	// Create line coverage maps: function -> line_number -> count
	baseCoverage := buildLineCoverageMap(baseFile)
	newCoverage := buildLineCoverageMap(newFile)

	// Get function demangled names
	funcNames := buildFunctionNameMap(newFile)

	// Find increased coverage
	for funcName, newLines := range newCoverage {
		baseLines, exists := baseCoverage[funcName]

		increasedLines := make([]int, 0)
		oldCoveredCount := 0
		newCoveredCount := 0

		for lineNum, newCount := range newLines {
			baseCount := 0
			if exists {
				baseCount = baseLines[lineNum]
			}

			// Count old coverage
			if baseCount > 0 {
				oldCoveredCount++
			}

			// Count new coverage
			if newCount > 0 {
				newCoveredCount++
			}

			// Coverage increased: was 0, now > 0
			if baseCount == 0 && newCount > 0 {
				increasedLines = append(increasedLines, lineNum)
			}
		}

		if len(increasedLines) > 0 {
			demangledName := funcNames[funcName]
			if demangledName == "" {
				demangledName = funcName
			}

			totalLines := getTotalFunctionLines(newFile, funcName)

			increases = append(increases, FunctionCoverageIncrease{
				File:                 newFile.FilePath,
				FunctionName:         funcName,
				DemangledName:        demangledName,
				LinesIncreased:       len(increasedLines),
				TotalLines:           totalLines,
				IncreasedLineNumbers: increasedLines,
				OldCoveredLines:      oldCoveredCount,
				NewCoveredLines:      newCoveredCount,
			})
		}
	}

	return increases
}

// buildLineCoverageMap creates a map of function -> line_number -> count
func buildLineCoverageMap(file *File) map[string]map[int]int {
	result := make(map[string]map[int]int)

	for _, line := range file.Lines {
		if _, exists := result[line.FunctionName]; !exists {
			result[line.FunctionName] = make(map[int]int)
		}
		result[line.FunctionName][line.LineNumber] = line.Count
	}

	return result
}

// buildFunctionNameMap creates a map of mangled -> demangled function names
func buildFunctionNameMap(file *File) map[string]string {
	result := make(map[string]string)

	for _, fn := range file.Functions {
		result[fn.Name] = fn.DemangledName
	}

	return result
}

// getTotalFunctionLines returns the total number of lines in a function
func getTotalFunctionLines(file *File, funcName string) int {
	count := 0
	for _, line := range file.Lines {
		if line.FunctionName == funcName {
			count++
		}
	}
	return count
}

// FormatReport formats the coverage increase report as a human-readable string
func FormatReport(report *CoverageIncreaseReport) string {
	if len(report.Increases) == 0 {
		return "No coverage increases found.\n"
	}

	result := fmt.Sprintf("Coverage Increase Report\n")
	result += fmt.Sprintf("=========================\n\n")
	result += fmt.Sprintf("Found %d function(s) with increased coverage:\n\n", len(report.Increases))

	for i, inc := range report.Increases {
		oldCoveragePercent := 0.0
		newCoveragePercent := 0.0
		if inc.TotalLines > 0 {
			oldCoveragePercent = float64(inc.OldCoveredLines) * 100.0 / float64(inc.TotalLines)
			newCoveragePercent = float64(inc.NewCoveredLines) * 100.0 / float64(inc.TotalLines)
		}

		result += fmt.Sprintf("%d. File: %s\n", i+1, inc.File)
		result += fmt.Sprintf("   Function: %s\n", inc.DemangledName)
		result += fmt.Sprintf("   Old Coverage: %d/%d lines (%.1f%%)\n", inc.OldCoveredLines, inc.TotalLines, oldCoveragePercent)
		result += fmt.Sprintf("   New Coverage: %d/%d lines (%.1f%%)\n", inc.NewCoveredLines, inc.TotalLines, newCoveragePercent)
		result += fmt.Sprintf("   Lines Increased: %d\n", inc.LinesIncreased)
		result += fmt.Sprintf("   Newly Covered Line Numbers: %v\n\n", inc.IncreasedLineNumbers)
	}

	return result
}
