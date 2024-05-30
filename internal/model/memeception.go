package model

import (
	"github.com/flexstack.ai/membots-be/internal/dto"
)

type Memeception struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	MemeID          uint       `json:"meme_id"`
	StartAt         uint       `json:"start_at"`
	Status          uint       `json:"status"` // 2: SOLD OUT
	Ama             bool       `json:"ama"`
	ContractAddress string     `json:"contract_address"`
	Meme            MemeCommon `json:"meme" gorm:"foreignKey:ID;references:MemeID"`
}

func (m *Memeception) TableName() string {
	return "memeception"
}

func (m *Memeception) ToDto() dto.Memeception {
	return dto.Memeception{
		StartAt:         m.StartAt,
		Status:          m.Status,
		Ama:             m.Ama,
		ContractAddress: m.ContractAddress,
	}
}

func (m *Memeception) ToCommonRespDto() dto.MemeceptionCommon {
	return dto.MemeceptionCommon{
		StartAt:         m.StartAt,
		Status:          m.Status,
		Ama:             m.Ama,
		ContractAddress: m.ContractAddress,
		Meme:            m.Meme.ToDto(),
	}
}
