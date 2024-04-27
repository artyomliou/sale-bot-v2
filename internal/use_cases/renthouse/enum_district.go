package renthouse

import (
	"fmt"
	"reflect"
)

type District int

const (
	Banqiao District = iota
	Sanchong
	Tamsui
	Zhonghe
	Yonghe
	Xizhi
	Xindian
	Tucheng
	Sanxia
	Shulin
	Yingge
	Xinzhuang
	Taishan
	Linkou
	Luzhou
	Wugu
)

var _DistrictStrings = map[District]string{
	Banqiao:   "板橋",
	Sanchong:  "三重",
	Tamsui:    "淡水",
	Zhonghe:   "中和",
	Yonghe:    "永和",
	Xizhi:     "汐止",
	Xindian:   "新店",
	Tucheng:   "土城",
	Sanxia:    "三峽",
	Shulin:    "樹林",
	Yingge:    "鶯歌",
	Xinzhuang: "新莊",
	Taishan:   "泰山",
	Linkou:    "林口",
	Luzhou:    "蘆洲",
	Wugu:      "五股",
}

func (i District) String() string {
	if val, ok := _DistrictStrings[i]; ok {
		return val
	}
	typeName := reflect.TypeOf(i).Name()
	return fmt.Sprintf("%s(%s)", typeName, i)
}

func NewDistrictFromString(s string) (District, error) {
	for key, val := range _DistrictStrings {
		if val == s {
			return key, nil
		}
	}
	return 0, fmt.Errorf("%s is not defined as %s enum", s, reflect.TypeOf(City(0)).Name())
}
