package main

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"artyomliou/sale-bot-v2/internal/db"
	"artyomliou/sale-bot-v2/internal/notification"
	"artyomliou/sale-bot-v2/internal/use_cases/camera"
	"artyomliou/sale-bot-v2/internal/use_cases/renthouse"
	"artyomliou/sale-bot-v2/internal/utils"
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"time"
)

var telegramBotKey string
var telegramChannelId int
var once sync.Once

func init() {
	flag.StringVar(&telegramBotKey, "tgBotKey", "", "telegram bot key")
	flag.IntVar(&telegramChannelId, "tgChannelId", 0, "telegram channel id")
}

func setupTerminableContext() context.Context {
	var ctx context.Context

	once.Do(func() {
		var cancel context.CancelFunc
		var sig chan os.Signal
		ctx, cancel = context.WithCancel(context.Background())
		sig = make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		go func() {
			<-sig
			cancel()
		}()
	})

	return ctx
}

func main() {
	flag.Parse()
	if telegramBotKey == "" || telegramChannelId == 0 {
		flag.Usage()
		os.Exit(1)
	}

	logger := utils.NewModuleLogger("main")
	ctx := setupTerminableContext()

	// Initialize...
	conn, err := db.NewConnection("db.sqlite")
	if err != nil {
		logger.Fatal(err)
	}

	telegramNotifier, err := notification.NewTelegramNotifier("", telegramBotKey, telegramChannelId)
	if err != nil {
		logger.Fatal(err)
	}

	// All crawling stuff
	adapters := []crawlers.CrawlerAdapter{
		camera.PttCrawlerDcSaleAdapter{
			Keywords: []string{
				"canon",
				"RF",
				"50",
				"f1.8",
			},
		},
		renthouse.PttCrawlerRentApartAdapter{
			Cities: []renthouse.City{
				renthouse.NewTaipei,
			},
			Districts: []renthouse.District{
				renthouse.Banqiao,
				renthouse.Xinzhuang,
				renthouse.Sanchong,
				renthouse.Luzhou,
				renthouse.Zhonghe,
			},
			Room: renthouse.TwoRooms,
		},
		renthouse.FiveNineOneAdapter{
			City: renthouse.NewTaipei,
			Districts: []renthouse.District{
				renthouse.Banqiao,
				renthouse.Xinzhuang,
				renthouse.Sanchong,
				renthouse.Luzhou,
				renthouse.Zhonghe,
			},
			Kind: renthouse.Apartment,
			Room: []renthouse.RoomCount{
				renthouse.TwoRooms,
				renthouse.ThreeRooms,
			},
			PriceRange: &renthouse.Range{
				Min: 20000,
				Max: 30000,
			},
			FloorRange: &renthouse.Range{
				Min: 1,
				Max: 10,
			},
			Options: []renthouse.Option{
				renthouse.AirConditioner,
				renthouse.Balcony,
				renthouse.NoRoofTop,
				renthouse.Refrigerator,
				renthouse.WashingMachine,
				renthouse.WaterHeater,
			},
		},
	}

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
