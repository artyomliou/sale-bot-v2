package renthouse

import (
	"fmt"
	"reflect"
)

type Option int

const (
	NoRoofTop Option = iota
	AirConditioner
	WashingMachine
	Refrigerator
	WaterHeater
	Internet
	Bed
	Gas
	Balcony
	Elevator
	AllowPet
	AllowCookWithFire
)

var _OptionStrings = map[Option]string{
	NoRoofTop:         "非頂樓加蓋",
	AirConditioner:    "冷氣",
	WashingMachine:    "洗衣機",
	Refrigerator:      "冰箱",
	WaterHeater:       "熱水器",
	Internet:          "網路",
	Bed:               "床",
	Gas:               "天然氣",
	Balcony:           "陽台",
	Elevator:          "電梯",
	AllowPet:          "允許寵物",
	AllowCookWithFire: "允許明火",
}

func (i Option) String() string {
	if val, ok := _OptionStrings[i]; ok {
		return val
	}
	typeName := reflect.TypeOf(i).Name()
	return fmt.Sprintf("%s(%s)", typeName, i)
}

func NewOptionFromString(s string) (Option, error) {
	for key, val := range _OptionStrings {
		if val == s {
			return key, nil
		}
	}
	return 0, fmt.Errorf("%s is not defined as %s enum", s, reflect.TypeOf(City(0)).Name())
}
