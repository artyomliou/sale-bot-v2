package renthouse

import (
	"fmt"
	"reflect"
)

type City int

const (
	Taipei City = iota
	NewTaipei
)

var _CityStrings = map[City]string{
	Taipei:    "台北",
	NewTaipei: "新北",
}

func (i City) String() string {
	if val, ok := _CityStrings[i]; ok {
		return val
	}
	typeName := reflect.TypeOf(i).Name()
	return fmt.Sprintf("%s(%s)", typeName, i)
}

func NewCityFromString(s string) (City, error) {
	for key, val := range _CityStrings {
		if val == s {
			return key, nil
		}
	}
	return 0, fmt.Errorf("%s is not defined as %s enum", s, reflect.TypeOf(City(0)).Name())
}
