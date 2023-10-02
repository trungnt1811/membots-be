package dto

type Cashback struct {
	Distributed float64 `json:"distributed"`
	Remain      float64 `json:"remain"`
}

type StatisticSummaryResponse struct {
	NumOfOrders      int      `json:"num_of_orders"`
	NumOfCustomers   int      `json:"num_of_customers"`
	ActiveCampaigns  int      `json:"active_campaigns"`
	TotalRevenue     float64  `json:"total_revenue"`
	TotalASACashback Cashback `json:"total_asa_cashback"`
}

type CampaignSummaryResponse struct {
	NumOfOrders      int      `json:"num_of_orders"`
	NumOfCustomers   int      `json:"num_of_customers"`
	TotalRevenue     float64  `json:"total_revenue"`
	TotalASACashback Cashback `json:"total_asa_cashback"`
}
