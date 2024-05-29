package util

func IsValidOrder(order string) bool {
	switch order {
	case "asc", "desc":
		return true
	}
	return false
}
