package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type StatisticRepository interface {
	FindOrdersInRange(campaignId uint, d dto.TimeRange, offset int, limit int) ([]model.AffOrder, error)
	TotalOrdersInRange(campaignId uint, d dto.TimeRange) (int64, error)
	CalculateCashbackInRange(campaignId uint, d dto.TimeRange) (dto.Cashback, error)
	TotalActiveCampaignsInRange(d dto.TimeRange) (int64, error)
}

type StatisticUCase interface {
	GetSummaryByTimeRange(d dto.TimeRange) (*dto.StatisticSummaryResponse, error)
	GetCampaignSummaryByTimeRange(campaignId uint, d dto.TimeRange) (*dto.CampaignSummaryResponse, error)
}
