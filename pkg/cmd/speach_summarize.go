package cmd

import (
	"fmt"
	"github.com/h2non/filetype"
	"github.com/joern1811/ai/pkg/core/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var speachSummarizeCmd = &cobra.Command{
	Use:                   "summarize",
	Short:                 "Tool to summarize text- or audio-files",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.Unmarshal(&appConfig)
		cobra.CheckErr(err)

		speachService := service.NewSpeachService(appConfig.OpenAIConfig.OpenAIAuthToken, appConfig.PromptConfig)

		inputPath := viper.GetString("input")
		if inputPath == "" {
			fmt.Println("Please provide an input file")
			os.Exit(1)
		}

		file, err := os.Open(inputPath)
		cobra.CheckErr(err)

		// We only have to pass the file header = first 261 bytes
		head := make([]byte, 261)
		_, err = file.Read(head)
		cobra.CheckErr(err)

		var summary string
		if filetype.IsAudio(head) {
			summary, err = speachService.SummarizeAudio(inputPath)
		} else {
			text, err := os.ReadFile(inputPath)
			cobra.CheckErr(err)
			summary, err = speachService.SummarizeText(string(text))
		}
		outputPath := viper.GetString("output")
		if outputPath != "" {
			err = os.WriteFile(outputPath, []byte(summary), 0644)
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
