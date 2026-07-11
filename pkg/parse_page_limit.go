package pkg

import (
	"strconv"
)

func ParsePageAndLimit(limitStr string, pageStr string) (limit int, page int) {
	var pageV int
	var limitV int

	pageQ, pageIntErr := strconv.Atoi(pageStr)
	limitQ, limitIntErr := strconv.Atoi(limitStr)

	if pageIntErr != nil {
		pageV = 1
	} else {
		pageV = pageQ
	}

	if limitIntErr != nil {
		limitV = 50
	} else {
		limitV = limitQ
	}

	return limitV, pageV
}
