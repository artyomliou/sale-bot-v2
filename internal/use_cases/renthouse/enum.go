package renthouse

//go:generate stringer -type=City -linecomment
type City int

const (
	Taipei    City = iota // 台北
	NewTaipei             // 新北
)

//go:generate stringer -type=District -linecomment
type District int

const (
	Banqiao   District = iota // 板橋
	Sanchong                  // 三重
	Tamsui                    // 淡水
	Zhonghe                   // 中和
	Yonghe                    // 永和
	Xizhi                     // 汐止
	Xindian                   // 新店
	Tucheng                   // 土城
	Sanxia                    // 三峽
	Shulin                    // 樹林
	Yingge                    // 鶯歌
	Xinzhuang                 // 新莊
	Taishan                   // 泰山
	Linkou                    // 林口
	Luzhou                    // 蘆洲
	Wugu                      // 五股
)

type Kind int

const (
	Apartment Kind = iota
	Studio
)

type RoomCount int

const (
	OneRoom RoomCount = iota
	TwoRooms
	ThreeRooms
	FourRooms
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
