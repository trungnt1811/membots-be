package dto

type Cashback struct {
	Distributed float32 `json:"distributed"`
	Remain      float32 `json:"remain"`
}

type StatisticSummaryResponse struct {
	NumOfOrders      int      `json:"num_of_orders"`
	NumOfCustomers   int      `json:"num_of_customers"`
	TotalRevenue     float32  `json:"total_revenue"`
	TotalASACashback Cashback `json:"total_asa_cashback"`
}
