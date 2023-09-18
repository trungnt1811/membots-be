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
	case "IN_PROGRESS":
		return []string{"IN_PROGRESS"}
	case "DRAFT":
		return []string{"DRAFT"}
	case "PAUSED":
		return []string{"PAUSED"}
	case "ENDED":
		return []string{"ENDED"}
	case "DELETED":
		return []string{"DELETED"}
	default:
		return []string{"IN_PROGRESS", "DRAFT", "PAUSED", "ENDED", "DELETED"}
	}
}

func NormalizeStatusActiveInActive(queryStatus string) []string {
	switch queryStatus {
	case "active":
		return []string{"active"}
	case "inactive":
		return []string{"inactive"}
	default:
		return []string{"active", "inactive"}
	}
}

func BuildOrderBy(queryBy, order string) string {
	switch queryBy {
	case "ctime":
		return "aff_campaign.created_at " + order
	case "top":
		return "aff_campaign.id " + order
	}
	return ""
}
