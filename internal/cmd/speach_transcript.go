package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/joern1811/ai/internal/core/service"
)

var speachTranscriptCmd = &cobra.Command{
	Use:                   "transcript",
	Short:                 "Tool to transcript audio files",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.Unmarshal(&appConfig)
		cobra.CheckErr(err)

		inputPath := cmd.Flags().Lookup("input").Value.String()
		if inputPath == "" {
			fmt.Println("Please provide an input file")
			os.Exit(1)
		}

		speachService := service.NewSpeachService(appConfig.OpenAIConfig.OpenAIAuthToken, appConfig.PromptConfig)
		transcript, err := speachService.Transcript(inputPath)
		cobra.CheckErr(err)

		outputPath := cmd.Flags().Lookup("output").Value.String()
		if outputPath != "" {
			err = os.WriteFile(filepath.Clean(outputPath), []byte(transcript), 0o600)
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
