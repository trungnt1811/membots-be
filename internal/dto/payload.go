package dto

type MemeInfoPayload struct {
	SwapFeePct      float64 `json:"swapFeePct"`
	TargetETH       float64 `json:"targetETH"`
	VestingAllocBps string  `json:"vestingAllocBps"`
	Ama             bool    `json:"ama"`
	Name            string  `json:"name"`
	Symbol          string  `json:"symbol"`
	Description     string  `json:"description"`
	Website         string  `json:"website"`
	Telegram        string  `json:"telegram"`
	Salt            string  `json:"salt"`
	LogoUrl         string  `json:"logoUrl"`
	BannerUrl       string  `json:"bannerUrl"`
	SwapFeeBps      string  `json:"swapFeeBps"`
	Meta            string  `json:"meta"`
}

type MemeceptionPayload struct {
	StartAt   string `json:"startAt"`
	Ema       bool   `json:"ema"`
	TargetETH string `json:"targetETH"`
}

type BlockchainPayload struct {
	ChainID        string `json:"chainId"`
	CreatorAddress string `json:"creatorAddress"`
}

type CreateMemePayload struct {
	MemeInfo    MemeInfoPayload    `json:"memeInfo"`
	Memeception MemeceptionPayload `json:"memeception"`
	Blockchain  BlockchainPayload  `json:"blockchain"`
	Socials     map[string]Social  `json:"socials"`
}
