package interfaces

import (
	"context"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

type ConvertPriceHandler interface {
	ConvertVndPriceToAstra(ctx context.Context, attribute model.AffCampaignAttribute) string
	GetStellaMaxCommission(ctx context.Context, attributes []model.AffCampaignAttribute) string
}
