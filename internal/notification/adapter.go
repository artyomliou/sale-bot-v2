package notification

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"fmt"
)

func PagesToHtml(pages []*crawlers.Page) string {
	text := ""
	for _, page := range pages {
		text += fmt.Sprintf("<a href=\"%s\">%s</a>\n", page.Link, page.NotificationTitle)
	}
	return text
}
