package handlers

import "strconv"

// ParsePagination 解析分页参数
func ParsePagination(offset, limit string) (int, int) {
	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		offsetInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 10
	} else if limitInt > 100 {
		limitInt = 100
	}

	return offsetInt, limitInt
}
