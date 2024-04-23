package crawlers_test

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"artyomliou/sale-bot-v2/internal/use_cases/renthouse"
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/591_luzhou_apartment.json
var json string

func TestFiveNineOneCrawler(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, json)
	}))
	defer svr.Close()

	adapter := renthouse.FiveNineOneAdapter{
		City: renthouse.NewTaipei,
		Districts: []renthouse.District{
			renthouse.Luzhou,
		},
		Kind: renthouse.Apartment,
		Room: []renthouse.RoomCount{
			renthouse.TwoRooms,
			renthouse.ThreeRooms,
		},
		PriceRange: &renthouse.Range{
			Min: 20000,
			Max: 30000,
		},
		FloorRange: &renthouse.Range{
			Min: 1,
			Max: 10,
		},
		Options: []renthouse.Option{
			renthouse.WashingMachine,
			renthouse.AirConditioner,
			renthouse.Refrigerator,
			renthouse.WaterHeater,
			renthouse.Bed,

			renthouse.NoRoofTop,
			renthouse.AllowPet,
			renthouse.Balcony,
		},
	}

	c := adapter.GetCrawler()
	if crawler, ok := c.(*crawlers.FiveNineOneCrawler); ok {
		crawler.BaseUrl = svr.URL
		results := []*crawlers.Page{}
		c.Crawl(context.Background(), &results)

		assert.Equal(t, 6, len(results))
	} else {
		t.Fatal("crawler is not FiveNineOneCrawler")
	}
}
