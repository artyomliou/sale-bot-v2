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
var html string

func TestPttCrawler(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, html)
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
