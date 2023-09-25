package test

import (
	"context"
	"testing"
	"time"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/order"
	"github.com/astraprotocol/affiliate-system/internal/app/reward"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/worker/kafkaconsumer"
	"github.com/stretchr/testify/assert"
)

type MockAstraPriceRepo struct {
	AstraPrice int64
}

func (m MockAstraPriceRepo) GetAstraPrice(ctx context.Context) (int64, error) {
	return m.AstraPrice, nil
}

type orderRewardTest struct {
	AffCommission    float64
	StellaCommission float64
	PriceRepo        interfaces.TokenPriceRepo
	RewardAmount     float64
}

var calculateOrderRewardTestSet = []orderRewardTest{
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

// ==========================================================================
// SECTION: test withdraw amount calculating
// ===========================================================================

type withdrawRewardTest struct {
	Reward       model.Reward
	ResultAmount float64
}

func initTestReward(amount, rewardedAmount float64, timeSinceCreated time.Duration, lockTime time.Duration) model.Reward {
	startAt := time.Now().Add(-1 * timeSinceCreated)
	return model.Reward{
		Amount:         amount,
		RewardedAmount: rewardedAmount,
		StartAt:        startAt,
		EndAt:          startAt.Add(lockTime),
	}
}

var withdrawRewardTestSet = []withdrawRewardTest{
	// rewarded = 0
	{
		initTestReward(100, 0, 0*model.OneDay, 60*model.OneDay),
		50,
	},
	{
		initTestReward(100, 0, 3*model.OneDay+2*time.Minute, 60*model.OneDay),
		52.5,
	},
	{
		initTestReward(100, 0, 30*model.OneDay+2*time.Minute, 60*model.OneDay),
		75,
	},
	{
		initTestReward(100, 0, 100*model.OneDay+2*time.Minute, 60*model.OneDay),
		100,
	},
	// rewarded != 0
	{
		initTestReward(100, 0, 0*model.OneDay, 60*model.OneDay),
		50,
	},
	{
		initTestReward(100, 5, 3*model.OneDay+2*time.Minute, 60*model.OneDay),
		47.5,
	},
	{
		initTestReward(100, 10.11, 30*model.OneDay+2*time.Minute, 60*model.OneDay),
		64.89,
	},
	{
		initTestReward(100, 99, 100*model.OneDay+2*time.Minute, 60*model.OneDay),
		1,
	},
}

func Test_CalculateWithdrawableReward(t *testing.T) {
	for _, test := range withdrawRewardTestSet {
		withdrawAmount, _ := test.Reward.WithdrawableReward()
		assert.Equal(t, test.ResultAmount, withdrawAmount)
	}
}
