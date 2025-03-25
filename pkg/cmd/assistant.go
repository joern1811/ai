package cmd

import "github.com/spf13/cobra"

var assistantCmd = &cobra.Command{
	Use:   "assistant",
	Short: "Some tools in context to assistants",
}

func init() {
	rootCmd.AddCommand(assistantCmd)
}
