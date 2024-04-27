package renthouse

import (
	"fmt"
	"reflect"
)

type RoomCount int

const (
	OneRoom RoomCount = iota
	TwoRooms
	ThreeRooms
	FourRooms
)

var _RoomCountStrings = map[RoomCount]string{
	OneRoom:    "一房",
	TwoRooms:   "兩房",
	ThreeRooms: "三房",
	FourRooms:  "四房",
}

func (i RoomCount) String() string {
	if val, ok := _RoomCountStrings[i]; ok {
		return val
	}
	typeName := reflect.TypeOf(i).Name()
	return fmt.Sprintf("%s(%s)", typeName, i)
}

func NewRoomCountFromString(s string) (RoomCount, error) {
	for key, val := range _RoomCountStrings {
		if val == s {
			return key, nil
		}
	}
	return 0, fmt.Errorf("%s is not defined as %s enum", s, reflect.TypeOf(City(0)).Name())
}
