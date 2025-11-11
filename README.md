# gcovr-json-util

A utility tool for processing and analyzing gcovr JSON coverage reports.

## Features

- **Coverage Diff**: Compare two gcovr JSON reports to identify coverage increases
- **Filtering Support**: Filter coverage tracking by specific files and functions using a YAML config
- Reports which functions have improved coverage
- Shows old and new coverage percentages
- Displays newly covered line numbers
- Uses demangled function names for readability
- Can be used both as a CLI tool and as a Go library

## Installation

```bash
go install github.com/zjy-dev/gcovr-json-util@latest
```

Or install a specific version:

```bash
# Install v2.0.0 (with filtering support)
go install github.com/zjy-dev/gcovr-json-util@v2.0.0

# Install v1.0.0 (basic functionality)
go install github.com/zjy-dev/gcovr-json-util@v1.0.0
```

### Updating to v2.0.0 in Your Go Project

If you're using this as a library in your Go project:

**Option 1: Update to latest version**
```bash
go get -u github.com/zjy-dev/gcovr-json-util@latest
go mod tidy
```

**Option 2: Update to specific v2.0.0**
```bash
go get github.com/zjy-dev/gcovr-json-util@v2.0.0
go mod tidy
```

**Option 3: Edit go.mod directly**
```go
require (
    github.com/zjy-dev/gcovr-json-util v2.0.0
)
```
Then run:
```bash
go mod tidy
```

**Check your current version:**
```bash
go list -m github.com/zjy-dev/gcovr-json-util
```

### Building from Source

Or build from source:

```bash
git clone https://github.com/zjy-dev/gcovr-json-util
cd gcovr-json-util
make build
# Or without make:
go build -o gcovr-util .
```

### Building with Version Information

To build with version information embedded:

```bash
make build
# This automatically embeds version, git commit, and build date

# Or manually:
go build -ldflags "-X main.Version=v1.0.0 -X main.GitCommit=$(git rev-parse --short HEAD) -X main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o gcovr-util .
```

Check version:
```bash
./gcovr-util --version
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
- `--filter, -f`: Filter config file (YAML) to specify target files and functions (optional)

#### Using Filter Configuration

You can use a YAML configuration file to filter which files and functions to track:

```bash
./gcovr-util diff --base base.json --new new.json --filter filter.yaml
```

**Filter Config Format** (`filter.yaml`):

```yaml
compiler:
  path: "/usr/bin/gcc"
  gcovr_exec_path: "/path/to/build"

targets:
  - file: "demo.cc"
    functions:
      - "f"
      - "g"
  - file: "another_file.cpp"
    functions:
      - "myFunction"
      - "anotherFunction"
```

This will only report coverage increases for the specified functions in the specified files. All other files and functions will be ignored.

**Note**: 
- File paths can be specified as relative paths, absolute paths, or just filenames
- Function names should match the demangled names (e.g., "f" instead of "_Z1fv")
- The `*.json` files and filter config file paths support both relative and absolute paths

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

// Optional: Apply filtering
filterConfig, err := gcovr.ParseFilterConfig("filter.yaml")
if err != nil {
    log.Fatal(err)
}
baseReport = gcovr.ApplyFilter(baseReport, filterConfig)
newReport = gcovr.ApplyFilter(newReport, filterConfig)

// Compute coverage increase
report, err := gcovr.ComputeCoverageIncrease(baseReport, newReport)
if err != nil {
    log.Fatal(err)
}

// Format and display results
output := gcovr.FormatReport(report)
fmt.Print(output)
```

## Version History & Migration

### v2.0.0 (Current) - November 12, 2025

**New Features:**
- ✨ Filter configuration support via YAML files
- 🎯 Selective coverage tracking by file and function
- 📋 New `--filter` flag
- 🔧 New library functions: `ParseFilterConfig()`, `ApplyFilter()`

**Dependencies:**
- Added: `gopkg.in/yaml.v3`

**Migration from v1.0.0:**
```bash
# Update dependency
go get github.com/zjy-dev/gcovr-json-util@v2.0.0
go mod tidy
```

**Code changes (optional, for new filtering feature):**
```go
// v1.0.0 code (still works in v2.0.0)
baseReport, _ := gcovr.ParseReport("base.json")
newReport, _ := gcovr.ParseReport("new.json")
report, _ := gcovr.ComputeCoverageIncrease(baseReport, newReport)

// v2.0.0 new feature (optional)
filterConfig, _ := gcovr.ParseFilterConfig("filter.yaml")
baseReport = gcovr.ApplyFilter(baseReport, filterConfig)
newReport = gcovr.ApplyFilter(newReport, filterConfig)
report, _ := gcovr.ComputeCoverageIncrease(baseReport, newReport)
```

**Backward Compatibility:** ✅ Fully backward compatible - all v1.0.0 code works in v2.0.0 without changes.

### v1.0.0 - November 11, 2025

**Features:**
- Basic coverage diff functionality
- Version information (`--version`)
- CLI and library usage

## Project Structure

```
.
├── main.go              # CLI entry point
├── version.go           # Version information
├── cmd/                 # CLI commands
│   ├── root.go         # Root command
│   └── diff.go         # Diff command implementation
├── pkg/
│   └── gcovr/          # Public library package
│       ├── types.go    # Data structures
│       ├── parser.go   # JSON parsing
│       ├── diff.go     # Coverage diff logic
│       └── filter.go   # Filter configuration
├── test_data/          # Sample test files
│   ├── f.json
│   ├── g.json
│   ├── m.json
│   ├── filter.yaml              # Example filter config
│   └── filter-f-only.yaml       # Another filter example
├── Makefile            # Build automation
├── CHANGELOG.md        # Version history
└── README.md           # This file
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
