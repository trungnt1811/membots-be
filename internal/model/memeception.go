package model

import (
	"time"

	"github.com/flexstack.ai/membots-be/internal/dto"
)

type Meme struct {
	ID         			uint    		`json:"id" gorm:"primaryKey"`
	Name     			string  		`json:"name"`
	Symbol      		string    		`json:"symbol"`
	Description       	string     		`json:"description"`
	TotalSupply        	string    		`json:"total_supply"`
	Decimals           	uint       		`json:"decimals"`
	LogoUrl            	string     		`json:"logo_url"`
	BannerUrl          	string     		`json:"banner_url"`
	CreatorAddress		string     		`json:"creator_address"`
	ContractAddress		string     		`json:"contract_address"`
	SwapFeeBps         	uint       		`json:"swap_fee_bps"`
	VestingAllocBps		uint       		`json:"vesting_alloc_bps"`
	Memerc20           	bool       		`json:"memerc20"`
	Live               	bool       		`json:"live"`
	NetworkID          	uint       		`json:"network_id"`
	Website            	string     		`json:"website"`
	Memeception        	Memeception		`json:"memeception" gorm:"foreignKey:MemeID;references:ID"`
	Social     			Social     		`json:"social" gorm:"foreignKey:MemeID;references:ID"`
	CreatedAt 			time.Time 		`json:"created_at"`
	UpdatedAt 			time.Time 		`json:"updated_at"`
}

func (m *Meme) TableName() string {
	return "meme"
}

func (m *Meme) ToDto() dto.Meme {
	return dto.Meme{
		ID: m.ID,
		Name: m.Name,
		Symbol: m.Symbol,
		Description: m.Description,
		TotalSupply: m.TotalSupply,
		Decimals: m.Decimals,
		LogoUrl: m.LogoUrl,
		BannerUrl: m.BannerUrl,
		CreatorAddress: m.CreatorAddress,
		ContractAddress: m.ContractAddress,
		SwapFeeBps: m.SwapFeeBps,
		VestingAllocBps: m.VestingAllocBps,
		Memerc20: m.Memerc20,
		Live: m.Live,
		NetworkID: m.NetworkID,
		Website: m.Website,
		Memeception: m.Memeception.ToDto(),
		Socials: m.Social.ToMapDto(),
	}
}

type Memeception struct {
	ID         			uint    	`json:"id" gorm:"primaryKey"`
	MemeID         		uint    	`json:"meme_id"`
	StartAt				uint		`json:"start_at"`
	Status             	uint    	`json:"status"`
	Ama                	bool   		`json:"ama"`
	ContractAddress		string     	`json:"contract_address"`
}

func (m *Memeception) TableName() string {
	return "memeception"
}

func (m *Memeception) ToDto() dto.Memeception {
	return dto.Memeception{
		StartAt: m.StartAt,
		Status: m.Status,
		Ama: m.Ama,
		ContractAddress: m.ContractAddress,
	}
}
