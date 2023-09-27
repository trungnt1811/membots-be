package reward

import (
	"context"
	"database/sql"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"

	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/gorm"
)

type rewardRepository struct {
	db *gorm.DB
}

func NewRewardRepository(db *gorm.DB) interfaces.RewardRepository {
	return &rewardRepository{
		db: db,
	}
}

func (r *rewardRepository) CreateReward(ctx context.Context, reward *model.Reward) error {
	return r.db.Create(reward).Error
}

func (r *rewardRepository) GetInProgressRewards(ctx context.Context, userId uint32) ([]model.Reward, error) {
	var rewards []model.Reward
	query := "SELECT r.id, r.user_id, r.accesstrade_order_id, r.amount, r.rewarded_amount, r.commission_fee, " +
		"r.immediate_release, r.end_at, r.start_at, r.created_at, r.updated_at " +
		"FROM aff_reward AS r " +
		"LEFT JOIN aff_order as o ON o.accesstrade_order_id = r.accesstrade_order_id " +
		"WHERE o.user_id = ? AND o.order_status = ?"
	err := r.db.Raw(query, userId, model.OrderStatusRewarding).Scan(&rewards).Error
	return rewards, err
}

func (r *rewardRepository) GetRewardsInDay(ctx context.Context) ([]model.Reward, error) {
	startDay := time.Now().UTC().Round(24 * time.Hour)           // 00:00
	endDay := startDay.Add(24 * time.Hour).Add(-1 * time.Second) // 23:59

	var rewards []model.Reward
	err := r.db.Where("created_at BETWEEN ? AND ?", startDay, endDay).Find(&rewards).Error
	if err != nil {
		return []model.Reward{}, err
	}

	return rewards, err
}

func (r *rewardRepository) SaveRewardWithdraw(ctx context.Context, rewardWithdraw *model.RewardWithdraw, rewards []model.Reward, orderRewardHistories []model.OrderRewardHistory, completeRwOrders []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rewardWithdraw).Error; err != nil {
			return err
		}

		for idx := range orderRewardHistories {
			orderRewardHistories[idx].RewardWithdrawID = rewardWithdraw.ID
		}
		if err := tx.Create(orderRewardHistories).Error; err != nil {
			return err
		}

		if err := tx.Save(rewards).Error; err != nil {
			return err
		}

		if len(completeRwOrders) > 0 {
			if err := tx.Model(&model.AffOrder{}).Where("accesstrade_order_id IN ?", completeRwOrders).Update("order_status", model.OrderStatusComplete).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *rewardRepository) UpdateWithdrawShippingStatus(ctx context.Context, shippingReqId, txHash, status string) error {
	updates := map[string]interface{}{"tx_hash": txHash, "shipping_status": status}
	return r.db.Model(&model.RewardWithdraw{}).Where("shipping_request_id = ?", shippingReqId).Updates(updates).Error
}

func (r *rewardRepository) GetWithdrawById(ctx context.Context, userId uint32, withdrawId uint) (model.RewardWithdraw, error) {
	var withdraw model.RewardWithdraw
	err := r.db.Model(&model.RewardWithdraw{}).Where("id = ? AND user_id = ?", withdrawId, userId).First(&withdraw).Error
	return withdraw, err
}

func (r *rewardRepository) GetWithdrawHistory(ctx context.Context, userId uint32, page, size int) ([]model.RewardWithdraw, error) {
	var withdrawHistory []model.RewardWithdraw
	offset := (page - 1) * size
	err := r.db.Model(&model.RewardWithdraw{}).Where("user_id = ?", userId).Limit(size + 1).Offset(offset).Scan(&withdrawHistory).Error
	return withdrawHistory, err
}

func (r *rewardRepository) CountWithdrawal(ctx context.Context, userId uint32) (int64, error) {
	var count int64
	err := r.db.Model(&model.RewardWithdraw{}).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}

func (r *rewardRepository) GetTotalWithdrewAmount(ctx context.Context, userId uint32) (float64, error) {
	var amount sql.NullFloat64
	err := r.db.Model(&model.RewardWithdraw{}).Select("SUM(amount)").Where("user_id = ?", userId).Scan(&amount).Error
	return amount.Float64, err
}
