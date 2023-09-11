package reward

import (
	"context"

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

func (r *RewardRepository) GetRewardByOrderId(ctx context.Context, affOrderId uint) (model.Reward, error) {
	var reward model.Reward
	err := r.db.Model(&model.Reward{}).Where("aff_order_id = ?", affOrderId).First(&reward).Error
	return reward, err
}

func (r *RewardRepository) GetAllReward(ctx context.Context, userId uint, page, size int) ([]model.Reward, error) {
	var rewards []model.Reward
	offset := (page - 1) * size
	err := r.db.Model(&model.Reward{}).Where("user_id = ?", userId).Order("id DESC").Limit(size + 1).Offset(offset).Scan(&rewards).Error
	return rewards, err
}

func (r *RewardRepository) CountReward(ctx context.Context, userId uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Reward{}).Where("user_id = ?", userId).Order("id DESC").Count(&count).Error
	return count, err
}

func (r *RewardRepository) GetRewardHistory(ctx context.Context, userId uint, page, size int) ([]model.RewardHistoryFull, error) {
	var rewardHistory []model.RewardHistoryFull
	offset := (page - 1) * size

	query := "SELECT rh.id, rh.reward_id, r.user_id, r.aff_order_id, rh.amount, rh.type, rh.created_at, rh.updated_at " +
		"FROM aff_reward_history AS rh " +
		"LEFT JOIN aff_reward AS r " +
		"ON rh.reward_id = r.id " +
		"WHERE r.user_id = ? " +
		"ORDER BY id DESC " +
		"LIMIT ? OFFSET ?"
	err := r.db.Raw(query, userId, size+1, offset).Scan(&rewardHistory).Error

	return rewardHistory, err
}

func (r *RewardRepository) CountRewardHistory(ctx context.Context, userId uint) (int64, error) {
	var count int64

	query := "SELECT count(*) " +
		"FROM aff_reward_history AS rh " +
		"LEFT JOIN aff_reward AS r " +
		"ON rh.reward_id = r.id " +
		"WHERE r.user_id = ? " +
		"ORDER BY id DESC"
	err := r.db.Raw(query, userId).Scan(&count).Error

	return count, err
}
