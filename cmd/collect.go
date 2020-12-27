package cmd

import (
	"errors"
	"fmt"

	"github.com/mdb/gossboss"
	"github.com/spf13/cobra"
)

var (
	// collectCmd is the cobra.Command defining the "gossboss collect" action.
	collectCmd = &cobra.Command{
		Use:   "collect",
		Short: "Collect and report goss test results",
		Long:  "Collect and report goss test results from multiple goss servers' '/healthz' endpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := gossboss.NewClient()
			servers, err := cmd.Flags().GetStringSlice("servers")
			if err != nil {
				return err
			}

			resps := c.CollectHealthzs(servers)

			hasErrors := false
			hasFailed := false
			for _, resp := range resps {
				if resp.Error != nil {
					hasErrors = true
					pFailure(resp.URL)
					fmt.Println(fmt.Sprintf(" \tError: %v", resp.Error.Error()))
					continue
				}

				if resp.Result.Summary.Failed > 0 {
					hasFailed = true
					pFailure(resp.URL)
					continue
				}

				fmt.Println(fmt.Sprintf(" \xE2\x9C\x94 %s", resp.URL))
			}

			// TODO: handle scenarios with both errors and failures
			if hasErrors {
				return errors.New("Goss test collection error")
			}

			if hasFailed {
				return errors.New("Goss test failed")
			}

			return nil
		},
	}
)

func init() {
	collectCmd.Flags().StringSliceP("servers", "s", []string{}, "A comma-separated list of goss servers")
	collectCmd.MarkFlagRequired("servers")
	rootCmd.AddCommand(collectCmd)
}

func pFailure(url string) {
	fmt.Println(fmt.Sprintf(" \u2718 %s", url))
}
