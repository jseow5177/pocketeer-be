package entity

import (
	"fmt"
	"sort"
)

func sortBudgetBreakdowns(slice []*BudgetBreakdown) {
	// Define a custom sorting function that sorts by year and then by month
	sort.Slice(slice, func(i, j int) bool {
		if slice[i].Year == slice[j].Year {
			return slice[i].Month < slice[j].Month
		}
		return slice[i].Year < slice[j].Year
	})
}

func getDateToBreakdownMap(breakdowns []*BudgetBreakdown) map[string]*BudgetBreakdown {
	_map := make(map[string]*BudgetBreakdown)

	for _, bd := range breakdowns {
		key := getYearMonthKey(bd.Year, bd.Month)
		_map[key] = bd
	}

	return _map
}

func getYearMonthKey(year, month uint32) string {
	return fmt.Sprintf("%d-%d", year, month)
}

func isDate1LargerEqual(year1, month1, year2, month2 uint32) bool {
	date1 := year1*100 + month1
	date2 := year2*100 + month2
	return date1 >= date2
}
