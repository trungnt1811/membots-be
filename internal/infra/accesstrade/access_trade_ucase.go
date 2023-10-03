package accesstrade

import (
	"github.com/astraprotocol/affiliate-system/internal/app/campaign"
	interfaces2 "github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
)

type accessTradeUCase struct {
	Repo         interfaces2.ATRepository
	CampaignRepo interfaces2.CampaignRepository
}

func NewAccessTradeUCase(repo interfaces2.ATRepository, campRepo interfaces2.CampaignRepository) interfaces2.ATUCase {
	return &accessTradeUCase{
		Repo:         repo,
		CampaignRepo: campRepo,
	}
}

func (u *accessTradeUCase) QueryAndSaveCampaigns(onlyApproval bool) (int, int, error) {
	// Then, query campaigns
	page := 1
	limit := 20
	totalPages := 1
	totalSynced := 0
	totalUpdated := 0

	for true {
		atResp, err := u.Repo.QueryCampaigns(true, page, limit)
		if err != nil {
			log.LG.Errorf("query AccessTrade campaign error: %v", err)
			break
		}

		totalPages = int(atResp.TotalPage)

		ids := make([]string, len(atResp.Data))
		for i, item := range atResp.Data {
			ids[i] = item.Id
		}

		savedCampaigns, err := u.CampaignRepo.RetrieveCampaignsByAccessTradeIds(ids)
		if err != nil {
			log.LG.Errorf("retrieve at campaign by ids error: %v", err)
			break
		}

		for _, atApproved := range atResp.Data {
			oldCampaign, ok := savedCampaigns[atApproved.Id]
			if !ok {
				// Not yet save, insert this new campaign
				// and find merchant
				err := u.CampaignRepo.SaveATCampaign(&atApproved)
				if err != nil {
					log.LG.Errorf("create campaign error: %v", err)
					continue
				}
				// TODO: push kafka message for synced campaign
				totalSynced += 1
			} else {
				// Update campaign if changed
				campChanges, descriptionChanges := campaign.CheckCampaignChanged(&oldCampaign, &atApproved)
				err := u.CampaignRepo.UpdateCampaignByID(oldCampaign.ID, campChanges, descriptionChanges)
				if err != nil {
					log.LG.Errorf("update campaign '%d' error: %v", oldCampaign.ID, err)
					continue
				}
				// TODO: push kafka message for updated campaign
				if len(campChanges) != 0 || len(descriptionChanges) != 0 {
					totalUpdated += 1
				}
				// Update aff_link active status if campaign ended
				if campChanges["stella_status"] == model.StellaStatusEnded {
					err := u.CampaignRepo.DeactivateCampaignLinks(oldCampaign.ID)
					if err != nil {
						log.LG.Errorf("deactivate campaign '%d' links error: %v", oldCampaign.ID, err)
						continue
					}
				}
			}
		}

		if page >= totalPages {
			log.LG.Infof("done sync %d pages", totalPages)
			break
		}
		page += 1
	}

	return totalSynced, totalUpdated, nil
}
