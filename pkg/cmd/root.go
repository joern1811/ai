package cmd

import (
	"fmt"
	"github.com/joern1811/ai/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	baseAppName = "ai"
)

var (
	cfgFile   string
	appConfig *config.AppConfig
)

var rootCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI is a tool to use OpenAI's GPT-3 API",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", fmt.Sprintf("config file (default $HOME/.%v)", baseAppName))
}

func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(fmt.Sprintf(".%v", baseAppName))
	}

	viper.AutomaticEnv() // read in Environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("cannot read config file: %v because: %v", viper.ConfigFileUsed(), err)
	}
}
