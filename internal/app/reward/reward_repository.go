package reward

import (
	"context"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type RewardRepository struct {
	db *gorm.DB
}

func NewRewardRepository(db *gorm.DB) *RewardRepository {
	return &RewardRepository{
		db: db,
	}
}

func (r *RewardRepository) CreateReward(ctx context.Context, reward *model.Reward) error {
	return r.db.Create(reward).Error
}

func (r *RewardRepository) GetInProgressRewards(ctx context.Context, userId uint32) ([]model.Reward, error) {
	var rewards []model.Reward
	err := r.db.Model(&model.Reward{}).Where("user_id = ? AND status = ?", userId, model.RewardStatusInProgress).Order("id DESC").Scan(&rewards).Error
	return rewards, err
}

func (r *RewardRepository) GetRewardsInDay(ctx context.Context) ([]model.Reward, error) {
	startDay := time.Now().UTC().Round(24 * time.Hour)            // 00:00
	nextDay := startDay.Add(24 * time.Hour).Add(-1 * time.Second) // 23:59

	var rewards []model.Reward
	err := r.db.Where("created_at BETWEEN ? AND ?", startDay, nextDay).Scan(&rewards).Error
	if err != nil {
		return []model.Reward{}, err
	}

	return rewards, err
}

func (r *RewardRepository) SaveRewardWithdraw(ctx context.Context, rewardClaim *model.RewardWithdraw, rewards []model.Reward, orderRewardHistories []model.OrderRewardHistory) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rewardClaim).Error; err != nil {
			return err
		}

		if err := tx.Create(orderRewardHistories).Error; err != nil {
			return err
		}

		if err := tx.Save(rewards).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *RewardRepository) GetWithdrawHistory(ctx context.Context, userId uint32, page, size int) ([]model.RewardWithdraw, error) {
	var withdrawHistory []model.RewardWithdraw
	offset := (page - 1) * size
	err := r.db.Model(&model.RewardWithdraw{}).Where("user_id = ?", userId).Limit(size + 1).Offset(offset).Scan(&withdrawHistory).Error
	return withdrawHistory, err
}

func (r *RewardRepository) CountWithdrawal(ctx context.Context, userId uint32) (int64, error) {
	var count int64
	err := r.db.Model(&model.RewardWithdraw{}).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}

type RewardWithdrawDetailsDto struct {
	ID                uint      `gorm:"primarykey" json:"id"`
	UserId            uint      `json:"user_id"`
	ShippingRequestID string    `json:"shipping_request_id"`
	Amount            float64   `json:"amount"`
	Fee               float64   `json:"fee"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (r *RewardRepository) GetWithdrawDetails(ctx context.Context, userId uint32, withdrawId uint) (dto.RewardWithdrawDetailsDto, error) {
	var withdrawDetails dto.RewardWithdrawDetailsDto
	query := "SELECT c.id, c.user_id, c.shipping_request_id, c.amount, c.fee, d.tx_hash, c.created_at, c.updated_at " +
		"FROM aff_reward_withdraw AS c " +
		"LEFT JOIN reward AS r ON c.shipping_request_id = r.request_id " +
		"LEFT JOIN deliveries AS d ON r.delivery_id = d.id " +
		"WHERE c.user_id = ? AND c.id = ?"
	err := r.db.Raw(query, userId, withdrawId).First(&withdrawDetails).Error
	return withdrawDetails, err
}

func (r *RewardRepository) GetTotalWithdrewAmount(ctx context.Context, userId uint32) (float64, error) {
	var amount float64
	err := r.db.Model(&model.RewardWithdraw{}).Select("COUNT(amount)").Where("user_id = ?", userId).Scan(&amount).Error
	return amount, err
}
