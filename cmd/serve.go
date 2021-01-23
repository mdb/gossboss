package cmd

import (
	"github.com/mdb/gossboss"
	"github.com/spf13/cobra"
)

var (
	// serveCmd is the cobra.Command defining the "gossboss serve" action.
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Collect and report goss test results via a web server JSON endpoint",
		Long:  "Collect and report goss test results from multiple goss servers' '/healthz' endpoints via a web server JSON endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			servers, err := cmd.Flags().GetStringSlice("servers")
			if err != nil {
				return err
			}

			port, err := cmd.Flags().GetString("port")
			if err != nil {
				return err
			}

			_ = gossboss.NewServer(":"+port, servers)

			return nil
		},
	}
)

func init() {
	serveCmd.Flags().StringSliceP("servers", "s", []string{}, "A comma-separated list of goss servers from which to collect test results")
	serveCmd.MarkFlagRequired("servers")
	serveCmd.Flags().StringP("port", "p", "8081", "The port on which to run the gossboss server")
	rootCmd.AddCommand(serveCmd)
}
