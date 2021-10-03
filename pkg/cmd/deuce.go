package cmd

import (
	"github.com/akannke/oflflev/pkg/deuce"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(deuceCmd)
}

var deuceCmd = &cobra.Command{
	Use:   "deuce",
	Short: "OFC/P 2-7 rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		numDealt, err := cmd.Flags().GetInt("dealt")
		if err != nil {
			return err
		}

		iteration, err := cmd.Flags().GetInt("iteration")
		if err != nil {
			return err
		}
		deuce.Solve(iteration, numDealt)

		return nil
	},
}
