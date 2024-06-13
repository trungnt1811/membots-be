package dto

type MemeceptionCommon struct {
	StartAt         uint64 `json:"startAt,omitempty"`
	Status          uint64 `json:"status,omitempty"`
	Ama             bool   `json:"ama,omitempty"`
	ContractAddress string `json:"contractAddress,omitempty"`
	TargetETH       string `json:"targetETH,omitempty"`
	CollectedETH    string `json:"collectedETH,omitempty"`
	Enabled         bool   `json:"enabled,omitempty"`
	MemeID          uint64 `json:"memeID,omitempty"`
	UpdatedAtEpoch  uint64 `json:"updatedAtEpoch,omitempty"`
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
