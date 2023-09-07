package util

func IsValidOrder(order string) bool {
	switch order {
	case
		"asc",
		"desc":
		return true
	}
	return false
}

func NormalizeStatus(queryStatus string) []string {
	switch queryStatus {
	case "EXPIRED":
		return []string{"EXPIRED"}
	case "IN_PROGRESS":
		return []string{"IN_PROGRESS"}
	case "SCHEDULED":
		return []string{"SCHEDULED"}
	case "DRAFT":
		return []string{"DRAFT"}
	case "PAUSED":
		return []string{"PAUSED"}
	case "ENDED":
		return []string{"ENDED"}
	case "DELETED":
		return []string{"DELETED"}
	default:
		return []string{"EXPIRED", "IN_PROGRESS", "SCHEDULED", "DRAFT", "PAUSED", "ENDED", "DELETED"}
	}
}
