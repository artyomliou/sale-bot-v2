package crawlers

import (
	"artyomliou/sale-bot-v2/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const FiveNineOneBaseUrl = "https://rent.591.com.tw"
const UserAgent = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0"

type FiveNineOneCrawler struct {
	logger  *log.Logger
	BaseUrl string
	Queries url.Values
}

func NewFiveNineOneCrawler() *FiveNineOneCrawler {
	return &FiveNineOneCrawler{
		logger:  utils.NewModuleLogger("FiveNineOneCrawler"),
		BaseUrl: FiveNineOneBaseUrl,
		Queries: url.Values{},
	}
}

func (c *FiveNineOneCrawler) indexUrl() string {
	indexUrl := c.BaseUrl + "/?" + c.Queries.Encode()
	indexUrl = strings.Replace(indexUrl, "%2C", ",", -1)
	return indexUrl
}

func (c *FiveNineOneCrawler) apiUrl() string {
	apiUrl, _ := url.JoinPath(c.BaseUrl, "/home/search/rsList")
	apiUrl = apiUrl + "?" + c.Queries.Encode()
	apiUrl = strings.Replace(apiUrl, "%2C", ",", -1)
	return apiUrl
}

func (c *FiveNineOneCrawler) pageUrl(postId int) string {
	pageUrl, _ := url.JoinPath(c.BaseUrl, fmt.Sprintf("rent-detail-%d.html", postId))
	return pageUrl
}

func (c *FiveNineOneCrawler) Crawl(ctx context.Context, results *[]*Page) {
	cookies := []*http.Cookie{}
	headers := http.Header{}
	headers["User-Agent"] = []string{UserAgent}

	c.crawlHtmlForCsrfTokenAndCookie(ctx, &headers, &cookies)

	// Prepare for API call
	headers["device"] = []string{"pc"}
	headers["Referer"] = []string{c.indexUrl()}
	headers["Accpet"] = []string{"application/json, text/javascript, */*; q=0.01"}
	headers["Accept-Language"] = []string{"en-US,en;q=0.5"}
	headers["DNT"] = []string{"1"}
	headers["X-Requested-With"] = []string{"XMLHttpRequest"}

	// Pagination control
	page := 1
	const PerPageRows = 30
	for {
		totalRows := c.crawlApi(ctx, &headers, &cookies, results)
		if totalRows == 0 {
			break
		}
		if totalRows-page*PerPageRows <= 0 {
			break
		}
		page++
		c.Queries.Set("firstRow", fmt.Sprintf("%d", (page-1)*PerPageRows))
		c.Queries.Set("totalRows", fmt.Sprintf("%d", totalRows))
	}
}

func (c *FiveNineOneCrawler) crawlHtmlForCsrfTokenAndCookie(ctx context.Context, headers *http.Header, cookies *[]*http.Cookie) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.indexUrl(), nil)
	if err != nil {
		c.logger.Printf("failed to new request %s", err)
		return
	}
	req.Header = *headers

	// Debug
	// bytes, _ := httputil.DumpRequestOut(req, false)
	// c.logger.Printf("%s\n", bytes)

	// Send
	c.logger.Printf("crawling %s...", c.indexUrl())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Printf("failed to request a target: %v", err)
		return
	}
	c.logger.Printf("status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		c.logger.Printf("status code is not 200")
		return
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		c.logger.Printf("Response of target is not HTML: %v", err)
		return
	}

	// headers: X-CSRF-TOKEN
	htmlNode, err := html.Parse(resp.Body)
	if err != nil {
		c.logger.Printf("failed to parse response as HTML node: %v", err)
		return
	}

	doc := goquery.NewDocumentFromNode(htmlNode)
	doc.Find("meta[name=csrf-token]").Each(func(i int, s *goquery.Selection) {
		csrfToken, ok := s.Attr("content")
		if !ok {
			c.logger.Println("failed to get meta[name=csrf-token]")
			return
		}
		(*headers)["X-CSRF-TOKEN"] = []string{csrfToken}
	})

	// headers: deviceid
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "T591_TOKEN" {
			(*headers)["deviceid"] = []string{cookie.Value}
		}
	}

	// Keep cookies for following API call
	*cookies = resp.Cookies()
	for _, cookie := range *cookies {
		if cookie.Name == "urlJumpIp" {
			cookie.Value = "3"
		}
	}
}

func (c *FiveNineOneCrawler) crawlApi(ctx context.Context, headers *http.Header, cookies *[]*http.Cookie, results *[]*Page) int {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiUrl(), nil)
	if err != nil {
		c.logger.Printf("failed to new request %s", err)
		return 0
	}
	req.Header = *headers
	for _, cookie := range *cookies {
		req.AddCookie(cookie)
	}

	// Debug
	// reqBytes, _ := httputil.DumpRequestOut(req, false)
	// c.logger.Printf("%s\n", reqBytes)

	// Send
	c.logger.Printf("crawling %s...", c.apiUrl())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Printf("failed to request a target: %v", err)
		return 0
	}
	c.logger.Printf("status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		c.logger.Print("status code is not 200")
		return 0
	}

	// Debug
	// respBytes, _ := httputil.DumpResponse(resp, true)
	// c.logger.Printf("%s\n", respBytes)

	// Update cookies for following API call
	for _, prevCookie := range *cookies {
		for _, newCookie := range resp.Cookies() {
			if prevCookie.Name == newCookie.Name {
				prevCookie.Value = newCookie.Value
			}
		}
	}

	// Parse Response
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Printf("failed to read body: %v", err)
		return 0
	}

	jsonResp := FiveNineOneResponse{}
	if err := json.Unmarshal(b, &jsonResp); err != nil {
		c.logger.Printf("failed to parse body to json: %v", err)
		return 0
	}

	for _, obj := range jsonResp.Data.Data {
		*results = append(*results, &Page{
			ID:    fmt.Sprintf("591-%d", obj.PostId),
			Title: obj.Title,
			Link:  c.pageUrl(obj.PostId),
		})
	}

	if jsonResp.Data.Records != "" {
		totalRows, err := strconv.ParseUint(jsonResp.Data.Records, 10, 32)
		if err != nil {
			c.logger.Printf("failed to parse .Records as uint: %v", err)
			return 0
		}
		return int(totalRows)
	}

	return 0
}

type FiveNineOneResponse struct {
	Data struct {
		Records string `json:"records"`
		Data    []struct {
			Title       string   `json:"title"`
			PostId      int      `json:"post_id"`
			KindName    string   `json:"kind_name"`
			RoomStr     string   `json:"room_str"`
			FloorStr    string   `json:"floor_str"`
			Price       string   `json:"price"`
			PhotoList   []string `json:"photo_list"`
			SectionName string   `json:"section_name"`
			StreetName  string   `json:"street_name"`
			Location    string   `json:"location"`
			RentTag     []struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"rent_tag"`
		} `json:"data"`
	} `json:"data"`
}
