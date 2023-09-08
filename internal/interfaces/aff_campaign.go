package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

type AffCampaignRepository interface {
	GetAllAffCampaign(ctx context.Context, page, size int) ([]model.AffCampaign, error)
	GetAffCampaignByAccesstradeId(ctx context.Context, accesstradeId uint64) (model.AffCampaign, error)
}
