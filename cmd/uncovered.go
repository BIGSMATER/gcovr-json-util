package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zjy-dev/gcovr-json-util/v2/pkg/gcovr"
)

var (
	uncoveredFilterFile string
)

// uncoveredCmd represents the uncovered command
var uncoveredCmd = &cobra.Command{
	Use:     "uncovered [gcovr-file]",
	Aliases: []string{"un"},
	Short:   "Report uncovered lines from a gcovr JSON report",
	Long: `Analyze a gcovr JSON report and display which lines are not covered,
grouped by file and function. This helps identify gaps in test coverage.

The tool will show:
- Which files have uncovered lines
- Which functions within those files have uncovered lines
- The specific line numbers that are not covered
- Coverage statistics for each function`,
	Args: cobra.ExactArgs(1),
	RunE: runUncovered,
}

func init() {
	rootCmd.AddCommand(uncoveredCmd)

	uncoveredCmd.Flags().StringVarP(&uncoveredFilterFile, "filter", "f", "",
		"Filter config file (YAML) to specify target files and functions")
}

func runUncovered(cmd *cobra.Command, args []string) error {
	reportFile := args[0]

	// Parse the gcovr JSON report
	fmt.Printf("Reading report: %s\n", reportFile)
	report, err := gcovr.ParseReport(reportFile)
	if err != nil {
		return fmt.Errorf("failed to parse report: %w", err)
	}

	// Apply filter if specified
	if uncoveredFilterFile != "" {
		fmt.Printf("Reading filter config: %s\n", uncoveredFilterFile)
		filterConfig, err := gcovr.ParseFilterConfig(uncoveredFilterFile)
		if err != nil {
			return fmt.Errorf("failed to parse filter config: %w", err)
		}

		fmt.Printf("Filtering enabled: tracking %d file(s)\n", len(filterConfig.Targets))
		report = gcovr.ApplyFilter(report, filterConfig)
		fmt.Println("Applying filters...")
	}

	// Find uncovered lines
	fmt.Println("\nAnalyzing coverage...\n")
	uncoveredReport, err := gcovr.FindUncoveredLines(report)
	if err != nil {
		return fmt.Errorf("failed to find uncovered lines: %w", err)
	}

	// Display results
	output := gcovr.FormatUncoveredReport(uncoveredReport)
	fmt.Print(output)

	return nil
}
