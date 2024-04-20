package renthouse

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"artyomliou/sale-bot-v2/internal/utils"
	"fmt"
	"net/url"
)

type FiveNineOneAdapter struct {
	City       City
	Districts  []District
	Kind       Kind
	Room       []RoomCount
	PriceRange *Range
	FloorRange *Range
	Options    []Option
}

type Range struct {
	Min uint
	Max uint
}

func (a FiveNineOneAdapter) GetCrawler() crawlers.Crawler {
	queries := url.Values{
		"is_format_data":  []string{"1"},
		"is_new_list":     []string{"1"},
		"showMore":        []string{"1"},
		"recom_community": []string{"1"},
		"searchtype":      []string{"1"},
		"type":            []string{"1"},
	}

	switch a.City {
	case Taipei:
		queries.Add("region", "1")
	case NewTaipei:
		queries.Add("region", "3")
	}

	switch a.Kind {
	case Apartment:
		queries.Add("kind", "1")
	case Studio:
		queries.Add("kind", "2")
	}

	if a.PriceRange != nil {
		queries.Add("multiPrice", fmt.Sprintf("%d_%d", a.PriceRange.Min, a.PriceRange.Max))
	}
	if a.FloorRange != nil {
		queries.Add("floor", fmt.Sprintf("%d_%d", a.FloorRange.Min, a.FloorRange.Max))
	}

	qSection := utils.MultiValueBuilder{Sep: ","}
	for _, district := range a.Districts {
		switch district {
		case Banqiao:
			qSection.Add("26")
		case Xizhi:
			qSection.Add("27")
		case Xindian:
			qSection.Add("34")
		case Yonghe:
			qSection.Add("37")
		case Zhonghe:
			qSection.Add("38")
		case Sanxia:
			qSection.Add("40")
		case Shulin:
			qSection.Add("41")
		case Yingge:
			qSection.Add("42")
		case Sanchong:
			qSection.Add("43")
		case Xinzhuang:
			qSection.Add("44")
		case Taishan:
			qSection.Add("45")
		case Linkou:
			qSection.Add("46")
		case Luzhou:
			qSection.Add("47")
		case Wugu:
			qSection.Add("48")
		case Tamsui:
			qSection.Add("50")
		}
	}
	if qSection.Len() > 0 {
		queries.Add("section", qSection.Encode())
	}

	qMultiRoom := utils.MultiValueBuilder{Sep: ","}
	for _, r := range a.Room {
		switch r {
		case OneRoom:
			qMultiRoom.Add("1")
		case TwoRooms:
			qMultiRoom.Add("2")
		case ThreeRooms:
			qMultiRoom.Add("3")
		case FourRooms:
			qMultiRoom.Add("4")
		}
	}
	if qMultiRoom.Len() > 0 {
		queries.Add("multiRoom", qMultiRoom.Encode())
	}

	qOption := utils.MultiValueBuilder{Sep: ","}
	qMultiNotice := utils.MultiValueBuilder{Sep: ","}
	qOther := utils.MultiValueBuilder{Sep: ","}
	for _, opt := range a.Options {
		switch opt {
		case NoRoofTop:
			qMultiNotice.Add("not_cover")
		case AirConditioner:
			qOption.Add("cold")
		case WashingMachine:
			qOption.Add("washer")
		case Refrigerator:
			qOption.Add("icebox")
		case WaterHeater:
			qOption.Add("hotwater")
		case Internet:
			qOption.Add("broadband")
		case Bed:
			qOption.Add("bed")
		case Gas:
			qOption.Add("naturalgas")
		case Balcony:
			qOther.Add("balcony_1")
		case Elevator:
			qOther.Add("lift")
		case AllowPet:
			qOther.Add("pet")
		case AllowCookWithFire:
			qOther.Add("cook")
		}
	}
	if qOption.Len() > 0 {
		queries.Add("option", qOption.Encode())
	}
	if qMultiNotice.Len() > 0 {
		queries.Add("multiNotice", qMultiNotice.Encode())
	}
	if qOther.Len() > 0 {
		queries.Add("other", qOther.Encode())
	}

	crawler := crawlers.NewFiveNineOneCrawler()
	crawler.BaseUrl = "https://rent.591.com.tw"
	crawler.Queries = queries
	return crawler
}
