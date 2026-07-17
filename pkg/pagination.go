package pkg

import (
	customerrors "espectro/custom_errors"
	"strconv"
)

func GetOffset(limit int, page int) int {
	return (page - 1) * limit
}

func ParsePageAndLimit(limitStr string, pageStr string) (limit int, page int, limitErr error) {
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
	} else if limit > 150 {
		return 0, 0, &customerrors.SizeError{OrgError: "Limit should be less than or equal to 150"}

	} else {
		limitV = limitQ
	}

	return limitV, pageV, nil
}
