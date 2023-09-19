package statistic

import (
	"fmt"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
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
	var rev float32 = 0
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
			rev += order.PubCommission
		}

		offset += BATCH_SIZE
	}

	resp.NumOfCustomers = len(customers)
	resp.NumOfOrders = int(count)
	resp.TotalRevenue = rev

	// TODO: Calculate total asa cashback

	return &resp, nil
}
