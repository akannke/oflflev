package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(regularCmd)
}

var regularCmd = &cobra.Command{
	Use:   "regular",
	Short: "OFC/P rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		numDealt, err := cmd.Flags().GetInt("dealt")
		if err != nil {
			return err
		}

		iteration, err := cmd.Flags().GetInt("iteration")
		if err != nil {
			return err
		}
		_ = numDealt
		_ = iteration

		return nil
	},
}
