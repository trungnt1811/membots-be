package dto

type HomePageDto struct {
	RecentlyVisited []AffCampaignLessDto `json:"recently_visited"`
	MostCommission  []AffCampaignLessDto `json:"most_commission"`
	Following       []AffCampaignLessDto `json:"following"`
	TopFavorited    []AffCampaignLessDto `json:"top_favorited"`
}
