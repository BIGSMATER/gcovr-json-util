package gcovr

import (
	"testing"
)

func TestComputeCoverageIncrease(t *testing.T) {
	tests := []struct {
		name              string
		baseReport        *GcovrReport
		newReport         *GcovrReport
		expectedIncreases int
		expectedFirstFile string
		expectedFirstFunc string
	}{
		{
			name: "Basic coverage increase",
			baseReport: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 0},
							{LineNumber: 2, FunctionName: "foo", Count: 0},
							{LineNumber: 3, FunctionName: "foo", Count: 1},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			newReport: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
							{LineNumber: 2, FunctionName: "foo", Count: 1},
							{LineNumber: 3, FunctionName: "foo", Count: 1},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			expectedIncreases: 1,
			expectedFirstFile: "test.cpp",
			expectedFirstFunc: "foo()",
		},
		{
			name: "No coverage increase",
			baseReport: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
							{LineNumber: 2, FunctionName: "foo", Count: 1},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			newReport: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
							{LineNumber: 2, FunctionName: "foo", Count: 1},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			expectedIncreases: 0,
		},
		{
			name: "New file in new report",
			baseReport: &GcovrReport{
				Files: []File{},
			},
			newReport: &GcovrReport{
				Files: []File{
					{
						FilePath: "new.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "bar", Count: 1},
							{LineNumber: 2, FunctionName: "bar", Count: 1},
						},
						Functions: []Function{
							{Name: "bar", DemangledName: "bar()"},
						},
					},
				},
			},
			expectedIncreases: 1,
			expectedFirstFile: "new.cpp",
			expectedFirstFunc: "bar()",
		},
		{
			name: "Multiple functions with increases",
			baseReport: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 0},
							{LineNumber: 5, FunctionName: "bar", Count: 0},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
							{Name: "bar", DemangledName: "bar()"},
						},
					},
				},
			},
			newReport: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
							{LineNumber: 5, FunctionName: "bar", Count: 1},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
							{Name: "bar", DemangledName: "bar()"},
						},
					},
				},
			},
			expectedIncreases: 2,
			expectedFirstFile: "test.cpp",
		},
		{
			name: "Empty reports",
			baseReport: &GcovrReport{
				Files: []File{},
			},
			newReport: &GcovrReport{
				Files: []File{},
			},
			expectedIncreases: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ComputeCoverageIncrease(tt.baseReport, tt.newReport)
			if err != nil {
				t.Fatalf("ComputeCoverageIncrease() error = %v", err)
			}

			if len(result.Increases) != tt.expectedIncreases {
				t.Errorf("Expected %d increases, got %d", tt.expectedIncreases, len(result.Increases))
			}

			if tt.expectedIncreases > 0 {
				firstInc := result.Increases[0]
				if firstInc.File != tt.expectedFirstFile {
					t.Errorf("Expected first file %s, got %s", tt.expectedFirstFile, firstInc.File)
				}
				if tt.expectedFirstFunc != "" && firstInc.DemangledName != tt.expectedFirstFunc {
					t.Errorf("Expected first function %s, got %s", tt.expectedFirstFunc, firstInc.DemangledName)
				}
			}
		})
	}
}

func TestComputeCoverageIncrease_DetailedMetrics(t *testing.T) {
	baseReport := &GcovrReport{
		Files: []File{
			{
				FilePath: "demo.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "foo", Count: 1},
					{LineNumber: 2, FunctionName: "foo", Count: 0},
					{LineNumber: 3, FunctionName: "foo", Count: 0},
					{LineNumber: 4, FunctionName: "foo", Count: 1},
				},
				Functions: []Function{
					{Name: "foo", DemangledName: "foo()"},
				},
			},
		},
	}

	newReport := &GcovrReport{
		Files: []File{
			{
				FilePath: "demo.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "foo", Count: 1},
					{LineNumber: 2, FunctionName: "foo", Count: 1},
					{LineNumber: 3, FunctionName: "foo", Count: 1},
					{LineNumber: 4, FunctionName: "foo", Count: 1},
				},
				Functions: []Function{
					{Name: "foo", DemangledName: "foo()"},
				},
			},
		},
	}

	result, err := ComputeCoverageIncrease(baseReport, newReport)
	if err != nil {
		t.Fatalf("ComputeCoverageIncrease() error = %v", err)
	}

	if len(result.Increases) != 1 {
		t.Fatalf("Expected 1 increase, got %d", len(result.Increases))
	}

	inc := result.Increases[0]

	if inc.OldCoveredLines != 2 {
		t.Errorf("Expected OldCoveredLines=2, got %d", inc.OldCoveredLines)
	}
	if inc.NewCoveredLines != 4 {
		t.Errorf("Expected NewCoveredLines=4, got %d", inc.NewCoveredLines)
	}
	if inc.LinesIncreased != 2 {
		t.Errorf("Expected LinesIncreased=2, got %d", inc.LinesIncreased)
	}
	if inc.TotalLines != 4 {
		t.Errorf("Expected TotalLines=4, got %d", inc.TotalLines)
	}
	if len(inc.IncreasedLineNumbers) != 2 {
		t.Errorf("Expected 2 increased line numbers, got %d", len(inc.IncreasedLineNumbers))
	}
	if inc.FunctionName != "foo" {
		t.Errorf("Expected FunctionName='foo', got '%s'", inc.FunctionName)
	}
	if inc.DemangledName != "foo()" {
		t.Errorf("Expected DemangledName='foo()', got '%s'", inc.DemangledName)
	}
}

func TestComputeCoverageIncrease_NewFunctionInExistingFile(t *testing.T) {
	baseReport := &GcovrReport{
		Files: []File{
			{
				FilePath: "test.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "foo", Count: 1},
				},
				Functions: []Function{
					{Name: "foo", DemangledName: "foo()"},
				},
			},
		},
	}

	newReport := &GcovrReport{
		Files: []File{
			{
				FilePath: "test.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "foo", Count: 1},
					{LineNumber: 5, FunctionName: "bar", Count: 1},
					{LineNumber: 6, FunctionName: "bar", Count: 1},
				},
				Functions: []Function{
					{Name: "foo", DemangledName: "foo()"},
					{Name: "bar", DemangledName: "bar()"},
				},
			},
		},
	}

	result, err := ComputeCoverageIncrease(baseReport, newReport)
	if err != nil {
		t.Fatalf("ComputeCoverageIncrease() error = %v", err)
	}

	// Should report increase for the new function 'bar'
	if len(result.Increases) != 1 {
		t.Fatalf("Expected 1 increase for new function, got %d", len(result.Increases))
	}

	inc := result.Increases[0]
	if inc.FunctionName != "bar" {
		t.Errorf("Expected FunctionName='bar', got '%s'", inc.FunctionName)
	}
	if inc.OldCoveredLines != 0 {
		t.Errorf("Expected OldCoveredLines=0 for new function, got %d", inc.OldCoveredLines)
	}
	if inc.NewCoveredLines != 2 {
		t.Errorf("Expected NewCoveredLines=2, got %d", inc.NewCoveredLines)
	}
}

func TestFormatReport(t *testing.T) {
	tests := []struct {
		name     string
		report   *CoverageIncreaseReport
		contains []string
	}{
		{
			name: "Single increase",
			report: &CoverageIncreaseReport{
				Increases: []FunctionCoverageIncrease{
					{
						File:                 "test.cpp",
						FunctionName:         "foo",
						DemangledName:        "foo()",
						LinesIncreased:       2,
						TotalLines:           5,
						IncreasedLineNumbers: []int{2, 3},
						OldCoveredLines:      3,
						NewCoveredLines:      5,
					},
				},
			},
			contains: []string{
				"Coverage Increase Report",
				"Found 1 function(s)",
				"File: test.cpp",
				"Function: foo()",
				"Old Coverage: 3/5 lines (60.0%)",
				"New Coverage: 5/5 lines (100.0%)",
				"Lines Increased: 2",
				"Newly Covered Line Numbers: [2 3]",
			},
		},
		{
			name: "No increases",
			report: &CoverageIncreaseReport{
				Increases: []FunctionCoverageIncrease{},
			},
			contains: []string{
				"No coverage increases found",
			},
		},
		{
			name: "Multiple increases",
			report: &CoverageIncreaseReport{
				Increases: []FunctionCoverageIncrease{
					{
						File:                 "a.cpp",
						DemangledName:        "funcA()",
						LinesIncreased:       1,
						TotalLines:           3,
						IncreasedLineNumbers: []int{5},
						OldCoveredLines:      2,
						NewCoveredLines:      3,
					},
					{
						File:                 "b.cpp",
						DemangledName:        "funcB()",
						LinesIncreased:       2,
						TotalLines:           4,
						IncreasedLineNumbers: []int{10, 11},
						OldCoveredLines:      2,
						NewCoveredLines:      4,
					},
				},
			},
			contains: []string{
				"Found 2 function(s)",
				"File: a.cpp",
				"Function: funcA()",
				"File: b.cpp",
				"Function: funcB()",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatReport(tt.report)

			for _, substr := range tt.contains {
				if !containsString(result, substr) {
					t.Errorf("Expected output to contain %q, but it doesn't.\nOutput: %s", substr, result)
				}
			}
		})
	}
}

func TestBuildLineCoverageMap(t *testing.T) {
	file := &File{
		Lines: []Line{
			{LineNumber: 1, FunctionName: "foo", Count: 1},
			{LineNumber: 2, FunctionName: "foo", Count: 0},
			{LineNumber: 5, FunctionName: "bar", Count: 5},
		},
	}

	result := buildLineCoverageMap(file)

	if len(result) != 2 {
		t.Errorf("Expected 2 functions in map, got %d", len(result))
	}

	if _, exists := result["foo"]; !exists {
		t.Error("Expected 'foo' to exist in map")
	}
	if _, exists := result["bar"]; !exists {
		t.Error("Expected 'bar' to exist in map")
	}

	if result["foo"][1] != 1 {
		t.Errorf("Expected foo line 1 count=1, got %d", result["foo"][1])
	}
	if result["foo"][2] != 0 {
		t.Errorf("Expected foo line 2 count=0, got %d", result["foo"][2])
	}
	if result["bar"][5] != 5 {
		t.Errorf("Expected bar line 5 count=5, got %d", result["bar"][5])
	}
}

func TestBuildFunctionNameMap(t *testing.T) {
	file := &File{
		Functions: []Function{
			{Name: "_Z3foov", DemangledName: "foo()"},
			{Name: "_Z3barv", DemangledName: "bar()"},
		},
	}

	result := buildFunctionNameMap(file)

	if len(result) != 2 {
		t.Errorf("Expected 2 functions in map, got %d", len(result))
	}

	if result["_Z3foov"] != "foo()" {
		t.Errorf("Expected '_Z3foov' -> 'foo()', got '%s'", result["_Z3foov"])
	}
	if result["_Z3barv"] != "bar()" {
		t.Errorf("Expected '_Z3barv' -> 'bar()', got '%s'", result["_Z3barv"])
	}
}

func TestGetTotalFunctionLines(t *testing.T) {
	file := &File{
		Lines: []Line{
			{LineNumber: 1, FunctionName: "foo", Count: 1},
			{LineNumber: 2, FunctionName: "foo", Count: 0},
			{LineNumber: 3, FunctionName: "foo", Count: 1},
			{LineNumber: 10, FunctionName: "bar", Count: 1},
			{LineNumber: 11, FunctionName: "bar", Count: 1},
		},
	}

	fooTotal := getTotalFunctionLines(file, "foo")
	if fooTotal != 3 {
		t.Errorf("Expected foo to have 3 lines, got %d", fooTotal)
	}

	barTotal := getTotalFunctionLines(file, "bar")
	if barTotal != 2 {
		t.Errorf("Expected bar to have 2 lines, got %d", barTotal)
	}

	nonExistent := getTotalFunctionLines(file, "nonexistent")
	if nonExistent != 0 {
		t.Errorf("Expected nonexistent function to have 0 lines, got %d", nonExistent)
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
