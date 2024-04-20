package notification

import (
	"artyomliou/sale-bot-v2/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type telegramNotifier struct {
	logger    *log.Logger
	url       string
	channelId int
}

func NewTelegramNotifier(host, botKey string, channelID int) (*telegramNotifier, error) {
	if host == "" {
		host = "https://api.telegram.org"
	}
	if botKey == "" {
		return nil, errors.New("bot key cannot be empty")
	}
	if channelID == 0 {
		return nil, errors.New("channelID cannot be empty")
	}
	return &telegramNotifier{
		logger:    utils.NewModuleLogger("telegram"),
		url:       fmt.Sprintf("%s/bot%s/sendMessage", host, botKey),
		channelId: channelID,
	}, nil
}

func (n *telegramNotifier) SendMessage(ctx context.Context, text string) error {
	type payload struct {
		ChatId                int    `json:"chat_id"`
		Text                  string `json:"text"`
		ParseMode             string `json:"parse_mode"`
		DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	}

	errorWrapper := func(err error) error {
		return fmt.Errorf("send message error: %s", err.Error())
	}

	var buf bytes.Buffer
	body := payload{
		ChatId:                int(n.channelId),
		Text:                  text,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
	}
	if err := json.NewEncoder(&buf).Encode(&body); err != nil {
		return errorWrapper(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.url, &buf)
	if err != nil {
		return errorWrapper(err)
	}
	req.Header.Set("Content-Type", "application/json")

	n.logger.Printf("sending message...")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errorWrapper(err)
	}
	n.logger.Printf("status code: %d", resp.StatusCode)
	return nil
}
