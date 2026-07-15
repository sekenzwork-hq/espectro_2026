package pkg

import "strconv"

func GetOffset(limit int, page int) int {
	return (page - 1) * limit
}

func ParsePageAndLimit(limitStr string, pageStr string) (limit int, page int) {
	var pageV int
	var limitV int

	pageQ, pageIntErr := strconv.Atoi(pageStr)
	limitQ, limitIntErr := strconv.Atoi(limitStr)

	if pageIntErr != nil || pageQ < 0 {
		pageV = 1
	} else {
		pageV = pageQ
	}

	if limitIntErr != nil || limitQ < 0 {
		limitV = 50
	} else {
		limitV = limitQ
	}

	return limitV, pageV
}
