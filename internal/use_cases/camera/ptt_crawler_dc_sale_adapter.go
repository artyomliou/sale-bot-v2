package camera

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"regexp"
)

type PttCrawlerDcSaleAdapter struct {
	RegexPatterns []string `mapstructure:"regex_patterns"`
}

func (a PttCrawlerDcSaleAdapter) GetCrawler() crawlers.Crawler {
	crawler := crawlers.NewPttCrawler()
	crawler.BaseUrl = crawlers.PttBaseUrl
	crawler.Board = "dc_sale"
	crawler.Patterns = []*regexp.Regexp{}

	for _, pattern := range a.RegexPatterns {
		crawler.Patterns = append(crawler.Patterns, regexp.MustCompile(pattern))
	}

	return crawler
}
