package statistic

import (
	"fmt"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
)

const (
	BATCH_SIZE = 50
)

type StatisticUcase struct {
	Repo interfaces.StatisticRepository
}

func NewStatisticUcase(repo interfaces.StatisticRepository) *StatisticUcase {
	return &StatisticUcase{
		Repo: repo,
	}
}

func (ucase *StatisticUcase) GetSummaryByTimeRange(d dto.TimeRange) (*dto.StatisticSummaryResponse, error) {
	var resp dto.StatisticSummaryResponse

	// First, get total order by range
	count, err := ucase.Repo.TotalOrdersInRange(d)
	if err != nil {
		return nil, fmt.Errorf("count orders err: %v", err)
	}

	offset := int(0)
	customers := map[uint]int{}
	var rev float64 = 0
	for {
		if offset >= int(count) {
			break
		}

		// Process orders by batch
		orders, err := ucase.Repo.FindOrdersInRange(d, offset, BATCH_SIZE)
		if err != nil {
			return nil, fmt.Errorf("find orders err: %v", err)
		}

		// Then, find total customer by range and calculate revenue
		for _, order := range orders {
			if prev, ok := customers[order.UserId]; !ok {
				customers[order.UserId] = 1
			} else {
				customers[order.UserId] = prev + 1
			}
			rev += float64(order.PubCommission)
		}

		offset += BATCH_SIZE
	}

	resp.NumOfCustomers = len(customers)
	resp.NumOfOrders = int(count)
	resp.TotalRevenue = rev

	total, err := ucase.Repo.CalculateCashbackInRange(d)
	if err != nil {
		log.LG.Errorf("calculate cashback err: %v", err)
	}
	resp.TotalASACashback.Distributed = total.Distributed
	resp.TotalASACashback.Remain = total.Remain

	totalCampaigns, err := ucase.Repo.TotalActiveCampaignsInRange(d)
	if err != nil {
		log.LG.Errorf("count active campaigns err: %v", err)
	}
	resp.ActiveCampaigns = int(totalCampaigns)

	return &resp, nil
}
