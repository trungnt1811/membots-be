package campaign

import (
	"errors"
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/dto"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/model"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util"
)

type CampaignUsecase struct {
	Repo   interfaces.CampaignRepository
	ATRepo interfaces.ATRepository
}

func NewCampaignUsecase(repo interfaces.CampaignRepository, atRepo interfaces.ATRepository) *CampaignUsecase {
	return &CampaignUsecase{
		Repo:   repo,
		ATRepo: atRepo,
	}
}

func (u *CampaignUsecase) GenerateAffLink(user *model.UserEntity, payload *dto.CreateLinkPayload) (*dto.CreateLinkResponse, error) {
	// First query the campaign
	campaigns, err := u.Repo.RetrieveCampaigns(map[string]any{
		"id": payload.CampaignId,
	})
	if err != nil {
		return nil, err
	}
	if len(campaigns) == 0 {
		return nil, errors.New("wrong campaign id")
	}

	campaign := campaigns[0]

	// Then find campaign link if exist
	affLinks, err := u.Repo.RetrieveAffLinks(campaign.ID)

	if len(affLinks) == 0 {
		// If campaign link not available, request to generate new one
		resp, err := u.ATRepo.CreateTrackingLinks(campaign.AccesstradeId, []string{
			payload.OriginalUrl,
		}, map[string]string{})
		if err != nil {
			return nil, fmt.Errorf("fail to create aff link: %v", err)
		}

		for _, atLink := range resp.Data.SuccessLink {
			affLinks = append(affLinks, model.AffLink{
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				AffLink:      atLink.AffLink,
				FirstLink:    atLink.FirstLink,
				ShortLink:    atLink.ShortLink,
				UrlOrigin:    atLink.UrlOrigin,
				CampaignId:   campaign.ID,
				ActiveStatus: model.AFF_LINK_STATUS_ACTIVE,
			})
		}

		err = u.Repo.CreateAffLinks(affLinks)
		if err != nil {
			return nil, fmt.Errorf("save created aff link fail: %v", err)
		}
		if len(affLinks) == 0 {
			return nil, fmt.Errorf("aff links empty")
		}
	}
	link := affLinks[0]

	// TODO: Create tracked tx

	// Add user id and other params
	additionalParams := map[string]string{
		"utm_content": fmt.Sprint(user.ID),
		"stella_tx":   fmt.Sprint("tx_"),
	}

	linkResp := dto.CreateLinkResponse{
		CampaignId:  link.CampaignId,
		AffLink:     util.PackQueryParamsToUrl(link.AffLink, additionalParams),
		ShortLink:   util.PackQueryParamsToUrl(link.ShortLink, additionalParams),
		OriginalUrl: link.UrlOrigin,
	}

	return &linkResp, nil
}
