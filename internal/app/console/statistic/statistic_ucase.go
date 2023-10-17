package statistic

import (
	"fmt"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
)

const (
	BATCH_SIZE = 50
)

type StatisticUCase struct {
	Repo interfaces.StatisticRepository
}

func NewStatisticUCase(repo interfaces.StatisticRepository) *StatisticUCase {
	return &StatisticUCase{
		Repo: repo,
	}
}

func (ucase *StatisticUCase) GetSummaryByTimeRange(d dto.TimeRange) (*dto.StatisticSummaryResponse, error) {
	var resp dto.StatisticSummaryResponse

	// First, get total order by range
	count, err := ucase.Repo.TotalOrdersInRange(0, d)
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
		orders, err := ucase.Repo.FindOrdersInRange(0, d, offset, BATCH_SIZE)
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
			if order.OrderStatus == model.OrderStatusPending || order.OrderStatus == model.OrderStatusApproved || order.OrderStatus == model.OrderStatusRewarding || order.OrderStatus == model.OrderStatusComplete {
				rev += float64(order.PubCommission)
			}
		}

		offset += BATCH_SIZE
	}

	resp.NumOfCustomers = len(customers)
	resp.NumOfOrders = int(count)
	resp.TotalRevenue = rev

	total, err := ucase.Repo.CalculateCashbackInRange(0, d)
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

func (ucase *StatisticUCase) GetCampaignSummaryByTimeRange(campaignId uint, d dto.TimeRange) (*dto.CampaignSummaryResponse, error) {
	var resp dto.CampaignSummaryResponse
	// First, get total order by range
	count, err := ucase.Repo.TotalOrdersInRange(campaignId, d)
	if err != nil {
		return nil, fmt.Errorf("count orders err: %v", err)
	}

	offset := int(0)
	customers := map[uint]int{}
	var rev float64 = 0
	ordersByStatus := map[string]int{
		model.OrderStatusPending:   0,
		model.OrderStatusApproved:  0,
		model.OrderStatusRejected:  0,
		model.OrderStatusCancelled: 0,
		model.OrderStatusRewarding: 0,
		model.OrderStatusComplete:  0,
	}
	for {
		if offset >= int(count) {
			break
		}

		// Process orders by batch
		orders, err := ucase.Repo.FindOrdersInRange(campaignId, d, offset, BATCH_SIZE)
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
			if _, ok := ordersByStatus[order.OrderStatus]; ok {
				ordersByStatus[order.OrderStatus] += 1
			} else {
				ordersByStatus[order.OrderStatus] = 1
			}
			if order.OrderStatus == model.OrderStatusPending || order.OrderStatus == model.OrderStatusApproved || order.OrderStatus == model.OrderStatusRewarding || order.OrderStatus == model.OrderStatusComplete {
				rev += float64(order.PubCommission)
			}
		}

		offset += BATCH_SIZE
	}

	resp.NumOfCustomers = len(customers)
	resp.NumOfOrders = int(count)
	resp.TotalRevenue = rev
	resp.OrdersByStatus = ordersByStatus

	total, err := ucase.Repo.CalculateCashbackInRange(campaignId, d)
	if err != nil {
		log.LG.Errorf("calculate cashback err: %v", err)
	}
	resp.TotalASACashback.Distributed = total.Distributed
	resp.TotalASACashback.Remain = total.Remain

	return &resp, nil
}
