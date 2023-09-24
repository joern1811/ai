package cmd

import (
	"fmt"
	"github.com/joern1811/ai/pkg/core/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var speachAskCmd = &cobra.Command{
	Use:                   "ask QUESTION",
	Short:                 "Tool to ask questions",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.Unmarshal(&appConfig)
		cobra.CheckErr(err)

		speachService := service.NewSpeachService(appConfig.OpenAIConfig.OpenAIAuthToken, appConfig.PromptConfig)
		response, err := speachService.Ask(args[0])
		cobra.CheckErr(err)

		fmt.Println(response)
	},
}

func init() {
	speachCmd.AddCommand(speachAskCmd)
}
