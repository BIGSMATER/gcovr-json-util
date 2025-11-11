# gcovr-json-util

A utility tool for processing and analyzing gcovr JSON coverage reports.

## Features

- **Coverage Diff**: Compare two gcovr JSON reports to identify coverage increases
- Reports which functions have improved coverage
- Shows old and new coverage percentages
- Displays newly covered line numbers
- Uses demangled function names for readability
- Can be used both as a CLI tool and as a Go library

## Installation

```bash
go get github.com/zjy-dev/gcovr-json-util
```

Or build from source:

```bash
git clone https://github.com/zjy-dev/gcovr-json-util
cd gcovr-json-util
go build -o gcovr-util .
```

## Usage

### CLI Tool

Compare two gcovr JSON reports:

```bash
./gcovr-util diff --base base_coverage.json --new new_coverage.json
```

#### Options

- `--base, -b`: Base gcovr JSON report file (required)
- `--new, -n`: New gcovr JSON report file (required)

#### Example Output

```
Coverage Increase Report
=========================

Found 2 function(s) with increased coverage:

1. File: demo.cc
   Function: g()
   Old Coverage: 0/3 lines (0.0%)
   New Coverage: 3/3 lines (100.0%)
   Lines Increased: 3
   Newly Covered Line Numbers: [9 10 11]

2. File: demo.cc
   Function: main
   Old Coverage: 4/5 lines (80.0%)
   New Coverage: 5/5 lines (100.0%)
   Lines Increased: 1
   Newly Covered Line Numbers: [17]
```

### Go Library

You can also use this tool as a Go library in your projects:

```go
import "github.com/zjy-dev/gcovr-json-util/pkg/gcovr"

// Parse coverage reports
baseReport, err := gcovr.ParseReport("base.json")
if err != nil {
    log.Fatal(err)
}

newReport, err := gcovr.ParseReport("new.json")
if err != nil {
    log.Fatal(err)
}

// Compute coverage increase
report, err := gcovr.ComputeCoverageIncrease(baseReport, newReport)
if err != nil {
    log.Fatal(err)
}

// Format and display results
output := gcovr.FormatReport(report)
fmt.Print(output)
```

## Project Structure

```
.
├── main.go              # CLI entry point
├── cmd/                 # CLI commands
│   ├── root.go         # Root command
│   └── diff.go         # Diff command implementation
├── pkg/
│   └── gcovr/          # Public library package
│       ├── types.go    # Data structures
│       ├── parser.go   # JSON parsing
│       └── diff.go     # Coverage diff logic
└── test_data/          # Sample test files
    ├── f.json
    ├── g.json
    └── m.json
```

## How It Works

1. **Parse**: Reads and parses two gcovr JSON reports
2. **Compare**: Compares line-by-line coverage for each function
3. **Identify**: Identifies lines that were uncovered in base but are covered in new report
4. **Report**: Generates a detailed report with:
   - Function names (demangled for C++)
   - Old coverage percentage and line count
   - New coverage percentage and line count
   - Number of newly covered lines
   - Specific line numbers that gained coverage

## Use Cases

- **CI/CD Pipelines**: Automatically track coverage improvements in pull requests
- **Code Review**: Verify that new tests actually improve coverage
- **Quality Gates**: Ensure new code increases overall test coverage
- **Test Analysis**: Identify which functions benefited from new test cases

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
