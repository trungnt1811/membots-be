package model

import (
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

const (
	RewardStatusInProgress = "in_progress"
	RewardStatusDone       = "done"
	FirstPartRewardPercent = 0.5
	OneDay                 = 24 * time.Hour
)

func (m *Reward) TableName() string {
	return "aff_reward"
}

type Reward struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	UserId         uint      `json:"user_id"`
	AtOrderID      string    `json:"accesstrade_order_id"`
	Amount         float64   `json:"amount"`
	RewardedAmount float64   `json:"rewarded_amount"`
	Fee            float64   `json:"fee"`
	EndedAt        time.Time `json:"ended_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (r *Reward) ToRewardDto() dto.RewardDto {
	return dto.RewardDto{
		ID:        r.ID,
		UserId:    r.UserId,
		AtOrderID: r.AtOrderID,
		Amount:    r.Amount,
		EndedAt:   r.EndedAt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func (r *Reward) GetClaimableReward() (rewardAmount float64, ended bool) {
	daysPassed := int(time.Since(r.CreatedAt) / OneDay)   // number of days passed since order created
	totalDays := int(r.EndedAt.Sub(r.CreatedAt) / OneDay) // total lock days
	claimablePercent := float64(daysPassed) / float64(totalDays)
	if claimablePercent > 1 {
		claimablePercent = 1
		ended = true
	}

	rewardAmount = FirstPartRewardPercent*r.Amount + (1-FirstPartRewardPercent)*r.Amount*claimablePercent - r.RewardedAmount
	return
}
