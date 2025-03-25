package cmd

import "github.com/spf13/cobra"

var speachCmd = &cobra.Command{
	Use:   "speach",
	Short: "Some tools in the context of speach",
}

func init() {
	rootCmd.AddCommand(speachCmd)
}
