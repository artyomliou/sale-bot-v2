package crawlers

import "context"

type Crawler interface {
	Crawl(context.Context, *[]*Page)
}

type CrawlerAdapter interface {
	GetCrawler() Crawler
}

type Page struct {
	ID                string
	Link              string
	Title             string
	NotificationTitle string
	PhotoUrls         []string
}
