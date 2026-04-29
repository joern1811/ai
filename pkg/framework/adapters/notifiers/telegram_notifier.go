package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/joern1811/ai/pkg/framework/adapters/utils"
)

type TelegramConfig struct {
	ChatID string `mapstructure:"chatID" json:"chatID"`
	Token  string `mapstructure:"token" json:"token"`
}

type TelegramNotifier struct {
	TelegramConfig
}

func NewTelegramNotifier(config TelegramConfig) *TelegramNotifier {
	return &TelegramNotifier{config}
}

func (t TelegramNotifier) Notify(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)
	body := map[string]any{
		"chat_id": t.ChatID,
		"text":    message,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req) //nolint:bodyclose // closed via utils.CloseResource (bodyclose linter does not detect helper-based closures)
	if err != nil {
		return err
	}
	defer utils.CloseResource(resp.Body, &err)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}
