package configs

import (
	"artyomliou/sale-bot-v2/internal/crawlers"
	"artyomliou/sale-bot-v2/internal/use_cases/camera"
	"artyomliou/sale-bot-v2/internal/use_cases/renthouse"
	"encoding/json"
	"log"
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Telegram struct {
		BotKey    string `json:"bot_key"`
		ChannelId int    `json:"channel_id"`
	}
	Db struct {
		SqliteFile string `json:"sqlite_file"`
	}
	Targets []map[string]interface{}
}

func GetConfig() *Config {
	var cfg Config

	cfgCandidates := []string{
		"config.json",
		"configs/config.json",
	}
	cfgCandidateSelected := false
	for _, cfgFilepath := range cfgCandidates {
		b, err := os.ReadFile(cfgFilepath)
		if err != nil {
			log.Print(err)
			continue
		}
		if err := json.Unmarshal(b, &cfg); err != nil {
			log.Print(err)
			continue
		}
		cfgCandidateSelected = true
	}
	if !cfgCandidateSelected {
		log.Fatalln("cannot load any valid config")
	}

	return &cfg
}

func (cfg *Config) GetAdapters() ([]crawlers.CrawlerAdapter, error) {
	adapters := []crawlers.CrawlerAdapter{}
	for _, target := range cfg.Targets {
		switch target["type"].(string) {
		case "ptt_dc_sale":
			var adapter camera.PttCrawlerDcSaleAdapter
			if err := decodeAdapter(target, &adapter); err != nil {
				return nil, err
			}
			adapters = append(adapters, adapter)
		case "ptt_rent_apart":
			var adapter renthouse.PttCrawlerRentApartAdapter
			if err := decodeAdapter(target, &adapter); err != nil {
				return nil, err
			}
			adapters = append(adapters, adapter)
		case "five_nine_one":
			var adapter renthouse.FiveNineOneAdapter
			if err := decodeAdapter(target, &adapter); err != nil {
				return nil, err
			}
			adapters = append(adapters, adapter)
		}
	}
	return adapters, nil
}

func decodeAdapter(input interface{}, output interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		DecodeHook: func(
			f reflect.Type,
			t reflect.Type,
			data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t == reflect.TypeOf(renthouse.City(0)) {
				return renthouse.NewCityFromString(data.(string))
			}
			if t == reflect.TypeOf(renthouse.District(0)) {
				return renthouse.NewDistrictFromString(data.(string))
			}
			if t == reflect.TypeOf(renthouse.Kind(0)) {
				return renthouse.NewKindFromString(data.(string))
			}
			if t == reflect.TypeOf(renthouse.Option(0)) {
				return renthouse.NewOptionFromString(data.(string))
			}
			if t == reflect.TypeOf(renthouse.RoomCount(0)) {
				return renthouse.NewRoomCountFromString(data.(string))
			}

			return data, nil
		},
	})
	if err != nil {
		return err
	}
	err = decoder.Decode(input)
	if err != nil {
		return err
	}
	return nil
}
