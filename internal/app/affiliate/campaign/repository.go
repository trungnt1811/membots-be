package campaign

import (
	"fmt"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/model"
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

func (repo *CampaignRepository) RetrieveCampaignsByAccessTradeIds(ids []string) (map[string]*model.AffCampaign, error) {
	var data []model.AffCampaign
	err := repo.Db.Table("aff_campaign").
		Joins("Description").
		Where("accesstrade_id IN ?", ids).Find(&data).Error
	if err != nil {
		return nil, err
	}
	mapped := map[string]*model.AffCampaign{}
	for _, it := range data {
		mapped[it.AccessTradeId] = &it
	}

	return mapped, nil
}

func (repo *CampaignRepository) SaveATCampaign(atCampaign *types.ATCampaign) error {
	// First create campaign
	fmt.Println("SaveATCampaign", atCampaign.Id)
	newCampaign := model.AffCampaign{
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
	}
	if !atCampaign.StartTime.IsZero() {
		newCampaign.StartTime = &atCampaign.StartTime.Time
	}
	if !atCampaign.EndTime.IsZero() {
		newCampaign.EndTime = &atCampaign.EndTime.Time
	}

	campaignDescription := model.CampaignDescription{
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

func (repo *CampaignRepository) RetrieveCampaigns(q map[string]any) ([]model.AffCampaign, error) {
	var data []model.AffCampaign
	err := repo.Db.Preload("Description").Find(&data, q).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *CampaignRepository) CreateCampaigns(data []model.AffCampaign) ([]model.AffCampaign, error) {
	err := repo.Db.Create(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *CampaignRepository) UpdateCampaigns(data []model.AffCampaign) ([]model.AffCampaign, error) {
	err := repo.Db.Updates(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *CampaignRepository) DeactivateCampaigns(data []model.AffCampaign) error {
	err := repo.Db.Updates(&data).Error

	if err != nil {
		return err
	}
	return nil
}

func (repo *CampaignRepository) RetrieveAffLinks(campaignId uint) ([]model.AffLink, error) {
	var links []model.AffLink
	err := repo.Db.Find(&links, "campaign_id = ? AND active_status = ?", campaignId, model.AFF_LINK_STATUS_ACTIVE).Error

	if err != nil {
		return nil, err
	}
	return links, nil
}

func (repo *CampaignRepository) CreateAffLinks(data []model.AffLink) error {
	err := repo.Db.Create(&data).Error
	return err
}
