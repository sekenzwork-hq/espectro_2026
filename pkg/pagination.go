package pkg

func GetOffset(limit int, page int) int {
	return (page - 1) * limit
}
