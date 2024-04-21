package notification

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"fmt"
)

func PageToTelegramMessageText(page *crawlers.Page) string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>\n", page.Link, page.NotificationTitle)
}
