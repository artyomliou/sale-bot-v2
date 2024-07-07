package pc

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"regexp"
)

type PttCrawlerHardwareSaleAdapter struct {
	RegexPatterns []string `mapstructure:"regex_patterns"`
}

func (a PttCrawlerHardwareSaleAdapter) GetCrawler() crawlers.Crawler {
	crawler := crawlers.NewPttCrawler()
	crawler.BaseUrl = crawlers.PttBaseUrl
	crawler.Board = "HardwareSale"
	crawler.Patterns = []*regexp.Regexp{}

	for _, pattern := range a.RegexPatterns {
		crawler.Patterns = append(crawler.Patterns, regexp.MustCompile(pattern))
	}

	return crawler
}
