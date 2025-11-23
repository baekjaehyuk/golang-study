package model

type Day string

const (
	Monday    Day = "MON"
	Tuesday   Day = "TUE"
	Wednesday Day = "WED"
	Thursday  Day = "THU"
	Friday    Day = "FRI"
)

func (d Day) ToKorean() string {
	switch d {
	case Monday:
		return "월요일"
	case Tuesday:
		return "화요일"
	case Wednesday:
		return "수요일"
	case Thursday:
		return "목요일"
	case Friday:
		return "금요일"
	default:
		return "undefined"
	}
}
