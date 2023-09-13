package reward

import (
	"context"
	"time"

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

func (r *RewardRepository) GetRewardByOrderId(ctx context.Context, userId uint32, affOrderId uint) (model.Reward, error) {
	var reward model.Reward
	err := r.db.Model(&model.Reward{}).Where("user_id = ? AND aff_order_id = ?", userId, affOrderId).First(&reward).Error
	return reward, err
}

func (r *RewardRepository) GetRewardById(ctx context.Context, userId uint32, affOrderId uint) (model.Reward, error) {
	var reward model.Reward
	err := r.db.Model(&model.Reward{}).Where("user_id = ? AND id = ?", affOrderId).First(&reward).Error
	return reward, err
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

func (r *RewardRepository) GetAllReward(ctx context.Context, userId uint32, page, size int) ([]model.Reward, error) {
	var rewards []model.Reward
	offset := (page - 1) * size
	err := r.db.Model(&model.Reward{}).Where("user_id = ?", userId).Order("id DESC").Limit(size + 1).Offset(offset).Scan(&rewards).Error
	return rewards, err
}

func (r *RewardRepository) CountReward(ctx context.Context, userId uint32) (int64, error) {
	var count int64
	err := r.db.Model(&model.Reward{}).Where("user_id = ?", userId).Order("id DESC").Count(&count).Error
	return count, err
}

func (r *RewardRepository) SaveRewardClaim(ctx context.Context, rewardClaim *model.RewardClaim, rewards []model.Reward, orderRewardHistories []model.OrderRewardHistory) error {
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

func (r *RewardRepository) GetClaimHistory(ctx context.Context, userId uint32, page, size int) ([]model.RewardClaim, error) {
	var claimHistory []model.RewardClaim
	offset := (page - 1) * size
	err := r.db.Model(&model.RewardClaim{}).Where("user_id = ?", userId).Limit(size + 1).Offset(offset).Scan(&claimHistory).Error
	return claimHistory, err
}

func (r *RewardRepository) CountClaimHistory(ctx context.Context, userId uint32) (int64, error) {
	var count int64
	err := r.db.Model(&model.RewardClaim{}).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}
