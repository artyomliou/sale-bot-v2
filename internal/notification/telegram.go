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
	"net/url"
)

const TelegramBaseUrl = "https://api.telegram.org"

type telegramNotifier struct {
	logger    *log.Logger
	baseUrl   string
	botKey    string
	channelId int
}

func NewTelegramNotifier(baseUrl, botKey string, channelID int) (*telegramNotifier, error) {
	if baseUrl == "" {
		baseUrl = TelegramBaseUrl
	}
	if botKey == "" {
		return nil, errors.New("bot key cannot be empty")
	}
	if channelID == 0 {
		return nil, errors.New("channelID cannot be empty")
	}
	return &telegramNotifier{
		logger:    utils.NewModuleLogger("telegram"),
		baseUrl:   baseUrl,
		botKey:    botKey,
		channelId: channelID,
	}, nil
}

func (n *telegramNotifier) sendMessageUrl() string {
	u, _ := url.JoinPath(n.baseUrl, fmt.Sprintf("bot%s", n.botKey), "sendMessage")
	return u
}

func (n *telegramNotifier) SendMessage(ctx context.Context, text string, photoUrl string) error {
	errorWrapper := func(err error) error {
		return fmt.Errorf("send message error: %s", err.Error())
	}

	var buf bytes.Buffer
	body := sendMessagePayload{
		ChatId:              int(n.channelId),
		Text:                text,
		ParseMode:           "HTML",
		DisableNotification: false,
		LinkPreviewOptions: linkPreviewOptions{
			IsDisabled: photoUrl == "",
			Url:        photoUrl,
		},
	}
	if err := json.NewEncoder(&buf).Encode(&body); err != nil {
		return errorWrapper(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.sendMessageUrl(), &buf)
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

type linkPreviewOptions struct {
	IsDisabled bool   `json:"is_disabled"`
	Url        string `json:"url"`
}
type sendMessagePayload struct {
	ChatId              int                `json:"chat_id"`
	Text                string             `json:"text"`
	ParseMode           string             `json:"parse_mode"`
	DisableNotification bool               `json:"disable_notification"`
	LinkPreviewOptions  linkPreviewOptions `json:"link_preview_options"`
}
