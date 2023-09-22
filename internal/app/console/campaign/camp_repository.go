package campaign

import (
	"encoding/json"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type affCampaignRepository struct {
	Db *gorm.DB
}

func (a *affCampaignRepository) UpdateCampaignAttribute(id uint, attributes []model.AffCampaignAttribute) error {
	var listAttributeInDB []model.AffCampaignAttribute
	if err := a.Db.Where("campaign_id = ?", id).Find(&listAttributeInDB).Error; err != nil {
		return err
	}
	mapIdInput := make(map[uint]uint)
	for _, attribute := range attributes {
		mapIdInput[attribute.ID] = attribute.ID
	}
	var listKeyDelete []uint
	for _, attribute := range listAttributeInDB {
		_, ok := mapIdInput[attribute.ID]
		if !ok {
			listKeyDelete = append(listKeyDelete, attribute.ID)
		}
	}
	var listCreate []model.AffCampaignAttribute
	var listUpdate []model.AffCampaignAttribute
	for _, attribute := range attributes {
		if attribute.ID == 0 {
			listCreate = append(listCreate, attribute)
		} else {
			listUpdate = append(listUpdate, attribute)
		}
	}
	return a.Db.Transaction(func(tx *gorm.DB) error {
		if len(listCreate) > 0 {
			if err := tx.Create(&listCreate).Error; err != nil {
				return err
			}
		}
		if len(listUpdate) > 0 {
			for _, update := range listUpdate {
				tx.Updates(&update)
			}
			if err := tx.Error; err != nil {
				return err
			}
		}
		if len(listKeyDelete) > 0 {
			if err := tx.Delete(&model.AffCampaignAttribute{}, listKeyDelete).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *affCampaignRepository) GetCampaignLessByAccessTradeId(accessTradeId string) (model.AffCampaignLess, error) {
	var affCampaign model.AffCampaignLess
	if err := a.Db.Where("accesstrade_id = ?", accessTradeId).First(&affCampaign).Error; err != nil {
		return affCampaign, err
	}
	return affCampaign, nil
}

func (a *affCampaignRepository) GetCampaignById(id uint) (model.AffCampaign, error) {
	var affCampaign model.AffCampaign
	if err := a.Db.Table("aff_campaign").
		Joins("Description").
		Joins("Brand").
		Preload("Attributes").
		Joins("StellaCategory").
		Where("aff_campaign.id = ?", id).First(&affCampaign).Error; err != nil {
		return affCampaign, err
	}
	return affCampaign, nil
}

func (a *affCampaignRepository) UpdateCampaign(id uint, updates map[string]interface{}) error {
	stellaDescription, ok := updates["stella_description"]
	// If the not key exists
	if !ok {
		return a.Db.Table("aff_campaign").Where("id = ?", id).Updates(updates).Error
	}
	var affCampaign model.AffCampaignApp
	if err := a.Db.Table("aff_campaign").
		Where("id = ?", id).First(&affCampaign).Error; err != nil {
		return err
	}
	var stellaDescriptionJsonInDB model.StellaDescriptionJson
	err := json.Unmarshal(affCampaign.StellaDescription, &stellaDescriptionJsonInDB)
	if err != nil {
		return err
	}
	return a.Db.Transaction(func(tx *gorm.DB) error {
		stellaDescriptionJsonInput, _ := stellaDescription.(map[string]interface{})
		actionPoint, ok := stellaDescriptionJsonInput["action_point"]
		if ok {
			stellaDescriptionJsonInDB.ActionPoint = actionPoint.(string)
		}
		introduction, ok := stellaDescriptionJsonInput["introduction"]
		if ok {
			stellaDescriptionJsonInDB.Introduction = introduction.(string)
		}
		commissionPolicy, ok := stellaDescriptionJsonInput["commission_policy"]
		if ok {
			stellaDescriptionJsonInDB.CommissionPolicy = commissionPolicy.(string)
		}
		otherNotice, ok := stellaDescriptionJsonInput["other_notice"]
		if ok {
			stellaDescriptionJsonInDB.OtherNotice = otherNotice.(string)
		}
		trafficBP, ok := stellaDescriptionJsonInput["traffic_building_policy"]
		if ok {
			stellaDescriptionJsonInDB.TrafficBuildingPolicy = trafficBP.(string)
		}
		rr, ok := stellaDescriptionJsonInput["rejected_reason"]
		if ok {
			stellaDescriptionJsonInDB.RejectedReason = rr.(string)
		}
		cookiePolicy, ok := stellaDescriptionJsonInput["cookie_policy"]
		if ok {
			stellaDescriptionJsonInDB.CookiePolicy = cookiePolicy.(string)
		}
		b, err1 := json.Marshal(stellaDescriptionJsonInDB)
		if err1 != nil {
			return err1
		}
		updates["stella_description"] = datatypes.JSON(b)
		return tx.Table("aff_campaign").Where("id = ?", id).Updates(updates).Error
	})
}

func (a *affCampaignRepository) GetAllCampaign(listStatus []string, q string, page, size int) ([]model.AffCampaignLessApp, error) {
	var listAffCampaign []model.AffCampaignLessApp
	offset := (page - 1) * size
	if q == "" {
		if err := a.Db.Table("aff_campaign").
			Joins("Brand").
			Where("aff_campaign.stella_status IN ?", listStatus).
			Limit(size + 1).
			Offset(offset).
			Find(&listAffCampaign).
			Error; err != nil {
			return listAffCampaign, err
		}
	} else {
		if err := a.Db.Table("aff_campaign").
			Joins("Brand").
			Where("aff_campaign.stella_status IN ? and "+
				"MATCH (aff_campaign.name) AGAINST(? IN NATURAL LANGUAGE MODE)", listStatus, q).
			Limit(size + 1).
			Offset(offset).
			Find(&listAffCampaign).
			Error; err != nil {
			return listAffCampaign, err
		}
	}
	return listAffCampaign, nil
}

func (a *affCampaignRepository) CountCampaign(listStatus []string, q string) (int64, error) {
	var total int64
	if q == "" {
		if err := a.Db.Table("aff_campaign").
			Where("stella_status IN ?", listStatus).
			Count(&total).Error; err != nil {
			return total, err
		}
	} else {
		if err := a.Db.Table("aff_campaign").
			Where("stella_status IN ? and MATCH (name) AGAINST(? IN NATURAL LANGUAGE MODE)", listStatus, q).
			Count(&total).Error; err != nil {
			return total, err
		}
	}
	return total, nil
}

func NewConsoleCampaignRepository(db *gorm.DB) interfaces.ConsoleCampRepository {
	return &affCampaignRepository{
		Db: db,
	}
}
