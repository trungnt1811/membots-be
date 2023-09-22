package model

import (
	"math"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

const (
	FirstPartRewardPercent = 0.5
	OneDay                 = 24 * time.Hour
)

func (m *Reward) TableName() string {
	return "aff_reward"
}

type Reward struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	UserId         uint      `json:"user_id"`
	AtOrderID      string    `json:"accesstrade_order_id" gorm:"column:accesstrade_order_id"`
	Amount         float64   `json:"amount"` // amount of reward after fee subtractions
	RewardedAmount float64   `json:"rewarded_amount"`
	CommissionFee  float64   `json:"commission_fee"` // commission fee (in percentage)
	EndedAt        time.Time `json:"ended_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (r *Reward) ToRewardDto() dto.RewardDto {
	return dto.RewardDto{
		ID:             r.ID,
		UserId:         r.UserId,
		AtOrderID:      r.AtOrderID,
		Amount:         r.Amount,
		RewardedAmount: r.RewardedAmount,
		CommissionFee:  r.CommissionFee,
		EndedAt:        r.EndedAt,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

func (r *Reward) WithdrawableReward() (rewardAmount float64, ended bool) {
	daysPassed := int(time.Since(r.CreatedAt) / OneDay)   // number of days passed since order created
	totalDays := int(r.EndedAt.Sub(r.CreatedAt) / OneDay) // total lock days
	withdrawablePercent := float64(daysPassed) / float64(totalDays)
	if withdrawablePercent > 1 {
		withdrawablePercent = 1
		ended = true
	}

	rewardAmount = FirstPartRewardPercent*r.Amount + (1-FirstPartRewardPercent)*r.Amount*withdrawablePercent - r.RewardedAmount
	rewardAmount = math.Round(rewardAmount*100) / 100
	return
}
