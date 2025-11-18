package gcovr

import (
	"testing"
)

func TestFindUncoveredLines(t *testing.T) {
	tests := []struct {
		name                string
		report              *GcovrReport
		expectedFileCount   int
		expectedFirstFile   string
		expectedFuncCount   int
		expectedUncoveredAt int // expected uncovered lines in first function
	}{
		{
			name: "Single file with uncovered lines",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
							{LineNumber: 2, FunctionName: "foo", Count: 0},
							{LineNumber: 3, FunctionName: "foo", Count: 0},
							{LineNumber: 4, FunctionName: "bar", Count: 1},
							{LineNumber: 5, FunctionName: "bar", Count: 0},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
							{Name: "bar", DemangledName: "bar()"},
						},
					},
				},
			},
			expectedFileCount:   1,
			expectedFirstFile:   "test.cpp",
			expectedFuncCount:   2,
			expectedUncoveredAt: 2,
		},
		{
			name: "Multiple files with uncovered lines",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "a.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "func1", Count: 0},
							{LineNumber: 2, FunctionName: "func1", Count: 0},
						},
						Functions: []Function{
							{Name: "func1", DemangledName: "func1()"},
						},
					},
					{
						FilePath: "b.cpp",
						Lines: []Line{
							{LineNumber: 10, FunctionName: "func2", Count: 0},
						},
						Functions: []Function{
							{Name: "func2", DemangledName: "func2()"},
						},
					},
				},
			},
			expectedFileCount:   2,
			expectedFirstFile:   "a.cpp",
			expectedFuncCount:   1, // functions per file
			expectedUncoveredAt: 2,
		},
		{
			name: "All lines covered",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "covered.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
							{LineNumber: 2, FunctionName: "foo", Count: 5},
							{LineNumber: 3, FunctionName: "foo", Count: 3},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			expectedFileCount:   0,
			expectedFirstFile:   "",
			expectedFuncCount:   0,
			expectedUncoveredAt: 0,
		},
		{
			name: "Empty report",
			report: &GcovrReport{
				Files: []File{},
			},
			expectedFileCount:   0,
			expectedFirstFile:   "",
			expectedFuncCount:   0,
			expectedUncoveredAt: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FindUncoveredLines(tt.report)
			if err != nil {
				t.Fatalf("FindUncoveredLines() error = %v", err)
			}

			if len(result.Files) != tt.expectedFileCount {
				t.Errorf("Expected %d files, got %d", tt.expectedFileCount, len(result.Files))
			}

			if tt.expectedFileCount > 0 {
				firstFile := result.Files[0]
				if firstFile.FilePath != tt.expectedFirstFile {
					t.Errorf("Expected first file %s, got %s", tt.expectedFirstFile, firstFile.FilePath)
				}

				if len(firstFile.UncoveredFunctions) != tt.expectedFuncCount {
					t.Errorf("Expected %d uncovered functions, got %d", tt.expectedFuncCount, len(firstFile.UncoveredFunctions))
				}

				if tt.expectedFuncCount > 0 {
					firstFunc := firstFile.UncoveredFunctions[0]
					if len(firstFunc.UncoveredLineNumbers) != tt.expectedUncoveredAt {
						t.Errorf("Expected %d uncovered lines in first function, got %d",
							tt.expectedUncoveredAt, len(firstFunc.UncoveredLineNumbers))
					}
				}
			}
		})
	}
}

func TestFindUncoveredLines_GroupedByFile(t *testing.T) {
	report := &GcovrReport{
		Files: []File{
			{
				FilePath: "file1.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "foo", Count: 0},
					{LineNumber: 2, FunctionName: "foo", Count: 0},
					{LineNumber: 10, FunctionName: "bar", Count: 0},
				},
				Functions: []Function{
					{Name: "foo", DemangledName: "foo()"},
					{Name: "bar", DemangledName: "bar()"},
				},
			},
			{
				FilePath: "file2.cpp",
				Lines: []Line{
					{LineNumber: 5, FunctionName: "baz", Count: 0},
				},
				Functions: []Function{
					{Name: "baz", DemangledName: "baz()"},
				},
			},
		},
	}

	result, err := FindUncoveredLines(report)
	if err != nil {
		t.Fatalf("FindUncoveredLines() error = %v", err)
	}

	// Verify files are grouped
	if len(result.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(result.Files))
	}

	// Verify file1.cpp has 2 functions
	file1 := result.Files[0]
	if file1.FilePath != "file1.cpp" {
		t.Errorf("Expected first file to be file1.cpp, got %s", file1.FilePath)
	}
	if len(file1.UncoveredFunctions) != 2 {
		t.Errorf("Expected file1.cpp to have 2 uncovered functions, got %d", len(file1.UncoveredFunctions))
	}

	// Verify file2.cpp has 1 function
	file2 := result.Files[1]
	if file2.FilePath != "file2.cpp" {
		t.Errorf("Expected second file to be file2.cpp, got %s", file2.FilePath)
	}
	if len(file2.UncoveredFunctions) != 1 {
		t.Errorf("Expected file2.cpp to have 1 uncovered function, got %d", len(file2.UncoveredFunctions))
	}
}

func TestFindUncoveredLines_Statistics(t *testing.T) {
	report := &GcovrReport{
		Files: []File{
			{
				FilePath: "test.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "foo", Count: 1},
					{LineNumber: 2, FunctionName: "foo", Count: 0},
					{LineNumber: 3, FunctionName: "foo", Count: 0},
					{LineNumber: 4, FunctionName: "foo", Count: 1},
					{LineNumber: 5, FunctionName: "foo", Count: 1},
				},
				Functions: []Function{
					{Name: "foo", DemangledName: "foo()"},
				},
			},
		},
	}

	result, err := FindUncoveredLines(report)
	if err != nil {
		t.Fatalf("FindUncoveredLines() error = %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(result.Files))
	}

	file := result.Files[0]
	if len(file.UncoveredFunctions) != 1 {
		t.Fatalf("Expected 1 function, got %d", len(file.UncoveredFunctions))
	}

	fn := file.UncoveredFunctions[0]

	if fn.TotalLines != 5 {
		t.Errorf("Expected TotalLines=5, got %d", fn.TotalLines)
	}
	if fn.CoveredLines != 3 {
		t.Errorf("Expected CoveredLines=3, got %d", fn.CoveredLines)
	}
	if len(fn.UncoveredLineNumbers) != 2 {
		t.Errorf("Expected 2 uncovered lines, got %d", len(fn.UncoveredLineNumbers))
	}

	// Verify line numbers are sorted
	if fn.UncoveredLineNumbers[0] != 2 || fn.UncoveredLineNumbers[1] != 3 {
		t.Errorf("Expected uncovered lines [2, 3], got %v", fn.UncoveredLineNumbers)
	}
}

func TestFormatUncoveredReport(t *testing.T) {
	tests := []struct {
		name     string
		report   *UncoveredReport
		contains []string
	}{
		{
			name: "Report with uncovered lines",
			report: &UncoveredReport{
				Files: []FileUncovered{
					{
						FilePath: "test.cpp",
						UncoveredFunctions: []FunctionUncovered{
							{
								FunctionName:         "foo",
								DemangledName:        "foo()",
								UncoveredLineNumbers: []int{2, 3},
								TotalLines:           5,
								CoveredLines:         3,
							},
						},
					},
				},
			},
			contains: []string{
				"Uncovered Lines Report",
				"Found 1 function(s)",
				"2 total uncovered lines",
				"File: test.cpp",
				"Function: foo()",
				"Coverage: 3/5 lines (60.0%)",
				"Uncovered Lines (2): [2 3]",
			},
		},
		{
			name: "Empty report",
			report: &UncoveredReport{
				Files: []FileUncovered{},
			},
			contains: []string{
				"No uncovered lines found",
			},
		},
		{
			name: "Multiple files",
			report: &UncoveredReport{
				Files: []FileUncovered{
					{
						FilePath: "a.cpp",
						UncoveredFunctions: []FunctionUncovered{
							{
								FunctionName:         "func1",
								DemangledName:        "func1()",
								UncoveredLineNumbers: []int{1},
								TotalLines:           3,
								CoveredLines:         2,
							},
						},
					},
					{
						FilePath: "b.cpp",
						UncoveredFunctions: []FunctionUncovered{
							{
								FunctionName:         "func2",
								DemangledName:        "func2()",
								UncoveredLineNumbers: []int{10},
								TotalLines:           2,
								CoveredLines:         1,
							},
						},
					},
				},
			},
			contains: []string{
				"Found 2 function(s)",
				"2 total uncovered lines",
				"File: a.cpp",
				"Function: func1()",
				"File: b.cpp",
				"Function: func2()",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatUncoveredReport(tt.report)

			for _, substr := range tt.contains {
				if !contains(result, substr) {
					t.Errorf("Expected output to contain %q, but it doesn't.\nOutput: %s", substr, result)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
