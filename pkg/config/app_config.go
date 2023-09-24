package config

import "github.com/joern1811/ai/pkg/core/domain"

type AppConfig struct {
	OpenAIConfig domain.OpenAIConfig `mapstructure:"openAIConfig"`
	PromptConfig domain.PromptConfig `mapstructure:"promptConfig"`
}
