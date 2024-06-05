package model

import (
	"fmt"

	"github.com/flexstack.ai/membots-be/internal/dto"
)

type Memeception struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	StartAt         uint       `json:"start_at"`
	Status          uint       `json:"status"` // 2: SOLD OUT
	Ama             bool       `json:"ama"`
	ContractAddress string     `json:"contract_address"`
	TargetETH       float64    `json:"target_eth"`
	CollectedETH    float64    `json:"collected_eth"`
	Enabled         bool       `json:"enabled"`
	MemeID          uint       `json:"meme_id"`
	UpdatedAtEpoch  uint       `json:"updated_at_epoch"`
	Meme            MemeCommon `json:"meme" gorm:"foreignKey:ID;references:MemeID"`
}

func (m *Memeception) TableName() string {
	return "memeception"
}

func (m *Memeception) ToCommonDto() dto.MemeceptionCommon {
	return dto.MemeceptionCommon{
		StartAt:         m.StartAt,
		Status:          m.Status,
		Ama:             m.Ama,
		ContractAddress: m.ContractAddress,
		TargetETH:       fmt.Sprintf("%f", m.TargetETH),
		CollectedETH:    fmt.Sprintf("%f", m.CollectedETH),
		Enabled:         m.Enabled,
		MemeID:          m.MemeID,
		UpdatedAtEpoch:  m.UpdatedAtEpoch,
	}
}

func (m *Memeception) ToDto() dto.Memeception {
	return dto.Memeception{
		StartAt:         m.StartAt,
		Status:          m.Status,
		Ama:             m.Ama,
		ContractAddress: m.ContractAddress,
		TargetETH:       fmt.Sprintf("%f", m.TargetETH),
		CollectedETH:    fmt.Sprintf("%f", m.CollectedETH),
		Enabled:         m.Enabled,
		MemeID:          m.MemeID,
		Meme:            m.Meme.ToCommonDto(),
	}
}
