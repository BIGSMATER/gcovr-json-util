package gcovr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFilterConfig(t *testing.T) {
	tests := []struct {
		name           string
		createFile     bool
		fileContent    string
		expectedError  bool
		expectedFiles  int
		expectedFuncs  int // functions in first target
	}{
		{
			name:       "Valid filter config",
			createFile: true,
			fileContent: `compiler:
  path: "/usr/bin/gcc"
  gcovr_exec_path: "/path/to/build"

targets:
  - file: "demo.cc"
    functions:
      - "f"
      - "g"
`,
			expectedError: false,
			expectedFiles: 1,
			expectedFuncs: 2,
		},
		{
			name:       "Multiple target files",
			createFile: true,
			fileContent: `compiler:
  path: "/usr/bin/gcc"
  gcovr_exec_path: "/path/to/build"

targets:
  - file: "file1.cpp"
    functions:
      - "funcA"
  - file: "file2.cpp"
    functions:
      - "funcB"
      - "funcC"
`,
			expectedError: false,
			expectedFiles: 2,
			expectedFuncs: 1,
		},
		{
			name:       "Empty targets",
			createFile: true,
			fileContent: `compiler:
  path: "/usr/bin/gcc"
  gcovr_exec_path: "/path/to/build"

targets: []
`,
			expectedError: false,
			expectedFiles: 0,
			expectedFuncs: 0,
		},
		{
			name:          "File does not exist",
			createFile:    false,
			expectedError: true,
		},
		{
			name:       "Invalid YAML",
			createFile: true,
			fileContent: `compiler:
  path: "/usr/bin/gcc"
  invalid yaml content [[[
`,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string

			if tt.createFile {
				tmpFile, err := os.CreateTemp("", "filter_test_*.yaml")
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
				filePath = "nonexistent_filter.yaml"
			}

			result, err := ParseFilterConfig(filePath)

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
				if len(result.Targets) != tt.expectedFiles {
					t.Errorf("Expected %d target files, got %d", tt.expectedFiles, len(result.Targets))
				}
				if tt.expectedFiles > 0 && len(result.Targets[0].Functions) != tt.expectedFuncs {
					t.Errorf("Expected %d functions in first target, got %d", 
						tt.expectedFuncs, len(result.Targets[0].Functions))
				}
			}
		})
	}
}

func TestParseFilterConfig_ActualTestData(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "test_data")
	testFiles := []struct {
		name     string
		filename string
	}{
		{"filter.yaml", "filter.yaml"},
		{"filter-f-only.yaml", "filter-f-only.yaml"},
	}

	for _, tt := range testFiles {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(testDataDir, tt.filename)

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Skipf("Test data file %s does not exist, skipping", filePath)
				return
			}

			result, err := ParseFilterConfig(filePath)
			if err != nil {
				t.Errorf("Failed to parse %s: %v", tt.filename, err)
			}
			if result == nil {
				t.Errorf("Expected non-nil result for %s", tt.filename)
			}
		})
	}
}

func TestApplyFilter(t *testing.T) {
	tests := []struct {
		name           string
		report         *GcovrReport
		filterConfig   *FilterConfig
		expectedFiles  int
		expectedFuncs  int // functions in first file (if exists)
	}{
		{
			name: "Filter single file and function",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "demo.cc",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "_Z1fv", Count: 1},
							{LineNumber: 5, FunctionName: "_Z1gv", Count: 0},
						},
						Functions: []Function{
							{Name: "_Z1fv", DemangledName: "f()"},
							{Name: "_Z1gv", DemangledName: "g()"},
						},
					},
				},
			},
			filterConfig: &FilterConfig{
				Targets: []TargetFile{
					{
						File:      "demo.cc",
						Functions: []string{"f"},
					},
				},
			},
			expectedFiles: 1,
			expectedFuncs: 1,
		},
		{
			name: "Filter multiple functions",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
							{LineNumber: 5, FunctionName: "bar", Count: 0},
							{LineNumber: 10, FunctionName: "baz", Count: 1},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
							{Name: "bar", DemangledName: "bar()"},
							{Name: "baz", DemangledName: "baz()"},
						},
					},
				},
			},
			filterConfig: &FilterConfig{
				Targets: []TargetFile{
					{
						File:      "test.cpp",
						Functions: []string{"foo", "bar"},
					},
				},
			},
			expectedFiles: 1,
			expectedFuncs: 2,
		},
		{
			name: "No matching files",
			report: &GcovrReport{
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
			},
			filterConfig: &FilterConfig{
				Targets: []TargetFile{
					{
						File:      "nonexistent.cpp",
						Functions: []string{"foo"},
					},
				},
			},
			expectedFiles: 0,
			expectedFuncs: 0,
		},
		{
			name: "Empty filter config",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines:    []Line{{LineNumber: 1, FunctionName: "foo", Count: 1}},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			filterConfig:  &FilterConfig{Targets: []TargetFile{}},
			expectedFiles: 1, // No filtering applied
			expectedFuncs: 1,
		},
		{
			name: "Nil filter config",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "test.cpp",
						Lines:    []Line{{LineNumber: 1, FunctionName: "foo", Count: 1}},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			filterConfig:  nil,
			expectedFiles: 1, // No filtering applied
			expectedFuncs: 1,
		},
		{
			name: "Match by filename only",
			report: &GcovrReport{
				Files: []File{
					{
						FilePath: "/path/to/demo.cc",
						Lines: []Line{
							{LineNumber: 1, FunctionName: "foo", Count: 1},
						},
						Functions: []Function{
							{Name: "foo", DemangledName: "foo()"},
						},
					},
				},
			},
			filterConfig: &FilterConfig{
				Targets: []TargetFile{
					{
						File:      "demo.cc",
						Functions: []string{"foo"},
					},
				},
			},
			expectedFiles: 1,
			expectedFuncs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyFilter(tt.report, tt.filterConfig)

			if len(result.Files) != tt.expectedFiles {
				t.Errorf("Expected %d files, got %d", tt.expectedFiles, len(result.Files))
			}

			if tt.expectedFiles > 0 && len(result.Files[0].Functions) != tt.expectedFuncs {
				t.Errorf("Expected %d functions in first file, got %d", 
					tt.expectedFuncs, len(result.Files[0].Functions))
			}
		})
	}
}

func TestApplyFilter_LinesFiltering(t *testing.T) {
	report := &GcovrReport{
		Files: []File{
			{
				FilePath: "test.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "foo", Count: 1},
					{LineNumber: 2, FunctionName: "foo", Count: 0},
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

	filterConfig := &FilterConfig{
		Targets: []TargetFile{
			{
				File:      "test.cpp",
				Functions: []string{"foo"},
			},
		},
	}

	result := ApplyFilter(report, filterConfig)

	if len(result.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(result.Files))
	}

	file := result.Files[0]

	// Should only have lines from function 'foo'
	if len(file.Lines) != 2 {
		t.Errorf("Expected 2 lines (from foo), got %d", len(file.Lines))
	}

	for _, line := range file.Lines {
		if line.FunctionName != "foo" {
			t.Errorf("Expected all lines to be from 'foo', got line from '%s'", line.FunctionName)
		}
	}

	// Should only have function 'foo'
	if len(file.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(file.Functions))
	}
	if file.Functions[0].Name != "foo" {
		t.Errorf("Expected function 'foo', got '%s'", file.Functions[0].Name)
	}
}

func TestNormalizeFilePath(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Unix path",
			input: "/path/to/file.cpp",
		},
		{
			name:  "Windows path",
			input: "C:\\path\\to\\file.cpp",
		},
		{
			name:  "Relative path",
			input: "src/test.cpp",
		},
		{
			name:  "Path with dots",
			input: "./src/../lib/file.cpp",
		},
		{
			name:  "Mixed slashes",
			input: "path\\to/file.cpp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeFilePath(tt.input)
			
			// Verify the function returns a consistent result
			expected := filepath.ToSlash(filepath.Clean(tt.input))
			if result != expected {
				t.Errorf("Expected '%s', got '%s'", expected, result)
			}
			
			// Verify result is not empty
			if result == "" {
				t.Error("Expected non-empty result")
			}
			
			// Verify the result is cleaned (no double slashes, no trailing slashes except root)
			if result != "." && result != "/" && len(result) > 1 && result[len(result)-1] == '/' {
				t.Errorf("Expected no trailing slash, got '%s'", result)
			}
		})
	}
}

func TestShouldIncludeFunction(t *testing.T) {
	tests := []struct {
		name             string
		demangledName    string
		mangledName      string
		allowedFunctions map[string]bool
		expected         bool
	}{
		{
			name:             "Match by simple demangled name",
			demangledName:    "foo()",
			mangledName:      "_Z3foov",
			allowedFunctions: map[string]bool{"foo": true},
			expected:         true,
		},
		{
			name:             "Match by full demangled name",
			demangledName:    "foo()",
			mangledName:      "_Z3foov",
			allowedFunctions: map[string]bool{"foo()": true},
			expected:         true,
		},
		{
			name:             "Match by mangled name",
			demangledName:    "foo()",
			mangledName:      "_Z3foov",
			allowedFunctions: map[string]bool{"_Z3foov": true},
			expected:         true,
		},
		{
			name:             "No match",
			demangledName:    "foo()",
			mangledName:      "_Z3foov",
			allowedFunctions: map[string]bool{"bar": true},
			expected:         false,
		},
		{
			name:             "Match with parameters",
			demangledName:    "calculate(int, double)",
			mangledName:      "_Z9calculateid",
			allowedFunctions: map[string]bool{"calculate": true},
			expected:         true,
		},
		{
			name:             "Empty allowed functions",
			demangledName:    "foo()",
			mangledName:      "_Z3foov",
			allowedFunctions: map[string]bool{},
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldIncludeFunction(tt.demangledName, tt.mangledName, tt.allowedFunctions)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestApplyFilter_MultipleFiles(t *testing.T) {
	report := &GcovrReport{
		Files: []File{
			{
				FilePath: "file1.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "func1", Count: 1},
				},
				Functions: []Function{
					{Name: "func1", DemangledName: "func1()"},
				},
			},
			{
				FilePath: "file2.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "func2", Count: 1},
				},
				Functions: []Function{
					{Name: "func2", DemangledName: "func2()"},
				},
			},
			{
				FilePath: "file3.cpp",
				Lines: []Line{
					{LineNumber: 1, FunctionName: "func3", Count: 1},
				},
				Functions: []Function{
					{Name: "func3", DemangledName: "func3()"},
				},
			},
		},
	}

	filterConfig := &FilterConfig{
		Targets: []TargetFile{
			{
				File:      "file1.cpp",
				Functions: []string{"func1"},
			},
			{
				File:      "file3.cpp",
				Functions: []string{"func3"},
			},
		},
	}

	result := ApplyFilter(report, filterConfig)

	// Should only have file1.cpp and file3.cpp
	if len(result.Files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(result.Files))
	}

	filePaths := make(map[string]bool)
	for _, file := range result.Files {
		filePaths[file.FilePath] = true
	}

	if !filePaths["file1.cpp"] {
		t.Error("Expected file1.cpp to be in result")
	}
	if filePaths["file2.cpp"] {
		t.Error("file2.cpp should not be in result")
	}
	if !filePaths["file3.cpp"] {
		t.Error("Expected file3.cpp to be in result")
	}
}

func TestApplyFilter_PreservesFormatVersion(t *testing.T) {
	report := &GcovrReport{
		FormatVersion: "0.5",
		Files: []File{
			{
				FilePath: "test.cpp",
				Lines:    []Line{{LineNumber: 1, FunctionName: "foo", Count: 1}},
				Functions: []Function{
					{Name: "foo", DemangledName: "foo()"},
				},
			},
		},
	}

	filterConfig := &FilterConfig{
		Targets: []TargetFile{
			{
				File:      "test.cpp",
				Functions: []string{"foo"},
			},
		},
	}

	result := ApplyFilter(report, filterConfig)

	if result.FormatVersion != "0.5" {
		t.Errorf("Expected FormatVersion='0.5', got '%s'", result.FormatVersion)
	}
}
