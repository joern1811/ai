package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joern1811/ai/pkg/framework/adapters/utils"
	"net/http"
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
	body := map[string]interface{}{
		"chat_id": t.ChatID,
		"text":    message,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	defer utils.CloseResource(resp.Body, &err)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}
