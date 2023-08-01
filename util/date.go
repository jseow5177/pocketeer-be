package util

import (
	"strconv"
	"time"
)

// YYYYMMDD
const layout = "20060102"

func GetYearRangeAsDate(date, timezone string) (start, end uint64, err error) {
	fd, ld, err := getYearRange(date, timezone)
	if err != nil {
		return 0, 0, nil
	}

	start = FormatDateAsInt(fd)
	end = FormatDateAsInt(ld)

	return
}

func GetYearRangeAsUnix(date, timezone string) (start, end uint64, err error) {
	fd, ld, err := getYearRange(date, timezone)
	if err != nil {
		return 0, 0, nil
	}

	start = uint64(fd.UnixMilli())
	end = uint64(ld.UnixMilli())

	return
}

func getYearRange(date, timezone string) (start, end time.Time, err error) {
	zeroTime := time.Time{}

	t, err := ParseDateStr(date)
	if err != nil {
		return zeroTime, zeroTime, err
	}

	l := t.Location()
	if timezone != "" {
		l, err = time.LoadLocation(timezone)
		if err != nil {
			return zeroTime, zeroTime, err
		}
	}

	y := t.Year()

	start = time.Date(y, time.January, 1, 0, 0, 0, 0, l)
	end = time.Date(y+1, time.January, 0, 0, 0, 0, 0, l)

	return
}

func GetMonthRangeAsDate(date, timezone string) (start, end uint64, err error) {
	fd, ld, err := getMonthRange(date, timezone)
	if err != nil {
		return 0, 0, nil
	}

	start = FormatDateAsInt(fd)
	end = FormatDateAsInt(ld)

	return
}

func GetMonthRangeAsUnix(date, timezone string) (start, end uint64, err error) {
	fd, ld, err := getMonthRange(date, timezone)
	if err != nil {
		return 0, 0, nil
	}

	start = uint64(fd.UnixMilli())
	end = uint64(ld.UnixMilli())

	return
}

func getMonthRange(date, timezone string) (start, end time.Time, err error) {
	zeroTime := time.Time{}

	t, err := ParseDateStr(date)
	if err != nil {
		return zeroTime, zeroTime, err
	}

	l := t.Location()
	if timezone != "" {
		l, err = time.LoadLocation(timezone)
		if err != nil {
			return zeroTime, zeroTime, err
		}
	}

	y, m, _ := t.Date()

	start = time.Date(y, m, 1, 0, 0, 0, 0, l)
	end = start.AddDate(0, 1, -1)

	return
}

func FormatDateAsStr(t time.Time) string {
	return t.Format(layout)
}

func FormatDateAsInt(t time.Time) uint64 {
	d, _ := ParseDateStrToInt(FormatDateAsStr(t))
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
