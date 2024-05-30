package dto

type Memeception struct {
	StartAt         uint   `json:"startAt"`
	Status          uint   `json:"status"`
	Ama             bool   `json:"ama"`
	ContractAddress string `json:"contractAddress"`
}

type MemeceptionCommon struct {
	StartAt         uint       `json:"startAt"`
	Status          uint       `json:"status"`
	Ama             bool       `json:"ama"`
	ContractAddress string     `json:"contractAddress"`
	Meme            MemeCommon `json:"meme"`
}
