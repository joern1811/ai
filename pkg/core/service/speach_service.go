package service

import (
	"context"
	"github.com/joern1811/ai/pkg/core/domain"
	"github.com/sashabaranov/go-openai"
)

type SpeachService struct {
	client       *openai.Client
	promptConfig domain.PromptConfig
	context      context.Context
}

func NewSpeachService(openaiAuthToken string, promptConfig domain.PromptConfig) *SpeachService {
	return &SpeachService{
		client:       openai.NewClient(openaiAuthToken),
		promptConfig: promptConfig,
		context:      context.Background(),
	}
}

func (s *SpeachService) Transcript(pathToAudioFile string) (string, error) {
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: pathToAudioFile,
	}
	response, err := s.client.CreateTranscription(s.context, req)
	if err != nil {
		return "", err
	}
	return response.Text, nil
}

func (s *SpeachService) Ask(text string) (string, error) {
	messages := make([]openai.ChatCompletionMessage, 0)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}
	content := resp.Choices[0].Message.Content
	return content, nil
}

func (s *SpeachService) SummarizeText(text string) (string, error) {
	question := s.promptConfig.SummarizePrompt + ":\n" + text
	return s.Ask(question)
}

func (s *SpeachService) SummarizeAudio(pathToAudioFile string) (string, error) {
	transcript, err := s.Transcript(pathToAudioFile)
	if err != nil {
		return "", err
	}
	return s.SummarizeText(transcript)
}
