package gcovr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseReport(t *testing.T) {
	tests := []struct {
		name          string
		createFile    bool
		fileContent   string
		expectedError bool
		expectedFiles int
	}{
		{
			name:       "Valid JSON report",
			createFile: true,
			fileContent: `{
				"gcovr/format_version": "0.5",
				"files": [
					{
						"file": "test.cpp",
						"lines": [
							{"line_number": 1, "function_name": "foo", "count": 1},
							{"line_number": 2, "function_name": "foo", "count": 0}
						],
						"functions": [
							{
								"name": "_Z3foov",
								"demangled_name": "foo()",
								"lineno": 1,
								"execution_count": 1,
								"blocks_percent": 100.0,
								"pos": ["1:1"]
							}
						]
					}
				]
			}`,
			expectedError: false,
			expectedFiles: 1,
		},
		{
			name:       "Empty files array",
			createFile: true,
			fileContent: `{
				"gcovr/format_version": "0.5",
				"files": []
			}`,
			expectedError: false,
			expectedFiles: 0,
		},
		{
			name:          "File does not exist",
			createFile:    false,
			expectedError: true,
		},
		{
			name:       "Invalid JSON",
			createFile: true,
			fileContent: `{
				"gcovr/format_version": "0.5",
				"files": [
					invalid json
				]
			}`,
			expectedError: true,
		},
		{
			name:       "Multiple files",
			createFile: true,
			fileContent: `{
				"gcovr/format_version": "0.5",
				"files": [
					{
						"file": "file1.cpp",
						"lines": [],
						"functions": []
					},
					{
						"file": "file2.cpp",
						"lines": [],
						"functions": []
					}
				]
			}`,
			expectedError: false,
			expectedFiles: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			
			if tt.createFile {
				// Create a temporary file
				tmpFile, err := os.CreateTemp("", "gcovr_test_*.json")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer os.Remove(tmpFile.Name())
				filePath = tmpFile.Name()

				if _, err := tmpFile.WriteString(tt.fileContent); err != nil {
					t.Fatalf("Failed to write to temp file: %v", err)
				}
				tmpFile.Close()
			} else {
				filePath = "nonexistent_file.json"
			}

			result, err := ParseReport(filePath)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Fatal("Expected non-nil result")
				}
				if len(result.Files) != tt.expectedFiles {
					t.Errorf("Expected %d files, got %d", tt.expectedFiles, len(result.Files))
				}
			}
		})
	}
}

func TestParseReport_ActualTestData(t *testing.T) {
	// Test with actual test data files if they exist
	testDataDir := filepath.Join("..", "..", "test_data")
	testFiles := []struct {
		name     string
		filename string
	}{
		{"f.json", "f.json"},
		{"g.json", "g.json"},
		{"m.json", "m.json"},
	}

	for _, tt := range testFiles {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(testDataDir, tt.filename)
			
			// Check if file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Skipf("Test data file %s does not exist, skipping", filePath)
				return
			}

			result, err := ParseReport(filePath)
			if err != nil {
				t.Errorf("Failed to parse %s: %v", tt.filename, err)
			}
			if result == nil {
				t.Errorf("Expected non-nil result for %s", tt.filename)
			}
			if result != nil && result.FormatVersion == "" {
				t.Errorf("Expected FormatVersion to be set for %s", tt.filename)
			}
		})
	}
}

func TestParseReport_DataStructure(t *testing.T) {
	fileContent := `{
		"gcovr/format_version": "0.5",
		"files": [
			{
				"file": "demo.cpp",
				"lines": [
					{"line_number": 5, "function_name": "_Z1fv", "count": 10},
					{"line_number": 6, "function_name": "_Z1fv", "count": 0},
					{"line_number": 7, "function_name": "_Z1fv", "count": 5}
				],
				"functions": [
					{
						"name": "_Z1fv",
						"demangled_name": "f()",
						"lineno": 5,
						"execution_count": 1,
						"blocks_percent": 75.5,
						"pos": ["5:1", "7:2"]
					}
				]
			}
		]
	}`

	tmpFile, err := os.CreateTemp("", "gcovr_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(fileContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	result, err := ParseReport(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseReport failed: %v", err)
	}

	// Verify format version
	if result.FormatVersion != "0.5" {
		t.Errorf("Expected FormatVersion='0.5', got '%s'", result.FormatVersion)
	}

	// Verify file structure
	if len(result.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(result.Files))
	}

	file := result.Files[0]
	if file.FilePath != "demo.cpp" {
		t.Errorf("Expected FilePath='demo.cpp', got '%s'", file.FilePath)
	}

	// Verify lines
	if len(file.Lines) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(file.Lines))
	}

	line := file.Lines[0]
	if line.LineNumber != 5 {
		t.Errorf("Expected LineNumber=5, got %d", line.LineNumber)
	}
	if line.FunctionName != "_Z1fv" {
		t.Errorf("Expected FunctionName='_Z1fv', got '%s'", line.FunctionName)
	}
	if line.Count != 10 {
		t.Errorf("Expected Count=10, got %d", line.Count)
	}

	// Verify functions
	if len(file.Functions) != 1 {
		t.Fatalf("Expected 1 function, got %d", len(file.Functions))
	}

	fn := file.Functions[0]
	if fn.Name != "_Z1fv" {
		t.Errorf("Expected Name='_Z1fv', got '%s'", fn.Name)
	}
	if fn.DemangledName != "f()" {
		t.Errorf("Expected DemangledName='f()', got '%s'", fn.DemangledName)
	}
	if fn.LineNo != 5 {
		t.Errorf("Expected LineNo=5, got %d", fn.LineNo)
	}
	if fn.ExecutionCount != 1 {
		t.Errorf("Expected ExecutionCount=1, got %d", fn.ExecutionCount)
	}
	if fn.BlocksPercent != 75.5 {
		t.Errorf("Expected BlocksPercent=75.5, got %f", fn.BlocksPercent)
	}
	if len(fn.Pos) != 2 {
		t.Errorf("Expected 2 positions, got %d", len(fn.Pos))
	}
}
