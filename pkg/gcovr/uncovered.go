package gcovr

import (
	"fmt"
	"sort"
)

// FindUncoveredLines analyzes a gcovr report and returns all uncovered lines
// grouped by file and function
func FindUncoveredLines(report *GcovrReport) (*UncoveredReport, error) {
	result := &UncoveredReport{
		UncoveredFunctions: make([]FunctionUncovered, 0),
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

	// Second pass: build FunctionUncovered structs with complete stats
	for _, file := range report.Files {
		funcUncovered := uncoveredMap[file.FilePath]

		for funcName, uncoveredLines := range funcUncovered {
			if len(uncoveredLines) == 0 {
				continue
			}

			// Calculate total lines and covered lines for this function
			totalLines := 0
			coveredLines := 0

			for _, line := range file.Lines {
				if line.FunctionName == funcName {
					totalLines++
					if line.Count > 0 {
						coveredLines++
					}
				}
			}

			// Get demangled name
			demangledName := funcMetadata[file.FilePath][funcName]
			if demangledName == "" {
				demangledName = funcName
			}

			// Sort line numbers for consistent output
			sort.Ints(uncoveredLines)

			result.UncoveredFunctions = append(result.UncoveredFunctions, FunctionUncovered{
				File:                 file.FilePath,
				FunctionName:         funcName,
				DemangledName:        demangledName,
				UncoveredLineNumbers: uncoveredLines,
				TotalLines:           totalLines,
				CoveredLines:         coveredLines,
			})
		}
	}

	return result, nil
}

// FormatUncoveredReport formats the uncovered lines report as a human-readable string
func FormatUncoveredReport(report *UncoveredReport) string {
	if len(report.UncoveredFunctions) == 0 {
		return "No uncovered lines found. All lines have coverage!\n"
	}

	// Group by file
	fileGroups := make(map[string][]FunctionUncovered)
	for _, fn := range report.UncoveredFunctions {
		fileGroups[fn.File] = append(fileGroups[fn.File], fn)
	}

	// Build output
	result := fmt.Sprintf("Uncovered Lines Report\n")
	result += fmt.Sprintf("======================\n\n")

	totalUncoveredLines := 0
	for _, fn := range report.UncoveredFunctions {
		totalUncoveredLines += len(fn.UncoveredLineNumbers)
	}

	result += fmt.Sprintf("Found %d function(s) with uncovered lines (%d total uncovered lines):\n\n",
		len(report.UncoveredFunctions), totalUncoveredLines)

	// Sort files for consistent output
	files := make([]string, 0, len(fileGroups))
	for file := range fileGroups {
		files = append(files, file)
	}
	sort.Strings(files)

	funcIdx := 1
	for _, filePath := range files {
		functions := fileGroups[filePath]

		for _, fn := range functions {
			coveragePercent := 0.0
			if fn.TotalLines > 0 {
				coveragePercent = float64(fn.CoveredLines) * 100.0 / float64(fn.TotalLines)
			}

			result += fmt.Sprintf("%d. File: %s\n", funcIdx, fn.File)
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
