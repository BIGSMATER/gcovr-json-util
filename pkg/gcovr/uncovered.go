package gcovr

import (
	"fmt"
	"sort"
)

// FindUncoveredLines analyzes a gcovr report and returns all uncovered lines
// grouped by file and function
func FindUncoveredLines(report *GcovrReport) (*UncoveredReport, error) {
	result := &UncoveredReport{
		Files: make([]FileUncovered, 0),
	}

	// Map structure: file -> function -> uncovered line numbers
	uncoveredMap := make(map[string]map[string][]int)

	// Map to store function metadata (demangled names)
	funcMetadata := make(map[string]map[string]string) // file -> funcName -> demangledName

	// First pass: collect uncovered lines
	for _, file := range report.Files {
		uncoveredMap[file.FilePath] = make(map[string][]int)
		funcMetadata[file.FilePath] = make(map[string]string)

		// Store function metadata
		for _, fn := range file.Functions {
			funcMetadata[file.FilePath][fn.Name] = fn.DemangledName
		}

		// Find uncovered lines
		for _, line := range file.Lines {
			if line.Count == 0 {
				if _, exists := uncoveredMap[file.FilePath][line.FunctionName]; !exists {
					uncoveredMap[file.FilePath][line.FunctionName] = make([]int, 0)
				}
				uncoveredMap[file.FilePath][line.FunctionName] = append(
					uncoveredMap[file.FilePath][line.FunctionName],
					line.LineNumber,
				)
			}
		}
	}

	// Sort file paths for consistent output
	filePaths := make([]string, 0, len(uncoveredMap))
	for filePath := range uncoveredMap {
		if len(uncoveredMap[filePath]) > 0 {
			filePaths = append(filePaths, filePath)
		}
	}
	sort.Strings(filePaths)

	// Second pass: build FileUncovered structs with FunctionUncovered
	for _, filePath := range filePaths {
		funcUncovered := uncoveredMap[filePath]

		fileResult := FileUncovered{
			FilePath:           filePath,
			UncoveredFunctions: make([]FunctionUncovered, 0),
		}

		// Get the file object for line stats
		var fileObj *File
		for i := range report.Files {
			if report.Files[i].FilePath == filePath {
				fileObj = &report.Files[i]
				break
			}
		}

		if fileObj == nil {
			continue
		}

		for funcName, uncoveredLines := range funcUncovered {
			if len(uncoveredLines) == 0 {
				continue
			}

			// Calculate total lines and covered lines for this function
			totalLines := 0
			coveredLines := 0

			for _, line := range fileObj.Lines {
				if line.FunctionName == funcName {
					totalLines++
					if line.Count > 0 {
						coveredLines++
					}
				}
			}

			// Get demangled name
			demangledName := funcMetadata[filePath][funcName]
			if demangledName == "" {
				demangledName = funcName
			}

			// Sort line numbers for consistent output
			sort.Ints(uncoveredLines)

			fileResult.UncoveredFunctions = append(fileResult.UncoveredFunctions, FunctionUncovered{
				FunctionName:         funcName,
				DemangledName:        demangledName,
				UncoveredLineNumbers: uncoveredLines,
				TotalLines:           totalLines,
				CoveredLines:         coveredLines,
			})
		}

		if len(fileResult.UncoveredFunctions) > 0 {
			result.Files = append(result.Files, fileResult)
		}
	}

	return result, nil
}

// FormatUncoveredReport formats the uncovered lines report as a human-readable string
func FormatUncoveredReport(report *UncoveredReport) string {
	if len(report.Files) == 0 {
		return "No uncovered lines found. All lines have coverage!\n"
	}

	// Calculate total statistics
	totalFunctions := 0
	totalUncoveredLines := 0
	for _, file := range report.Files {
		totalFunctions += len(file.UncoveredFunctions)
		for _, fn := range file.UncoveredFunctions {
			totalUncoveredLines += len(fn.UncoveredLineNumbers)
		}
	}

	// Build output
	result := fmt.Sprintf("Uncovered Lines Report\n")
	result += fmt.Sprintf("======================\n\n")

	result += fmt.Sprintf("Found %d function(s) with uncovered lines (%d total uncovered lines):\n\n",
		totalFunctions, totalUncoveredLines)

	funcIdx := 1
	for _, file := range report.Files {
		for _, fn := range file.UncoveredFunctions {
			coveragePercent := 0.0
			if fn.TotalLines > 0 {
				coveragePercent = float64(fn.CoveredLines) * 100.0 / float64(fn.TotalLines)
			}

			result += fmt.Sprintf("%d. File: %s\n", funcIdx, file.FilePath)
			result += fmt.Sprintf("   Function: %s\n", fn.DemangledName)
			result += fmt.Sprintf("   Coverage: %d/%d lines (%.1f%%)\n",
				fn.CoveredLines, fn.TotalLines, coveragePercent)
			result += fmt.Sprintf("   Uncovered Lines (%d): %v\n\n",
				len(fn.UncoveredLineNumbers), fn.UncoveredLineNumbers)

			funcIdx++
		}
	}

	return result
}
