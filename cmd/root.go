package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "gossboss",
	Short:        "gossboss collects goss test results",
	Long:         "A tool for collecting goss test results from multiple goss servers' '/healthz' endpoints",
	SilenceUsage: true,
}

// Execute executes the gossboss command.
func Execute(version string) {
	rootCmd.Version = version

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
