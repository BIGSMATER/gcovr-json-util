package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zjy-dev/gcovr-json-util/pkg/gcovr"
)

var (
	baseFile string
	newFile  string
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare two gcovr JSON reports and show coverage increases",
	Long: `Compare a base gcovr JSON report with a new report to identify 
functions with increased line coverage. The tool will report:
- Which functions have coverage increases
- How many lines were newly covered
- Total lines in each function
- Demangled function names for readability`,
	RunE: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)

	diffCmd.Flags().StringVarP(&baseFile, "base", "b", "", "Base gcovr JSON report file (required)")
	diffCmd.Flags().StringVarP(&newFile, "new", "n", "", "New gcovr JSON report file (required)")

	diffCmd.MarkFlagRequired("base")
	diffCmd.MarkFlagRequired("new")
}

func runDiff(cmd *cobra.Command, args []string) error {
	// Parse base report
	fmt.Printf("Reading base report: %s\n", baseFile)
	baseReport, err := gcovr.ParseReport(baseFile)
	if err != nil {
		return fmt.Errorf("failed to parse base report: %w", err)
	}

	// Parse new report
	fmt.Printf("Reading new report: %s\n", newFile)
	newReport, err := gcovr.ParseReport(newFile)
	if err != nil {
		return fmt.Errorf("failed to parse new report: %w", err)
	}

	// Compute coverage increase
	fmt.Println("\nComputing coverage increases...\n")
	report, err := gcovr.ComputeCoverageIncrease(baseReport, newReport)
	if err != nil {
		return fmt.Errorf("failed to compute coverage increase: %w", err)
	}

	// Display results
	output := gcovr.FormatReport(report)
	fmt.Print(output)

	return nil
}
