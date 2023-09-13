package goutil

func BinarySearch(itemCount int, checkFn func(index int) bool) int {
	var (
		left  = 0
		right = itemCount - 1
		index = -1
	)

	for left <= right {
		mid := left + (right-left)/2

		if checkFn(mid) {
			index = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return index
}
