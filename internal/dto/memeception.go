package dto

type MemeceptionCommon struct {
	StartAt         uint   `json:"startAt"`
	Status          uint   `json:"status"`
	Ama             bool   `json:"ama"`
	ContractAddress string `json:"contractAddress"`
	TargetETH       string `json:"targetETH"`
	CollectedETH    string `json:"collectedETH"`
	Enabled         bool   `json:"enabled"`
	MemeID          uint   `json:"memeID"`
	UpdatedAtEpoch  uint   `json:"updatedAtEpoch"`
}

type Memeception struct {
	StartAt         uint       `json:"startAt"`
	Status          uint       `json:"status"`
	Ama             bool       `json:"ama"`
	ContractAddress string     `json:"contractAddress"`
	TargetETH       string     `json:"targetETH"`
	CollectedETH    string     `json:"collectedETH"`
	Enabled         bool       `json:"enabled"`
	MemeID          uint       `json:"memeID"`
	UpdatedAtEpoch  uint       `json:"updatedAtEpoch"`
	Meme            MemeCommon `json:"meme"`
}

type MemeceptionByStatus struct {
	Live []Memeception `json:"live"`
	Past []Memeception `json:"past"`
}
