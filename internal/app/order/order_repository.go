package order

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"

	"github.com/astraprotocol/affiliate-system/internal/model"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) interfaces.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (repo *orderRepository) QueryOrdersConfirmedBefore(t time.Time, q map[string]any) ([]model.AffOrder, error) {
	var orders []model.AffOrder
	sql := repo.db.Model(&orders)
	if !t.IsZero() {
		sql.Where("sales_time <= ?", t)
	}

	sql.Where("order_status IN ?", []string{
		model.OrderStatusInitial,
		model.OrderStatusPending,
		model.OrderStatusApproved,
		model.OrderStatusRewarding,
	})

	err := sql.Find(&orders, q).Error
	return orders, err
}

func (repo *orderRepository) CreatePostBackLog(req *model.AffPostBackLog) error {
	return repo.db.Create(req).Error
}

func (repo *orderRepository) UpdatePostBackLog(id uint, changes map[string]any) error {
	return repo.db.Model(&model.AffPostBackLog{}).Where("id = ?", id).Updates(changes).Error
}

func (repo *orderRepository) CreateOrder(order *model.AffOrder) error {
	return repo.db.Create(order).Error
}

func (repo *orderRepository) UpdateOrder(updated *model.AffOrder) (int, error) {
	result := repo.db.Model(updated).Where("id = ?", updated.ID).Updates(updated)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func (repo *orderRepository) UpdateOrCreateATTransactions(newTxs []model.AffTransaction) error {
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

func (repo *orderRepository) FindOrderByAccessTradeId(atOrderId string) (*model.AffOrder, error) {
	var order model.AffOrder
	err := repo.db.First(&order, "accesstrade_order_id = ?", atOrderId).Error
	return &order, err
}

func (repo *orderRepository) FindOrderMappedByAccessTradeIds(atOrderIds []string) (map[string]model.AffOrder, error) {
	var orders []model.AffOrder
	ret := map[string]model.AffOrder{}
	err := repo.db.Find(&orders, "accesstrade_order_id IN ?", atOrderIds).Error
	if err != nil {
		return nil, err
	}
	for idx := range orders {
		ret[orders[idx].AccessTradeOrderId] = orders[idx]
	}
	return ret, nil
}

func (repo *orderRepository) GetCampaignByTrackedClick(trackedId uint64) (*model.AffCampaign, error) {
	var tracked model.AffTrackedClick

	err := repo.db.Joins("Campaign").First(&tracked, "aff_tracked_click.id = ?", trackedId).Error
	if err != nil {
		return nil, err
	}
	if tracked.Campaign == nil {
		return nil, errors.New("no campaign mapped for click")
	}

	return tracked.Campaign, nil
}

func (repo *orderRepository) UpdateTrackedClickOrder(trackedId uint64, order *model.AffOrder) error {
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

var selectOrderDetails = "SELECT o.id, o.user_id, o.order_status, o.billing, o.category_name, o.merchant, " +
	"o.accesstrade_order_id, o.pub_commission, o.update_time, o.cancelled_time, o.sales_time, o.confirmed_time, o.created_at, " +
	"r.amount, r.rewarded_amount, r.commission_fee, r.immediate_release, r.end_at, r.start_at, " +
	"b.logo, b.name AS brand_name " +
	"FROM aff_order AS o " +
	"LEFT JOIN aff_reward AS r ON r.accesstrade_order_id = o.accesstrade_order_id " +
	"LEFT JOIN brand AS b ON b.id = o.brand_id "

func (repo *orderRepository) GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*model.OrderDetails, error) {
	var o model.OrderDetails
	query := selectOrderDetails + "WHERE o.user_id = ? AND o.id = ?"

	rows, err := repo.db.Raw(query, userId, orderId).Rows()
	if err != nil {
		return &model.OrderDetails{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var rewardAmount sql.NullFloat64
		var rewardedAmount sql.NullFloat64
		var commissionFee sql.NullFloat64
		var immediateRelease sql.NullFloat64
		var cancelledTime sql.NullTime
		var rewardEndAt sql.NullTime
		var rewardStartAt sql.NullTime
		var brandLogo sql.NullString
		var brandName sql.NullString
		err = rows.Scan(&o.ID, &o.UserId, &o.OrderStatus, &o.Billing, &o.CategoryName, &o.Merchant,
			&o.AccessTradeOrderId, &o.PubCommission, &o.UpdateTime, &cancelledTime, &o.SalesTime, &o.ConfirmedTime, &o.CreatedAt,
			&rewardAmount, &rewardedAmount, &commissionFee, &immediateRelease, &rewardEndAt, &rewardStartAt,
			&brandLogo, &brandName)
		if err != nil {
			return &model.OrderDetails{}, err
		}
		o.CancelledTime = cancelledTime.Time
		o.RewardAmount = rewardAmount.Float64
		o.RewardedAmount = rewardedAmount.Float64
		o.CommissionFee = commissionFee.Float64
		o.ImmediateRelease = immediateRelease.Float64
		o.RewardEndAt = rewardEndAt.Time
		o.RewardStartAt = rewardStartAt.Time
		o.Brand.Logo = brandLogo.String
		o.Brand.Name = brandName.String
	}

	return &o, err
}

func (repo *orderRepository) GetOrderHistory(ctx context.Context, since time.Time, userId uint32, status string, page, size int) ([]model.OrderDetails, error) {
	orderHistory := []model.OrderDetails{}

	limit := size + 1
	offset := (page - 1) * size

	statusQuery, statusParams := model.BuildOrderStatusQuery(status)

	query := selectOrderDetails +
		"WHERE o.user_id = ? AND o.created_at > ? "
	if status != "" {
		query += statusQuery + " "
	}
	query += "ORDER BY o.sales_time DESC " +
		"LIMIT ? OFFSET ?"

	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = repo.db.Raw(query, userId, since, statusParams, limit, offset).Rows()
	} else {
		rows, err = repo.db.Raw(query, userId, since, limit, offset).Rows()
	}
	if err != nil {
		return []model.OrderDetails{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var o model.OrderDetails
		var rewardAmount sql.NullFloat64
		var rewardedAmount sql.NullFloat64
		var commissionFee sql.NullFloat64
		var immediateRelease sql.NullFloat64
		var cancelledTime sql.NullTime
		var rewardEndAt sql.NullTime
		var rewardStartAt sql.NullTime
		var brandLogo sql.NullString
		var brandName sql.NullString
		err = rows.Scan(&o.ID, &o.UserId, &o.OrderStatus, &o.Billing, &o.CategoryName, &o.Merchant,
			&o.AccessTradeOrderId, &o.PubCommission, &o.UpdateTime, &cancelledTime, &o.SalesTime, &o.ConfirmedTime, &o.CreatedAt,
			&rewardAmount, &rewardedAmount, &commissionFee, &immediateRelease, &rewardEndAt, &rewardStartAt,
			&brandLogo, &brandName)
		if err != nil {
			return []model.OrderDetails{}, err
		}
		o.CancelledTime = cancelledTime.Time
		o.RewardAmount = rewardAmount.Float64
		o.RewardedAmount = rewardedAmount.Float64
		o.CommissionFee = commissionFee.Float64
		o.ImmediateRelease = immediateRelease.Float64
		o.RewardEndAt = rewardEndAt.Time
		o.RewardStartAt = rewardStartAt.Time
		o.Brand.Logo = brandLogo.String
		o.Brand.Name = brandName.String

		orderHistory = append(orderHistory, o)
	}

	return orderHistory, nil
}

func (repo *orderRepository) CountOrders(ctx context.Context, since time.Time, userId uint32, status string) (int64, error) {
	var count int64
	statusQuery, statusParams := model.BuildOrderStatusQuery(status)
	query := "SELECT COUNT(o.accesstrade_order_id) " +
		"FROM aff_order AS o " +
		"LEFT JOIN aff_reward AS r ON r.accesstrade_order_id = o.accesstrade_order_id " +
		"WHERE o.user_id = ? AND o.created_at > ? "
	if status != "" {
		query += statusQuery
	}

	var err error
	if status != "" {
		err = repo.db.Raw(query, userId, since, statusParams).Count(&count).Error
	} else {
		err = repo.db.Raw(query, userId, since).Count(&count).Error
	}
	return count, err
}

func (repo *orderRepository) GetCampaignByATId(atId string) (*model.AffCampaign, error) {
	var camp model.AffCampaign
	err := repo.db.First(&camp, "accesstrade_id = ?", atId).Error
	if err != nil {
		return nil, err
	}
	return &camp, nil
}
