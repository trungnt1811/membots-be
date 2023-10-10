package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"

	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/order/mocks"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	atMocks "github.com/astraprotocol/affiliate-system/internal/infra/accesstrade/mocks"
	"github.com/astraprotocol/affiliate-system/internal/model"
	logger "github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type OrderUcaseTestSuite struct {
	suite.Suite
	repoMock *mocks.MockOrderRepository
	uCase    interfaces.OrderUCase
}

func NewOrderUcaseTestSuite() *OrderUcaseTestSuite {
	logger.LG = logger.NewZerologLogger(os.Stdout, zerolog.InfoLevel)
	atRepo := atMocks.NewAccessTradeRepoMock()
	repo := mocks.NewMockOrderRepository([]model.AffOrder{})
	ucase := order.NewOrderUCase(repo, atRepo, nil)
	return &OrderUcaseTestSuite{
		uCase:    ucase,
		repoMock: repo,
	}
}

func TestRunOrderUcaseSuite(t *testing.T) {
	suite.Run(t, NewOrderUcaseTestSuite())
}

func (s *OrderUcaseTestSuite) TestPostBackReceived() {
	var postBackReq dto.ATPostBackRequest
	err := json.Unmarshal([]byte(atMocks.SAMPLE_POST_BACK), &postBackReq)
	s.NoError(err)

	// run post back received handle
	m, err := s.uCase.PostBackUpdateOrder(&postBackReq)
	s.NoError(err)
	s.Equal(uint(1), m.ID)
	s.Equal(uint(1), m.UserId)
	s.Equal(1, len(s.repoMock.Orders))
	s.Equal(1, len(s.repoMock.Logs))
}

func (s *OrderUcaseTestSuite) TestSyncTransactionsByOrder() {
	var postBackReq dto.ATPostBackRequest
	err := json.Unmarshal([]byte(atMocks.SAMPLE_POST_BACK), &postBackReq)
	s.NoError(err)

	// run post back received handle first
	_, err = s.uCase.PostBackUpdateOrder(&postBackReq)
	s.NoError(err)

	nSynced, err := s.uCase.SyncTransactionsByOrder(postBackReq.OrderId)
	s.NoError(err)
	s.Equal(4, nSynced)
}

func (s *OrderUcaseTestSuite) TestCheckOrderConfirmed() {
	// First, try create order
	affOrder := model.AffOrder{
		AccessTradeOrderId: "230824H3K7H4V5",
		OrderStatus:        "approved",
		IsConfirmed:        0,
	}

	err := s.repoMock.CreateOrder(&affOrder)
	s.NoError(err)

	// run post back received handle first
	updatedCount, err := s.uCase.CheckOrderConfirmed()
	s.NoError(err)
	s.Equal(1, updatedCount)
}

func (s *OrderUcaseTestSuite) TestCheckOrderListAndSync() {
	// run check order with one mocked order not synced
	createdCount, err := s.uCase.CheckOrderListAndSync()
	s.NoError(err)
	s.Equal(1, createdCount)
}
