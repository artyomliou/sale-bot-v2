package camera

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"fmt"
	"regexp"
	"strings"
)

type PttCrawlerDcSaleAdapter struct {
	Keywords []string `mapstructure:"keywords"`
}

func (a PttCrawlerDcSaleAdapter) GetCrawler() crawlers.Crawler {
	crawler := crawlers.NewPttCrawler()
	crawler.BaseUrl = crawlers.PttBaseUrl
	crawler.Board = "dc_sale"
	crawler.Patterns = []*regexp.Regexp{}

	keywords := strings.Join(a.Keywords, "|")
	crawler.Patterns = append(crawler.Patterns, regexp.MustCompile(fmt.Sprintf("(?i)(%s)", keywords)))

	return crawler
}
