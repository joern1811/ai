package cmd

import (
	"fmt"
	"github.com/joern1811/ai/pkg/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		long, err := cmd.Flags().GetBool("long")
		if err != nil {
			return err
		}
		if long {
			fmt.Printf("%#v", version.Get())
		} else {
			fmt.Println(version.Get().Version)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().Bool("long", false, "output long version")
}
