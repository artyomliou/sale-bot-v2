package notification_test

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"artyomliou/sale-bot-v2/internal/notification"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTelegram(t *testing.T) {
	// mock server
	didSentToMockEndpoint := false

	router := gin.Default()
	router.Any("/*proxyPath", func(ctx *gin.Context) {
		didSentToMockEndpoint = true
		ctx.AbortWithStatus(200)
	})
	svr := httptest.NewServer(router)
	defer svr.Close()

	// prepare
	notifier, err := notification.NewTelegramNotifier(svr.URL, "12312312323", 123123)
	if err != nil {
		t.Fatal(err)
	}

	// test
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	text := notification.PagesToHtml([]*crawlers.Page{
		{
			Title: "Testing title",
			Link:  "https://www.ptt.cc/bbs/DC_SALE/M.1691245105.A.729.html",
		},
	})
	err = notifier.SendMessage(ctx, text)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, true, didSentToMockEndpoint)
}
