package crawlers

import (
	"artyomliou/sale-bot-v2/internal/utils"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const PttBaseUrl = "https://www.ptt.cc"

type PttCrawler struct {
	logger   *log.Logger
	BaseUrl  string
	Board    string
	Patterns []*regexp.Regexp
}

func NewPttCrawler() *PttCrawler {
	return &PttCrawler{
		logger: utils.NewModuleLogger("PttCrawler"),
	}
}

func (c *PttCrawler) boardIndexUrl() string {
	res, _ := url.JoinPath(c.BaseUrl, "bbs", c.Board, "index.html")
	return res
}

func (c *PttCrawler) Crawl(ctx context.Context, results *[]*Page) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.boardIndexUrl(), nil)
	if err != nil {
		c.logger.Printf("failed to init a request: %v", err)
		return
	}

	c.logger.Printf("crawling %s...", c.boardIndexUrl())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Printf("failed to request a target: %v", err)
		return
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		c.logger.Print("Response of target is not HTML")
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
			ID:                fmt.Sprintf("ptt-%s", link),
			Link:              PttBaseUrl + link,
			Title:             title,
			NotificationTitle: fmt.Sprintf("ptt %s", title),
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

	for _, page := range matchedPages {
		if err := c.crawlSpecificPageForMoreData(ctx, page); err != nil {
			c.logger.Printf("failed to crawl specific page: %s", err)
		}
	}

	*results = append(*results, matchedPages...)
}

func (c *PttCrawler) crawlSpecificPageForMoreData(ctx context.Context, page *Page) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, page.Link, nil)
	if err != nil {
		return err
	}

	c.logger.Printf("crawling %s...", page.Link)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return errors.New("response of target is not HTML")
	}

	// parse photos
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	pattern, err := regexp.Compile("(?i)https://i.imgur.com/[0-9A-Z]+.(?:jpg|jpeg|png)")
	if err != nil {
		return err
	}
	matched := pattern.FindAllString(string(bytes), -1)
	if len(matched) > 0 {
		page.PhotoUrls = matched
	}
	return nil
}
