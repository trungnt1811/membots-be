package campaign

import (
	"errors"
	"fmt"
	"time"

	interfaces2 "github.com/astraprotocol/affiliate-system/internal/interfaces"
	model2 "github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"

	"github.com/astraprotocol/affiliate-system/internal/dto"
)

type CampaignUsecase struct {
	Repo   interfaces2.CampaignRepository
	ATRepo interfaces2.ATRepository
}

func NewCampaignUsecase(repo interfaces2.CampaignRepository, atRepo interfaces2.ATRepository) *CampaignUsecase {
	return &CampaignUsecase{
		Repo:   repo,
		ATRepo: atRepo,
	}
}

func (u *CampaignUsecase) GenerateAffLink(userId uint64, payload *dto.CreateLinkPayload) (*dto.CreateLinkResponse, error) {
	// First query the campaignLess
	campaignLess, err := u.Repo.GetCampaignLessById(payload.CampaignId)
	if err != nil {
		return nil, err
	}

	// Then find campaignLess link if exist
	isJustCreated := false
	affLinks, err := u.Repo.RetrieveAffLinks(campaignLess.ID, payload.OriginalUrl)
	if err != nil {
		return nil, fmt.Errorf("retrieve aff link fail: %v", err)
	}
	if len(affLinks) == 0 {
		// If campaignLess link not available, request to generate new one
		urls := []string{}
		if payload.OriginalUrl != "" {
			urls = append(urls, payload.OriginalUrl)
		}
		resp, err := u.ATRepo.CreateTrackingLinks(campaignLess.AccessTradeId, payload.ShortenLink, urls, map[string]string{})
		if err != nil {
			return nil, fmt.Errorf("fail to generate: %v", err)
		}

		for _, atLink := range resp.Data.SuccessLink {
			affLinks = append(affLinks, model2.AffLink{
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				AffLink:      atLink.AffLink,
				FirstLink:    atLink.FirstLink,
				ShortLink:    atLink.ShortLink,
				UrlOrigin:    atLink.UrlOrigin,
				CampaignId:   campaignLess.ID,
				ActiveStatus: model2.AFF_LINK_STATUS_ACTIVE,
			})
		}
		if len(affLinks) == 0 {
			return nil, errors.New("empty access trade links resp")
		}

		err = u.Repo.CreateAffLinks(affLinks)
		if err != nil {
			return nil, fmt.Errorf("save created aff link fail: %v", err)
		}
		if len(affLinks) == 0 {
			return nil, fmt.Errorf("aff links empty")
		}
		isJustCreated = true
	}
	link := affLinks[0]

	// Create tracked click item
	tracked := &model2.AffTrackedClick{
		UserId:     uint(userId),
		CampaignId: campaignLess.ID,
		LinkId:     link.ID,
		AffLink:    link.AffLink,
		UrlOrigin:  payload.OriginalUrl,
		CreatedAt:  time.Now(),
	}
	err = u.Repo.CreateTrackedClick(tracked)
	if err != nil {
		log.LG.Errorf("create tracked click failed: %v", err)
	}
	// Add user id and other params
	additionalParams := map[string]string{
		"utm_content": util.StringifyUTMContent(uint(userId), tracked.ID),
	}
	clickLink := util.PackQueryParamsToUrl(link.AffLink, additionalParams)

	linkResp := dto.CreateLinkResponse{
		CampaignId:  link.CampaignId,
		AffLink:     clickLink,
		OriginalUrl: link.UrlOrigin,
		BrandNew:    isJustCreated,
	}

	return &linkResp, nil
}
