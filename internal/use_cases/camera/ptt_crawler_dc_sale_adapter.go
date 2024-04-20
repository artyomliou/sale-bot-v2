package camera

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"fmt"
	"regexp"
	"strings"
)

type PttCrawlerDcSaleAdapter struct {
	Keywords []string
}

func (a PttCrawlerDcSaleAdapter) GetCrawler() crawlers.Crawler {
	crawler := crawlers.NewPttCrawler()
	crawler.Url = "https://www.ptt.cc/bbs/dc_sale/index.html"
	crawler.Patterns = []*regexp.Regexp{}

	keywords := strings.Join(a.Keywords, "|")
	crawler.Patterns = append(crawler.Patterns, regexp.MustCompile(fmt.Sprintf("(?i)(%s)", keywords)))

	return crawler
}
