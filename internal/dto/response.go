package dto

type MemeceptionDetailResp struct {
	Meme  MemeDetail `json:"meme"`
	Price uint64     `json:"price"`
}

type MemeceptionsResp struct {
	MemeceptionsByStatus MemeceptionByStatus `json:"memeceptionsByStatus"`
	Price                uint64              `json:"price"`
	LatestLaunchpadTx    []LaunchpadTx       `json:"latestLaunchpadTx"`
	LatestCoins          []Memeception       `json:"latestCoins"`
}

type SwapHistoryByAddressResp struct {
	Swaps []Swap `json:"swaps"`
}

type LaunchpadInfoResp struct {
	LaunchpadInfo LaunchpadInfo `json:"launchpadInfo"`
}
