package util

import "time"

type MonthType uint32

const (
	Constant_JAN MonthType = iota + 1
	Constant_FEB
	Constant_MAR
	Constant_APR
	Constant_MAY
	Constant_JUN
	Constant_JUL
	Constant_AUG
	Constant_SEP
	Constant_OCT
	Constant_NOV
	Constant_DEC
)

var MonthTypes = map[uint32]string{
	uint32(Constant_JAN): "JAN",
	uint32(Constant_FEB): "FEB",
	uint32(Constant_MAR): "MAR",
	uint32(Constant_APR): "APR",
	uint32(Constant_MAY): "MAY",
	uint32(Constant_JUN): "JUN",
	uint32(Constant_JUL): "JUL",
	uint32(Constant_AUG): "AUG",
	uint32(Constant_SEP): "SEP",
	uint32(Constant_OCT): "OCT",
	uint32(Constant_NOV): "NOV",
	uint32(Constant_DEC): "DEC",
}

const layout = "2006-01-02"

func IsYearMonthAfterCurrent(year, month uint32) bool {
	currYear, currMonth, _ := time.Now().Date()
	if year > uint32(currYear) {
		return true
	} else if year == uint32(currYear) && month > uint32(currMonth) {
		return true
	}
	return false
}

// YYYY-MM-DD
func ValidateDateStr(dateString string) error {
	_, err := time.Parse(layout, dateString)
	if err != nil {
		return err
	}

	return nil
}

// YYYY-MM-DD
func DateStrToDate(dateString string) (time.Time, error) {
	date, err := time.Parse(layout, dateString)
	if err != nil {
		return date, err
	}

	return date, nil
}
