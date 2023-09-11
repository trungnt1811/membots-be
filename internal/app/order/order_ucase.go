package order

import (
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
)

type OrderUcase struct {
	Repo   interfaces.OrderRepository
	ATRepo interfaces.ATRepository
}

func NewOrderUcase(repo interfaces.OrderRepository, atRepo interfaces.ATRepository) *OrderUcase {
	return &OrderUcase{
		Repo:   repo,
		ATRepo: atRepo,
	}
}

func (u *OrderUcase) PostBackUpdateOrder(postBackReq *dto.ATPostBackRequest) (*model.AffOrder, error) {
	// First, save log
	err := u.Repo.SavePostBackLog(&model.AffPostBackLog{
		OrderId:   postBackReq.OrderId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Data:      postBackReq,
	})
	if err != nil {
		log.LG.Errorf("save aff_post_back_log error: %v", err)
	}

	// Parse sale time for later request
	salesTime, err := util.ParsePostBackTime(postBackReq.SalesTime)
	if err != nil {
		return nil, fmt.Errorf("cannot parse sales_time: %v", err)
	}
	since, until := util.GetSinceUntilTime(salesTime, 1)

	// Then, find if order is exist
	resp, err := u.ATRepo.QueryOrders(types.ATOrderQuery{
		Since: since,
		Until: until,
	}, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("query order-list error: %v", err)
	}
	var atOrder *types.ATOrder
	for _, item := range resp.Data {
		if item.OrderId == postBackReq.OrderId {
			atOrder = &item
		}
	}
	if atOrder == nil {
		return nil, fmt.Errorf("not found order \"%s\" in order-list", postBackReq.OrderId)
	}

	order, err := u.Repo.FindOrderByAccessTradeId(atOrder.OrderId)
	if err != nil {
		if err.Error() == "record not found" {
			// Order not exist, create new one
			order = model.NewOrderFromATOrder(atOrder)
			crErr := u.Repo.CreateOrder(order)
			if crErr != nil {
				return nil, fmt.Errorf("create order failed: %v", crErr)
			}
			return order, nil
		} else {
			return nil, err
		}
	} else {
		// Or update exist order
		updated := model.NewOrderFromATOrder(atOrder)
		updated.ID = order.ID
		_, err = u.Repo.UpdateOrder(updated)
		if err != nil {
			return nil, fmt.Errorf("update order failed: %v", err)
		}
	}

	// TODO: After create order, push Kafka msg for sync transactions
	_, err = u.SyncTransactionsByOrder(order.AccessTradeOrderId)
	if err != nil {
		log.LG.Errorf("sync txs failed: %v", err)
	}
	// TODO: Send Kafka message if order approved

	return order, nil
}

// The `SyncTransactionsByOrder` function is responsible for synchronizing transactions for a given
// order.
func (u *OrderUcase) SyncTransactionsByOrder(atOrderId string) (int, error) {
	// First find created order
	order, err := u.Repo.FindOrderByAccessTradeId(atOrderId)
	if err != nil {
		return 0, fmt.Errorf("find order id \"%s\" failed: %v", atOrderId, err)
	}
	since, until := util.GetSinceUntilTime(order.SalesTime, 1)
	txs, err := u.ATRepo.QueryTransactions(types.ATTransactionQuery{
		Since:         since,
		Until:         until,
		TransactionId: atOrderId, // NOTE: post back req will return old transaction_id as order_id!
	}, 0, 0)
	if err != nil {
		return 0, fmt.Errorf("query txs failed: %v", err)
	}
	if len(txs.Data) == 0 {
		return 0, fmt.Errorf("txs empty for order: %s", atOrderId)
	}

	// Save transactions if found
	transactions := make([]model.AffTransaction, len(txs.Data))
	for i, it := range txs.Data {
		transactions[i] = *model.NewAffTransactionFromAT(order, &it)
	}

	err = u.Repo.UpdateOrCreateATTransactions(transactions)
	if err != nil {
		return 0, fmt.Errorf("save txs failed: %v", err)
	}

	return len(txs.Data), nil
}
