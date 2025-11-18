package gcovr

// GcovrReport represents the top-level structure of a gcovr JSON report
type GcovrReport struct {
	FormatVersion string `json:"gcovr/format_version"`
	Files         []File `json:"files"`
}

// File represents a source file in the coverage report
type File struct {
	FilePath  string     `json:"file"`
	Lines     []Line     `json:"lines"`
	Functions []Function `json:"functions"`
}

// Line represents a single line of code with coverage information
type Line struct {
	LineNumber   int    `json:"line_number"`
	FunctionName string `json:"function_name"`
	Count        int    `json:"count"`
}

// Function represents a function in the source code
type Function struct {
	Name           string   `json:"name"`
	DemangledName  string   `json:"demangled_name"`
	LineNo         int      `json:"lineno"`
	ExecutionCount int      `json:"execution_count"`
	BlocksPercent  float64  `json:"blocks_percent"`
	Pos            []string `json:"pos"`
}

// FunctionCoverageIncrease represents coverage increase for a specific function
type FunctionCoverageIncrease struct {
	File                 string
	FunctionName         string // Mangled name
	DemangledName        string
	LinesIncreased       int
	TotalLines           int
	IncreasedLineNumbers []int
	OldCoveredLines      int // Number of lines covered in base report
	NewCoveredLines      int // Number of lines covered in new report
}

// CoverageIncreaseReport contains all coverage increases between two reports
type CoverageIncreaseReport struct {
	Increases []FunctionCoverageIncrease
}

// FunctionUncovered represents the uncovered lines within a single function
type FunctionUncovered struct {
	FunctionName         string // Mangled name
	DemangledName        string
	UncoveredLineNumbers []int
	TotalLines           int
	CoveredLines         int
}

// FileUncovered represents all uncovered functions within a single file
type FileUncovered struct {
	FilePath           string
	UncoveredFunctions []FunctionUncovered
}

// UncoveredReport represents a complete report of all uncovered functions and lines, grouped by file
type UncoveredReport struct {
	Files []FileUncovered
}
