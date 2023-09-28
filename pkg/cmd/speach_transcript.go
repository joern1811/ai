package cmd

import (
	"fmt"
	"github.com/joern1811/ai/pkg/core/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var speachTranscriptCmd = &cobra.Command{
	Use:                   "transcript",
	Short:                 "Tool to transcript audio files",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.Unmarshal(&appConfig)
		cobra.CheckErr(err)

		inputPath := viper.GetString("input")
		if inputPath == "" {
			fmt.Println("Please provide an input file")
			os.Exit(1)
		}

		speachService := service.NewSpeachService(appConfig.OpenAIConfig.OpenAIAuthToken, appConfig.PromptConfig)
		transcript, err := speachService.Transcript(inputPath)
		cobra.CheckErr(err)

		outputPath := viper.GetString("output")
		if outputPath != "" {
			err = os.WriteFile(outputPath, []byte(transcript), 0644)
			cobra.CheckErr(err)
		} else {
			fmt.Println(transcript)
		}
	},
}

func init() {
	speachCmd.AddCommand(speachTranscriptCmd)

	speachTranscriptCmd.Flags().StringP("input", "i", "", "input file")
	speachTranscriptCmd.Flags().StringP("output", "o", "", "output file")
}
