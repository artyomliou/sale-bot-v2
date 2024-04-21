package crawlers_test

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/dc_sale_index.html
var indexHtml string

//go:embed testdata/dc_sale_page.html
var pageHtml string

func TestPttCrawler(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, indexHtml)
	}))
	defer svr.Close()

	c := crawlers.PttCrawler{
		BaseUrl: svr.URL,
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("(?i)(sony|sigma)"),
			regexp.MustCompile("(?i)50(mm)?"),
			regexp.MustCompile("(?i)f2.8"),
		},
	}

	results := []*crawlers.Page{}
	c.Crawl(context.TODO(), &results)
	t.Logf("%+v", results)
	assert.Equal(t, 1, len(results))
}

func TestCrawlImgur(t *testing.T) {
	pattern, err := regexp.Compile("(?i)https://i.imgur.com/[0-9A-Z]+.(?:jpg|jpeg|png)")
	if err != nil {
		t.Fatal(err)
	}

	matched := pattern.FindAllString(pageHtml, -1)
	assert.Equal(t, 8, len(matched))
}
