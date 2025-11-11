# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[v1.0.0]: https://github.com/zjy-dev/gcovr-json-util/releases/tag/v1.0.0
