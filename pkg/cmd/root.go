package cmd

import (
	"fmt"
	"github.com/joern1811/ai/pkg/config"
	_ "github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", fmt.Sprintf("config file (default $HOME/.config/.%v/config)", baseAppName))
}

func initConfig() {

	if cfgFile != "" {
		// Verwende Config-Datei vom Flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Bestimme Config-Verzeichnis nach XDG Base Directory
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			// Fallback: ~/.config
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			configHome = filepath.Join(home, ".config")
		}

		// Unterverzeichnis für die Anwendung in XDG_CONFIG_HOME
		configDir := filepath.Join(configHome, baseAppName)

		// Erstelle den Unterordner, falls er nicht existiert
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			err = os.MkdirAll(configDir, 0755)
			cobra.CheckErr(err)
		}

		// Viper soll in diesem Verzeichnis nach der Datei "config.json" suchen
		viper.AddConfigPath(configDir)
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	// Erlaube Environment-Variablen (ersetze '.' durch '_')
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()

	// Falls eine Config-Datei gefunden wird, lese sie ein
	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Fehler beim Einlesen der Config-Datei: %s\n", viper.ConfigFileUsed())
	}
}
