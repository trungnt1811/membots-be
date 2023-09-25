package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (repo *OrderRepository) FindNonRewardOrders(
	customerId, sellerId int,
	fromDate time.Time,
	minValue int64,
	additionalFilter map[string]interface{},
) ([]model.AffOrder, error) {
	var entities []model.AffOrder
	err := repo.db.
		Table("order").
		Where("customer_id = ?", customerId).
		Where("seller_id = ?", sellerId).
		Where("amount >= ?", minValue).
		Where("reward_id", nil).
		Where("initialized_at >= ?", fromDate).
		Find(&entities, additionalFilter).Error

	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (repo *OrderRepository) SavePostBackLog(req *model.AffPostBackLog) error {
	return repo.db.Create(req).Error
}

func (repo *OrderRepository) CreateOrder(order *model.AffOrder) error {
	return repo.db.Create(order).Error
}

func (repo *OrderRepository) UpdateOrder(updated *model.AffOrder) (int, error) {
	result := repo.db.Model(updated).Where("id = ?", updated.ID).Updates(updated)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func (repo *OrderRepository) UpdateOrCreateATTransactions(newTxs []model.AffTransaction) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		for _, newTx := range newTxs {
			// Find by id
			var oldTx model.AffTransaction
			err := tx.First(&oldTx, "accesstrade_conversion_id = ?", newTx.AccessTradeConversionId).Error
			if err != nil {
				if err.Error() == "record not found" {
					// If tx is not found, create one
					newTx.CreatedAt = time.Now()
					newTx.UpdatedAt = time.Now()
					crErr := tx.Create(&newTx).Error
					if crErr != nil {
						return crErr
					}
					continue
				} else {
					return err
				}
			}
			// When found one, compare and update
			newTx.ID = oldTx.ID
			newTx.UpdatedAt = time.Now()
			upErr := tx.Model(&oldTx).Updates(&newTx).Error
			if upErr != nil {
				return upErr
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (repo *OrderRepository) FindOrderByAccessTradeId(atOrderId string) (*model.AffOrder, error) {
	var order model.AffOrder
	err := repo.db.First(&order, "accesstrade_order_id = ?", atOrderId).Error
	return &order, err
}

func (repo *OrderRepository) UpdateTrackedClickOrder(trackedId uint64, order *model.AffOrder) error {
	// Only update empty order_id item
	err := repo.db.Model(&model.AffTrackedClick{}).Where(map[string]any{
		"id":       trackedId,
		"order_id": "",
	}).Updates(map[string]any{
		"order_id":   order.AccessTradeOrderId,
		"aff_link":   order.AffLink,
		"updated_at": time.Now(),
	}).Error

	return err
}

func (repo *OrderRepository) GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*model.OrderDetails, error) {
	var o model.OrderDetails
	query := "SELECT o.user_id, o.order_status, o.at_product_link, o.billing, o.category_name, o.merchant, " +
		"o.accesstrade_order_id, o.pub_commission, o.sales_time, o.confirmed_time, " +
		"lg.created_at, lg.data, " +
		"r.amount, r.rewarded_amount, r.commission_fee, r.ended_at, r.created_at " +
		"FROM aff_order AS o " +
		"LEFT JOIN aff_postback_log AS lg ON lg.order_id = o.accesstrade_order_id " +
		"LEFT JOIN aff_reward AS r ON r.accesstrade_order_id = o.accesstrade_order_id " +
		"WHERE o.user_id = ? AND o.id = ?"

	rows, err := repo.db.Raw(query, userId, orderId).Rows()
	if err != nil {
		return &model.OrderDetails{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var postbackDataJson sql.NullString
		var postbackCreatedAt sql.NullTime
		var rewardAmount sql.NullFloat64
		var rewardedAmount sql.NullFloat64
		var commissionFee sql.NullFloat64
		var rewardEndedAt sql.NullTime
		var rewardCreatedAt sql.NullTime
		err = rows.Scan(&o.UserId, &o.OrderStatus, &o.ATProductLink, &o.Billing, &o.CategoryName, &o.Merchant,
			&o.AccessTradeOrderId, &o.PubCommission, &o.SalesTime, &o.ConfirmedTime, &postbackCreatedAt, &postbackDataJson,
			&rewardAmount, &rewardedAmount, &commissionFee, &rewardEndedAt, &rewardCreatedAt)
		if err != nil {
			return &model.OrderDetails{}, err
		}
		o.RewardAmount = rewardAmount.Float64
		o.RewardedAmount = rewardedAmount.Float64
		o.CommissionFee = commissionFee.Float64
		o.EndedAt = rewardEndedAt.Time
		o.CreatedAt = rewardCreatedAt.Time

		var postbackData dto.ATPostBackRequest
		err = json.Unmarshal([]byte(postbackDataJson.String), &postbackData)
		if err != nil {
			return &model.OrderDetails{}, err
		}

		if postbackData.Status == dto.REQ_STATUS_APPROVED {
			o.ApprovedTime = postbackCreatedAt.Time
		} else if postbackData.Status == dto.REQ_STATUS_REJECTED {
			o.RejectedTime = postbackCreatedAt.Time
		}
	}

	if o.UserId == 0 {
		return &o, gorm.ErrRecordNotFound
	}

	return &o, err
}

func (repo *OrderRepository) GetOrderHistory(ctx context.Context, userId uint32, page, size int) ([]model.OrderDetails, error) {
	orderHistory := []model.OrderDetails{}
	// limit := size + 1
	// offset := (page - 1) * size
	// query := "SELECT o.user_id, o.order_status, o.at_product_link, o.billing, o.category_name, o.confirmed_time, o.merchant, " +
	// 	"o.accesstrade_order_id, o.pub_commission, o.sales_time, " +
	// 	"r.id, r.user_id, r.accesstrade_order_id, r.amount, r.rewarded_amount, " +
	// 	"r.commission_fee, r.ended_at, r.created_at, r.updated_at " +
	// 	"FROM aff_order AS o " +
	// 	"LEFT JOIN aff_reward AS r ON r.accesstrade_order_id = o.accesstrade_order_id " +
	// 	"WHERE o.user_id = ? " +
	// 	"ORDER BY o.id DESC " +
	// 	"LIMIT ? OFFSET ?"

	// rows, err := repo.db.Raw(query, userId, limit, offset).Rows()
	// if err != nil {
	// 	return []model.OrderDetails{}, err
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var o model.OrderDetails
	// 	err = rows.Scan(&o.UserId, &o.OrderStatus, &o.ATProductLink, &o.Billing, &o.CategoryName, &o.ConfirmedTime, &o.Merchant,
	// 		&o.AccessTradeOrderId, &o.PubCommission, &o.SalesTime,
	// 		&o.Reward.ID, &o.Reward.UserId, &o.Reward.AtOrderID, &o.Reward.Amount, &o.Reward.RewardedAmount,
	// 		&o.Reward.CommissionFee, &o.Reward.EndedAt, &o.Reward.CreatedAt, &o.Reward.UpdatedAt)
	// 	if err != nil {
	// 		return []model.OrderDetails{}, err
	// 	}

	// 	orderHistory = append(orderHistory, o)
	// }

	return orderHistory, nil
}

func (repo *OrderRepository) CountOrder(ctx context.Context, userId uint32) (int64, error) {
	var count int64
	query := "SELECT o.id " +
		"FROM aff_order AS o " +
		"LEFT JOIN aff_reward AS r ON r.accesstrade_order_id = o.accesstrade_order_id " +
		"WHERE o.user_id = ? "
	err := repo.db.Raw(query, userId).Count(&count).Error
	return count, err
}
