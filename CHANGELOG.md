# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v2.1.0] - 2025-11-19

### Added

- **Uncovered Lines Reporting** - New feature to identify lines without coverage
- `uncovered` CLI command (alias: `un`) to report uncovered lines from a gcovr JSON file
- `FindUncoveredLines()` API function to analyze coverage gaps
- `FormatUncoveredReport()` API function to format uncovered lines report
- New data structures: `FunctionUncovered` and `UncoveredReport`
- Support for filtering uncovered lines using `--filter` flag
- Detailed coverage statistics per function (covered/total lines, percentage)

### Changed

- Updated README with uncovered command documentation and examples
- Added library usage examples for finding uncovered lines
- Enhanced feature list to highlight uncovered lines reporting

## [v2.0.0] - 2025-11-12

### Added

- **Filter configuration support** via YAML files - **BREAKING FEATURE**
- `--filter` flag to specify target files and functions to track
- `ApplyFilter` and `ParseFilterConfig` functions in library API
- Support for filtering by file and function names
- Flexible path matching (relative, absolute, or filename only)
- Example filter configuration files (gcc-x64-canary.yaml, filter-example.yaml)
- Smart function name matching (handles mangled and demangled names)

### Changed

- Updated README with comprehensive filtering documentation
- Enhanced CLI help messages with filtering examples
- Added gopkg.in/yaml.v3 dependency for YAML parsing

### Why v2.0.0?

This is a major version bump because:

- New filtering feature significantly changes how the tool can be used
- New command-line flag changes the interface
- New dependency added (gopkg.in/yaml.v3)
- Enhanced API with new public functions

## [v1.0.0] - 2025-11-11

### Added

- Initial stable release
- Coverage diff functionality to compare two gcovr JSON reports
- Identify functions with increased line coverage
- Display old and new coverage percentages
- Show newly covered line numbers
- Demangled function names for C++ functions
- CLI tool built with Cobra framework
- Public Go library API in `pkg/gcovr` package
- Version information via `--version` flag
- Makefile for easy building with version embedding
- Complete documentation and examples
- MIT License

### Features

- **Coverage Increase Detection**: Compare base and new coverage reports
- **Detailed Reporting**: Show file, function, old/new coverage stats
- **Line-by-line Analysis**: Track specific lines that gained coverage
- **Name Demangling**: Display readable C++ function names
- **Dual Usage**: Both CLI tool and importable Go library

[v2.1.0]: https://github.com/zjy-dev/gcovr-json-util/releases/tag/v2.1.0
[v2.0.0]: https://github.com/zjy-dev/gcovr-json-util/releases/tag/v2.0.0
[v1.0.0]: https://github.com/zjy-dev/gcovr-json-util/releases/tag/v1.0.0
