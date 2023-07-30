package util

import (
	"strconv"
	"time"
)

// YYYYMMDD
const layout = "20060102"

func GetYearDateRange(date string) (start, end uint64, err error) {
	t, err := ParseDateStr(date)
	if err != nil {
		return 0, 0, err
	}

	y := t.Year()

	fd := time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())
	ld := time.Date(y+1, time.January, 0, 0, 0, 0, 0, t.Location())

	start = FormatDateAsInt(fd)
	end = FormatDateAsInt(ld)

	return
}

func GetMonthDateRange(date string) (start, end uint64, err error) {
	t, err := ParseDateStr(date)
	if err != nil {
		return 0, 0, err
	}

	y, m, _ := t.Date()

	fd := time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
	ld := fd.AddDate(0, 1, -1)

	start = FormatDateAsInt(fd)
	end = FormatDateAsInt(ld)

	return
}

func FormatDate(t time.Time) string {
	return t.Format(layout)
}

func FormatDateAsInt(t time.Time) uint64 {
	d, _ := ParseDateStrToInt(FormatDate(t))
	return d
}

func ParseDateStr(s string) (time.Time, error) {
	return time.Parse(layout, s)
}

func ParseDateStrToInt(s string) (uint64, error) {
	_, err := ParseDateStr(s)
	if err != nil {
		return 0, nil
	}

	di, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, nil
	}

	return di, nil
}
