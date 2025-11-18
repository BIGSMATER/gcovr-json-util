package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zjy-dev/gcovr-json-util/v2/pkg/gcovr"
)

var (
	baseFile   string
	newFile    string
	filterFile string
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
- Demangled function names for readability

Optionally, you can specify a filter configuration file to only track
specific files and functions defined in the targets.`,
	RunE: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)

	diffCmd.Flags().StringVarP(&baseFile, "base", "b", "", "Base gcovr JSON report file (required)")
	diffCmd.Flags().StringVarP(&newFile, "new", "n", "", "New gcovr JSON report file (required)")
	diffCmd.Flags().StringVarP(&filterFile, "filter", "f", "", "Filter config file (YAML) to specify target files and functions")

	diffCmd.MarkFlagRequired("base")
	diffCmd.MarkFlagRequired("new")
}

func runDiff(cmd *cobra.Command, args []string) error {
	// Parse filter config if provided
	var filterConfig *gcovr.FilterConfig
	if filterFile != "" {
		fmt.Printf("Reading filter config: %s\n", filterFile)
		var err error
		filterConfig, err = gcovr.ParseFilterConfig(filterFile)
		if err != nil {
			return fmt.Errorf("failed to parse filter config: %w", err)
		}
		fmt.Printf("Filtering enabled: tracking %d file(s)\n", len(filterConfig.Targets))
	}

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

	// Apply filter if provided
	if filterConfig != nil {
		fmt.Println("Applying filters...")
		baseReport = gcovr.ApplyFilter(baseReport, filterConfig)
		newReport = gcovr.ApplyFilter(newReport, filterConfig)
	}

	// Compute coverage increase
	fmt.Println("Computing coverage increases...")
	report, err := gcovr.ComputeCoverageIncrease(baseReport, newReport)
	if err != nil {
		return fmt.Errorf("failed to compute coverage increase: %w", err)
	}

	// Display results
	output := gcovr.FormatReport(report)
	fmt.Print(output)

	return nil
}
