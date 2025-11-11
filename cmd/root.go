package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gcovr-util",
	Short: "A utility tool for processing gcovr JSON reports",
	Long: `gcovr-util is a CLI tool for analyzing and comparing gcovr JSON coverage reports.
It provides functionality to measure coverage increases between test runs and
identify which functions have improved coverage.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
}
