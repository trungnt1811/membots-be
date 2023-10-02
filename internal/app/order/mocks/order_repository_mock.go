package mocks

import (
	"context"
	"errors"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

type MockOrderRepository struct {
	Orders []model.AffOrder
	Logs   []model.AffPostBackLog
	Txs    []model.AffTransaction
}

func NewMockOrderRepository(orders []model.AffOrder) *MockOrderRepository {
	return &MockOrderRepository{
		Orders: orders,
		Logs:   []model.AffPostBackLog{},
		Txs:    []model.AffTransaction{},
	}
}

func (repo *MockOrderRepository) SavePostBackLog(req *model.AffPostBackLog) error {
	req.ID = uint(len(repo.Logs) + 1)
	repo.Logs = append(repo.Logs, *req)
	return nil
}

func (repo *MockOrderRepository) CreateOrder(order *model.AffOrder) error {
	order.ID = uint(len(repo.Orders) + 1)
	repo.Orders = append(repo.Orders, *order)
	return nil
}

func (repo *MockOrderRepository) UpdateOrder(updated *model.AffOrder) (int, error) {
	var effected = 0
	for i, order := range repo.Orders {
		if order.AccessTradeOrderId == updated.AccessTradeOrderId {
			repo.Orders[i] = *updated
			effected += 1
		}
	}

	return effected, nil
}

func (repo *MockOrderRepository) UpdateOrCreateATTransactions(txs []model.AffTransaction) error {
	repo.Txs = append(repo.Txs, txs...)
	return nil
}

func (repo *MockOrderRepository) FindOrderByAccessTradeId(atOrderId string) (*model.AffOrder, error) {
	var found *model.AffOrder
	for _, order := range repo.Orders {
		if order.AccessTradeOrderId == atOrderId {
			found = &order
		}
	}

	if found == nil {
		return nil, errors.New("record not found")
	}

	return found, nil
}

func (repo *MockOrderRepository) UpdateTrackedClickOrder(trackedId uint64, order *model.AffOrder) error {
	return nil
}

func (repo *MockOrderRepository) GetOrderDetails(ctx context.Context, userId uint32, orderId uint) (*model.OrderDetails, error) {
	return nil, nil
}

func (repo *MockOrderRepository) GetOrderHistory(ctx context.Context, since time.Time, userId uint32, status string, page, size int) ([]model.OrderDetails, error) {
	return []model.OrderDetails{}, nil
}
func (repo *MockOrderRepository) CountOrders(ctx context.Context, since time.Time, userId uint32, status string) (int64, error) {
	return 0, nil
}

func (repo *MockOrderRepository) GetCampaignByATId(id string) (*model.AffCampaign, error) {
	return &model.AffCampaign{}, nil
}
