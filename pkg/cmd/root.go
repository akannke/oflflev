package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "oflflev",
	Short: "oflflev is app to solve open face chinese poker",

	RunE: func(cmd *cobra.Command, args []string) error {
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
