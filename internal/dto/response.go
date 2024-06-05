package dto

type MemeceptionDetailResp struct {
	Meme  MemeDetail `json:"meme"`
	Price uint64     `json:"price"`
}

type MemeceptionsResp struct {
	Live     []MemeceptionCommon `json:"live"`
	Upcoming []MemeceptionCommon `json:"upcoming"`
	Past     []MemeceptionCommon `json:"past"`
}

type SwapHistoryByAddressResp struct {
	Swaps []Swap `json:"swaps"`
}

type LaunchpadInfoResp struct {
	LaunchpadInfo LaunchpadInfo `json:"launchpadInfo"`
}
