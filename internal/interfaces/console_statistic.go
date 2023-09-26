package interfaces

import (
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type StatisticRepository interface {
	FindOrdersInRange(d dto.TimeRange, offset int, limit int) ([]model.AffOrder, error)
	TotalOrdersInRange(d dto.TimeRange) (int64, error)
	CalculateCashbackInRange(d dto.TimeRange) (dto.Cashback, error)
	TotalActiveCampaignsInRange(d dto.TimeRange) (int64, error)
}

type StatisticUcase interface {
	GetSummaryByTimeRange(d dto.TimeRange) (*dto.StatisticSummaryResponse, error)
}
