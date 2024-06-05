package model

import (
	"time"

	"github.com/flexstack.ai/membots-be/internal/dto"
)

type Meme struct {
	ID              uint64      `json:"id" gorm:"primaryKey"`
	Name            string      `json:"name"`
	Symbol          string      `json:"symbol"`
	Description     string      `json:"description"`
	TotalSupply     string      `json:"total_supply"`
	Decimals        uint64      `json:"decimals"`
	LogoUrl         string      `json:"logo_url"`
	BannerUrl       string      `json:"banner_url"`
	CreatorAddress  string      `json:"creator_address"`
	ContractAddress string      `json:"contract_address"`
	SwapFeeBps      uint64      `json:"swap_fee_bps"`
	VestingAllocBps uint64      `json:"vesting_alloc_bps"`
	Meta            string      `json:"meta"`
	Live            bool        `json:"live"`
	NetworkID       uint64      `json:"network_id"`
	Website         string      `json:"website"`
	Memeception     Memeception `json:"memeception" gorm:"foreignKey:MemeID;references:ID"`
	Social          Social      `json:"social" gorm:"foreignKey:MemeID;references:ID"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

func (m *Meme) TableName() string {
	return "meme"
}

func (m *Meme) ToDto() dto.MemeDetail {
	return dto.MemeDetail{
		ID:              m.ID,
		Name:            m.Name,
		Symbol:          m.Symbol,
		Description:     m.Description,
		TotalSupply:     m.TotalSupply,
		Decimals:        m.Decimals,
		LogoUrl:         m.LogoUrl,
		BannerUrl:       m.BannerUrl,
		CreatorAddress:  m.CreatorAddress,
		ContractAddress: m.ContractAddress,
		SwapFeeBps:      m.SwapFeeBps,
		VestingAllocBps: m.VestingAllocBps,
		Meta:            m.Meta,
		Live:            m.Live,
		NetworkID:       m.NetworkID,
		Website:         m.Website,
		Memeception:     m.Memeception.ToDto(),
		Socials:         m.Social.ToMapDto(),
	}
}

type MemeCommon struct {
	ID              uint64 `json:"id" gorm:"primaryKey"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Description     string `json:"description"`
	LogoUrl         string `json:"logo_url"`
	BannerUrl       string `json:"banner_url"`
	ContractAddress string `json:"contract_address"`
	Meta            string `json:"meta"`
}

func (m *MemeCommon) TableName() string {
	return "meme"
}

func (m *MemeCommon) ToCommonDto() dto.MemeCommon {
	return dto.MemeCommon{
		Name:            m.Name,
		Symbol:          m.Symbol,
		Description:     m.Description,
		LogoUrl:         m.LogoUrl,
		BannerUrl:       m.BannerUrl,
		ContractAddress: m.ContractAddress,
		Meta:            m.Meta,
	}
}
