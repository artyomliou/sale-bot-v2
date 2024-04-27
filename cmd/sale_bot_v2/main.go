package main

import (
	"artyomliou/sale-bot-v2/configs"
	"artyomliou/sale-bot-v2/internal/crawlers"
	"artyomliou/sale-bot-v2/internal/db"
	"artyomliou/sale-bot-v2/internal/notification"
	"artyomliou/sale-bot-v2/internal/utils"
	"context"
	"os"
	"os/signal"
	"sync"
	"time"
)

var once sync.Once

func setupTerminableContext() context.Context {
	var ctx context.Context

	once.Do(func() {
		var cancel context.CancelFunc
		var sig chan os.Signal
		ctx, cancel = context.WithCancel(context.Background())
		sig = make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		go func() {
			for {
				<-sig
				cancel()
			}
		}()
	})

	return ctx
}

func main() {
	logger := utils.NewModuleLogger("main")
	ctx := setupTerminableContext()
	cfg := configs.GetConfig()

	// Initialize...
	conn, err := db.NewConnection("db.sqlite")
	if err != nil {
		logger.Fatal(err)
	}

	telegramNotifier, err := notification.NewTelegramNotifier("", cfg.Telegram.BotKey, cfg.Telegram.ChannelId)
	if err != nil {
		logger.Fatal(err)
	}

	adapters, err := cfg.GetAdapters()
	if err != nil {
		logger.Fatal(err)
	}

	// Execute...
	allCrawlers := []crawlers.Crawler{}
	for _, adapter := range adapters {
		allCrawlers = append(allCrawlers, adapter.GetCrawler())
	}

	ticker := time.NewTicker(10 * time.Second)
	for {
		crawledPages := []*crawlers.Page{}
		notifyPages := []*crawlers.Page{}
		select {
		case <-ticker.C:
			for _, crawler := range allCrawlers {
				pages := []*crawlers.Page{}
				crawler.Crawl(ctx, &pages)
				logger.Printf("append %d pages...", len(pages))
				crawledPages = append(crawledPages, pages...)
			}
			for _, page := range crawledPages {
				exists, err := conn.CheckPageExists(page)
				if err != nil {
					logger.Print(err)
					continue
				}
				if !exists {
					logger.Printf("create page %s...", page.ID)
					if err := conn.CreatePage(page); err != nil {
						logger.Print(err)
						continue
					}
					notifyPages = append(notifyPages, page)
				}
			}
			if len(notifyPages) > 0 {
				logger.Printf("will send %d pages to telegram", len(notifyPages))
				for _, page := range notifyPages {
					logger.Printf("sending %s to telegram", page.ID)
					// pick one photo as preview if exists
					photoUrl := ""
					if len(page.PhotoUrls) > 0 {
						photoUrl = page.PhotoUrls[0]
					}
					if err := telegramNotifier.SendMessage(ctx, notification.PageToTelegramMessageText(page), photoUrl); err != nil {
						logger.Print(err)
					}
				}
			}

		case <-ctx.Done():
			ticker.Stop()
			logger.Print("ticker stopped, exit for-loop")
			return
		}
	}

}
