package dto

type MemeceptionCommon struct {
	StartAt         uint64 `json:"startAt"`
	Status          uint64 `json:"status"`
	Ama             bool   `json:"ama"`
	ContractAddress string `json:"contractAddress"`
	TargetETH       string `json:"targetETH"`
	CollectedETH    string `json:"collectedETH"`
	Enabled         bool   `json:"enabled"`
	MemeID          uint64 `json:"memeID"`
	UpdatedAtEpoch  uint64 `json:"updatedAtEpoch"`
}

type Memeception struct {
	StartAt         uint64     `json:"startAt"`
	Status          uint64     `json:"status"`
	Ama             bool       `json:"ama"`
	ContractAddress string     `json:"contractAddress"`
	TargetETH       string     `json:"targetETH"`
	CollectedETH    string     `json:"collectedETH"`
	Enabled         bool       `json:"enabled"`
	MemeID          uint64     `json:"memeID"`
	UpdatedAtEpoch  uint64     `json:"updatedAtEpoch"`
	Meme            MemeCommon `json:"meme"`
}

type MemeceptionByStatus struct {
	Live []Memeception `json:"live"`
	Past []Memeception `json:"past"`
}
