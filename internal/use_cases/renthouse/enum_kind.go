package renthouse

import (
	"fmt"
	"reflect"
)

type Kind int

const (
	Apartment Kind = iota
	Studio
)

var _KindStrings = map[Kind]string{
	Apartment: "家庭式",
	Studio:    "套房",
}

func (i Kind) String() string {
	if val, ok := _KindStrings[i]; ok {
		return val
	}
	typeName := reflect.TypeOf(i).Name()
	return fmt.Sprintf("%s(%s)", typeName, i)
}

func NewKindFromString(s string) (Kind, error) {
	for key, val := range _KindStrings {
		if val == s {
			return key, nil
		}
	}
	return 0, fmt.Errorf("%s is not defined as %s enum", s, reflect.TypeOf(City(0)).Name())
}
