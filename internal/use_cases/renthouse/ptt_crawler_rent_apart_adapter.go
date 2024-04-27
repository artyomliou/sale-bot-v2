package renthouse

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"fmt"
	"regexp"
)

type PttCrawlerRentApartAdapter struct {
	Cities    []City     `mapstructure:"cities"`
	Districts []District `mapstructure:"districts"`
	Room      RoomCount  `mapstructure:"room"`
}

func (a PttCrawlerRentApartAdapter) GetCrawler() crawlers.Crawler {
	crawler := crawlers.NewPttCrawler()
	crawler.BaseUrl = crawlers.PttBaseUrl
	crawler.Board = "rent_apart"
	crawler.Patterns = []*regexp.Regexp{}

	if len(a.Districts) > 0 {
		for _, district := range a.Districts {
			crawler.Patterns = append(crawler.Patterns, regexp.MustCompile(fmt.Sprintf("(?i)%s", district)))
		}
	} else if len(a.Cities) > 0 {
		for _, city := range a.Cities {
			crawler.Patterns = append(crawler.Patterns, regexp.MustCompile(fmt.Sprintf("(?i)%s", city)))
		}
	}

	switch a.Room {
	case OneRoom:
		crawler.Patterns = append(crawler.Patterns, regexp.MustCompile("(?i)(1|一)房"))
	case TwoRooms:
		crawler.Patterns = append(crawler.Patterns, regexp.MustCompile("(?i)(2|二)房"))
	case ThreeRooms:
		crawler.Patterns = append(crawler.Patterns, regexp.MustCompile("(?i)(3|三)房"))
	case FourRooms:
		crawler.Patterns = append(crawler.Patterns, regexp.MustCompile("(?i)(4|四)房"))
	}
	return crawler
}
