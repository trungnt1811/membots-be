package aff_search

import (
	"context"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
)

type affSearchUCase struct {
	Repo         interfaces.AffSearchRepository
	ConvertPrice interfaces.ConvertPriceHandler
}

func (a *affSearchUCase) Search(ctx context.Context, q string, page, size int) (dto.AffSearchResponseDto, error) {
	results, err := a.Repo.Search(ctx, q, page, size)
	if err != nil {
		return dto.AffSearchResponseDto{}, err
	}
	var affCampaigns []dto.AffCampaignLessDto

	for i := range results.AffCampaign {
		if i >= size {
			break
		}
		tmpCampaign := results.AffCampaign[i]
		tmpCampaign.StellaMaxCom = a.ConvertPrice.GetStellaMaxCommission(ctx, tmpCampaign.Attributes)
		affCampaigns = append(affCampaigns, tmpCampaign.ToDto())
	}
	nextPage := page
	if len(results.AffCampaign) > size {
		nextPage = page + 1
	}
	return dto.AffSearchResponseDto{
		NextPage: nextPage,
		Page:     page,
		Size:     size,
		Data: dto.AffSearchDto{
			AffCampaigns:  affCampaigns,
			TotalCampaign: results.TotalCampaign,
		},
	}, nil
}

func NewAffSearchUCase(repo interfaces.AffSearchRepository, convertPrice interfaces.ConvertPriceHandler) interfaces.AffSearchUCase {
	return &affSearchUCase{
		Repo:         repo,
		ConvertPrice: convertPrice,
	}
}
