package campaign

import (
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade/types"
	model2 "github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/datatypes"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CampaignRepository struct {
	Db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) *CampaignRepository {
	return &CampaignRepository{
		Db: db,
	}
}

func (repo *CampaignRepository) RetrieveCampaignsByAccessTradeIds(ids []string) (map[string]*model2.AffCampaign, error) {
	var data []model2.AffCampaign
	err := repo.Db.Table("aff_campaign").
		Joins("Description").
		Where("accesstrade_id IN ?", ids).Find(&data).Error
	if err != nil {
		return nil, err
	}
	mapped := map[string]*model2.AffCampaign{}
	for _, it := range data {
		mapped[it.AccessTradeId] = &it
	}

	return mapped, nil
}

func (repo *CampaignRepository) SaveATCampaign(atCampaign *types.ATCampaign) error {
	// First create campaign
	fmt.Println("SaveATCampaign", atCampaign.Id)
	newCampaign := model2.AffCampaign{
		ActiveStatus:   1,
		AccessTradeId:  atCampaign.Id,
		Logo:           atCampaign.Logo,
		MaxCom:         atCampaign.MaxCom,
		Merchant:       atCampaign.Merchant,
		Name:           atCampaign.Name,
		Scope:          atCampaign.Scope,
		Approval:       atCampaign.Approval,
		Status:         atCampaign.Status,
		Type:           atCampaign.Type,
		Url:            atCampaign.Url,
		Category:       atCampaign.Category,
		SubCategory:    atCampaign.SubCategory,
		CookiePolicy:   atCampaign.CookiePolicy,
		CookieDuration: atCampaign.CookieDuration,
		StellaStatus:   "DRAFT",
		StellaDescription: datatypes.JSON(`{
			"action_point": "",
			"commission_policy": "",
			"introduction": "",
			"other_notice": "",
			"rejected_reason": "",
			"traffic_building_policy": ""}`),
	}
	if !atCampaign.StartTime.IsZero() {
		newCampaign.StartTime = &atCampaign.StartTime.Time
	}
	if !atCampaign.EndTime.IsZero() {
		newCampaign.EndTime = &atCampaign.EndTime.Time
	}

	campaignDescription := model2.CampaignDescription{
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		ActionPoint:           atCampaign.Description.ActionPoint,
		CommissionPolicy:      atCampaign.Description.CommissionPolicy,
		CookiePolicy:          atCampaign.Description.CookiePolicy,
		Introduction:          atCampaign.Description.Introduction,
		OtherNotice:           atCampaign.Description.OtherNotice,
		RejectedReason:        atCampaign.Description.RejectedReason,
		TrafficBuildingPolicy: atCampaign.Description.TrafficBuildingPolicy,
	}

	err := repo.Db.Transaction(func(tx *gorm.DB) error {
		// First find brand, and create if not exist
		err := tx.Create(&newCampaign).Error
		if err != nil {
			return err
		}
		campaignDescription.CampaignId = newCampaign.ID
		err = tx.Create(&campaignDescription).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.Errorf("campaign tx error: %v", err)
	}

	return nil
}

func (repo *CampaignRepository) RetrieveCampaigns(q map[string]any) ([]model2.AffCampaign, error) {
	var data []model2.AffCampaign
	err := repo.Db.Preload("Description").Find(&data, q).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *CampaignRepository) CreateCampaigns(data []model2.AffCampaign) ([]model2.AffCampaign, error) {
	err := repo.Db.Create(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *CampaignRepository) UpdateCampaigns(data []model2.AffCampaign) ([]model2.AffCampaign, error) {
	err := repo.Db.Updates(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *CampaignRepository) DeactivateCampaigns(data []model2.AffCampaign) error {
	err := repo.Db.Updates(&data).Error

	if err != nil {
		return err
	}
	return nil
}

func (repo *CampaignRepository) RetrieveAffLinks(campaignId uint, originalUrl string) ([]model2.AffLink, error) {
	var links []model2.AffLink
	m := repo.Db.Model(&links)
	if originalUrl != "" {
		m.Where("url_origin = ?", originalUrl)
	}
	m.Where("campaign_id = ? AND active_status = ?", campaignId, model2.AFF_LINK_STATUS_ACTIVE)
	err := m.Find(&links).Error

	if err != nil {
		return nil, err
	}
	return links, nil
}

func (repo *CampaignRepository) CreateAffLinks(data []model2.AffLink) error {
	err := repo.Db.Create(&data).Error
	return err
}

func (repo *CampaignRepository) CreateTrackedClick(m *model2.AffTrackedClick) error {
	err := repo.Db.Create(&m).Error
	return err
}
