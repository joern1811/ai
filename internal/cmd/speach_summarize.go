package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/joern1811/ai/internal/core/service"
)

var speachSummarizeCmd = &cobra.Command{
	Use:                   "summarize",
	Short:                 "Tool to summarize text- or audio-files",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.Unmarshal(&appConfig)
		cobra.CheckErr(err)

		speachService := service.NewSpeachService(appConfig.OpenAIConfig.OpenAIAuthToken, appConfig.PromptConfig)

		inputPath := cmd.Flags().Lookup("input").Value.String()
		if inputPath == "" {
			fmt.Println("Please provide an input file")
			os.Exit(1)
		}
		inputPath = filepath.Clean(inputPath)

		file, err := os.Open(inputPath) //nolint:gosec // inputPath comes from a CLI flag
		cobra.CheckErr(err)

		// We only have to pass the file header = first 261 bytes
		head := make([]byte, 261)
		_, err = file.Read(head)
		cobra.CheckErr(err)

		var summary string
		if filetype.IsAudio(head) {
			summary, err = speachService.SummarizeAudio(inputPath)
		} else {
			text, readErr := os.ReadFile(inputPath) //nolint:gosec // inputPath comes from a CLI flag
			cobra.CheckErr(readErr)
			summary, err = speachService.SummarizeText(string(text))
		}
		cobra.CheckErr(err)

		outputPath := cmd.Flags().Lookup("output").Value.String()
		if outputPath != "" {
			err = os.WriteFile(filepath.Clean(outputPath), []byte(summary), 0o600)
			cobra.CheckErr(err)
		} else {
			fmt.Println(summary)
		}
	},
}

func init() {
	speachCmd.AddCommand(speachSummarizeCmd)

	speachSummarizeCmd.Flags().StringP("input", "i", "", "input file")
	speachSummarizeCmd.Flags().StringP("output", "o", "", "output file")
}
