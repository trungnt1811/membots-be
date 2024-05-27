package dto

import "time"

type Meme20MetaDto struct {
	ID                           uint      `json:"id"`
	CreatedAt                    time.Time `json:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at"`
	Ticker                       string    `json:"ticker"`
	Description                  string    `json:"description"`
	TwitterURL                   string    `json:"twitter_url"`
	TelegramURL                  string    `json:"telegram_url"`
	CreatorSwapFee               string    `json:"creator_swap_fee"`
	CreatorAllocationVestedy     string    `json:"creator_allocation_vestedy"`
	MemeType                     string    `json:"meme_type"`
	FairLauchType                string    `json:"fair_lauch_type"`
}
