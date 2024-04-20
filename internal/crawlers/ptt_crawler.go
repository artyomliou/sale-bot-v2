package crawlers

import (
	"artyomliou/sale-bot-v2/internal/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type PttCrawler struct {
	logger   *log.Logger
	Url      string
	Patterns []*regexp.Regexp
}

func NewPttCrawler() *PttCrawler {
	return &PttCrawler{
		logger: utils.NewModuleLogger("PttCrawler"),
	}
}

func PttBoardUrl(board string) string {
	return fmt.Sprintf("https://www.ptt.cc/bbs/%s/index.html", board)
}

func (c *PttCrawler) Crawl(ctx context.Context, results *[]*Page) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Url, nil)
	if err != nil {
		c.logger.Printf("failed to init a request: %v", err)
		return
	}

	c.logger.Printf("crawling %s...", c.Url)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Printf("failed to request a target: %v", err)
		return
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		c.logger.Printf("Response of target is not HTML: %v", err)
		return
	}

	htmlNode, err := html.Parse(resp.Body)
	if err != nil {
		c.logger.Printf("failed to parse response as HTML node: %v", err)
		return
	}

	crawledPages := []*Page{}
	doc := goquery.NewDocumentFromNode(htmlNode)
	doc.Find(".r-list-container .r-ent .title").Each(func(i int, s *goquery.Selection) {
		link, ok := s.Find("a").Attr("href")
		if !ok {
			c.logger.Println("failed to get a:href")
			return
		}
		title := s.Find("a").Text()
		crawledPages = append(crawledPages, &Page{
			ID:    link,
			Link:  "https://www.ptt.cc" + link,
			Title: title,
		})
	})

	matchedPages := []*Page{}
	for _, page := range crawledPages {
		matchedCount := 0
		for _, pattern := range c.Patterns {
			if pattern.MatchString(page.Title) {
				matchedCount++
			}
		}
		if matchedCount == len(c.Patterns) {
			matchedPages = append(matchedPages, page)
		}
	}

	*results = append(*results, matchedPages...)
}