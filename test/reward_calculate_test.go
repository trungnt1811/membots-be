package test

import (
	"context"
	"testing"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/worker/kafkaconsumer"
	"github.com/stretchr/testify/assert"
)

type MockAstraPriceRepo struct {
	AstraPrice int64
}

func (m MockAstraPriceRepo) GetAstraPrice(ctx context.Context) (int64, error) {
	return m.AstraPrice, nil
}

type addTest struct {
	AffCommission    float64
	StellaCommission float64
	PriceRepo        interfaces.TokenPriceRepo
	RewardAmount     float64
}

var calculateOrderRewardTestSet = []addTest{
	{10000, 0, MockAstraPriceRepo{200}, 50},
	{2505, 0, MockAstraPriceRepo{200}, 12.53},
	{10000, 10, MockAstraPriceRepo{200}, 45},
	{2505, 10, MockAstraPriceRepo{200}, 11.27},
	{10000, 10, MockAstraPriceRepo{140}, 64.29},
	{2505, 10, MockAstraPriceRepo{140}, 16.1},
}

func Test_CalculateOrderReward(t *testing.T) {
	db := conf.DBConn()
	rewardRepo := reward.NewRewardRepository(db)
	orderRepo := order.NewOrderRepository(db)

	for _, test := range calculateOrderRewardTestSet {
		rewardMaker := kafkaconsumer.NewRewardMaker(rewardRepo, orderRepo, test.PriceRepo, nil, nil)
		rewardAmt, err := rewardMaker.CalculateRewardAmt(test.AffCommission, test.StellaCommission)
		assert.Nil(t, err)
		assert.Equal(t, test.RewardAmount, rewardAmt)
	}
}
