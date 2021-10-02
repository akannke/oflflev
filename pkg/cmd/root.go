package cmd

import (
	"fmt"
	"log"

	"github.com/akannke/oflflev/pkg/oflflev"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "oflflev",
	Short: "oflflev is app to solve open face chinese poker",

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Hello")
		numDealt, err := cmd.Flags().GetInt("dealt")
		if err != nil {
			return err
		}

		iteration, err := cmd.Flags().GetInt("iteration")
		if err != nil {
			return err
		}

		oflflev.Solve(iteration, numDealt)

		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().IntP("iteration", "i", 10000, "number of iterations")
	RootCmd.PersistentFlags().IntP("dealt", "d", 14, "number of cards dealt")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("cannot execute command")
	}
}
