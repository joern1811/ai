package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"log"
	"time"
)

type AssistantService struct {
	client *openai.Client
	ctx    context.Context
}

func NewAssistantService(authToken string) *AssistantService {
	return &AssistantService{
		client: openai.NewClient(authToken),
		ctx:    context.Background(),
	}
}

func (a *AssistantService) Call(assistantId string, prompt string) (string, error) {
	run, err := a.client.CreateThreadAndRun(a.ctx, openai.CreateThreadAndRunRequest{
		RunRequest: openai.RunRequest{AssistantID: assistantId},
		Thread: openai.ThreadRequest{
			Messages: []openai.ThreadMessage{{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			}},
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "error on creating thread and run")
	}

	for run.Status == openai.RunStatusQueued || run.Status == openai.RunStatusInProgress {
		run, err = a.client.RetrieveRun(a.ctx, run.ThreadID, run.ID)
		if err != nil {
			return "", errors.Wrap(err, "error on retrieving run")
		}
		time.Sleep(100 * time.Millisecond)
	}
	if run.Status != openai.RunStatusCompleted {
		return "", errors.New("run failed")
	}

	numMessages := 1
	messages, err := a.client.ListMessage(a.ctx, run.ThreadID, &numMessages, nil, nil, nil, nil)
	if err != nil {
		return "", errors.Wrap(err, "error on listing messages")
	}

	return messages.Messages[0].Content[0].Text.Value, nil
}

func (a *AssistantService) GetAssistants() ([]openai.Assistant, error) {
	assistants, err := a.client.ListAssistants(a.ctx, nil, nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	return assistants.Assistants, nil
}

func (a *AssistantService) DetermineAssistantId(assistantName string) (string, error) {
	assistants, err := a.client.ListAssistants(a.ctx, nil, nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, assistant := range assistants.Assistants {
		if *assistant.Name == assistantName {
			return assistant.ID, nil
		}
	}
	return "", nil
}
