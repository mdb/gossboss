package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "gossboss",
	Short:        "gossboss collects goss test results",
	Long:         "A tool for collecting goss test results from multiple goss servers' '/healthz' endpoints",
	Version:      "0.0.1",
	SilenceUsage: true,
}

// Execute executes the gossboss command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
