package pkg

import (
	"strconv"
)

func GetOffset(limit int, page int) int {
	return (page - 1) * limit
}

func ParsePageSetMin1(pageStr string) int {

	page, _ := strconv.Atoi(pageStr)

	if page <= 0 {
		page = 1
	}
	return page
}
