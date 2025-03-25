package cmd

import (
	"fmt"
	"github.com/joern1811/ai/pkg/core/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	availableAssistants = make(map[string]string)
)

var assistantCallCmd = &cobra.Command{
	Use:                   "call",
	Short:                 "Tool to call an assistant",
	DisableFlagsInUseLine: true,
	Args:                  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.Unmarshal(&appConfig)
		cobra.CheckErr(err)

		assistantName, err := cmd.Flags().GetString("assistantName")
		cobra.CheckErr(err)

		assistantId := getAssistantID(assistantName)

		prompt, err := cmd.Flags().GetString("prompt")
		cobra.CheckErr(err)

		assistantService := service.NewAssistantService(appConfig.OpenAIConfig.OpenAIAuthToken)

		answer, err := assistantService.Call(assistantId, prompt)
		cobra.CheckErr(err)

		fmt.Println(answer)
	},
}

func init() {
	assistantCmd.AddCommand(assistantCallCmd)
	assistants, err := initAssistants()
	if err != nil {
		panic(err)
	}

	assistantCallCmd.Flags().StringP("assistantName", "a", "", "Name of the assistant to call")
	_ = assistantCallCmd.MarkFlagRequired("assistantName")
	_ = assistantCallCmd.RegisterFlagCompletionFunc("assistantName", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		for key := range assistants {
			completions = append(completions, key)
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})
	assistantCallCmd.Flags().StringP("prompt", "p", "", "Prompt for the assistant")
	_ = assistantCallCmd.MarkFlagRequired("prompt")

}

func getAssistantID(assistantName string) string {
	return availableAssistants[assistantName]
}

func initAssistants() (map[string]string, error) {
	if os.Getenv("OPEN_AI_AUTH_TOKEN") == "" {
		return availableAssistants, nil
	}
	assistantService := service.NewAssistantService(os.Getenv("OPEN_AI_AUTH_TOKEN"))
	assistants, err := assistantService.GetAssistants()
	if err != nil {
		return availableAssistants, err
	}

	for _, assistant := range assistants {
		availableAssistants[*assistant.Name] = assistant.ID
	}

	return availableAssistants, nil
}
