package accesstrade

import (
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util/log"
)

type AccessTradeUsecase struct {
	Repo         interfaces.ATRepository
	CampaignRepo interfaces.CampaignRepository
}

func NewAccessTradeUsecase(repo interfaces.ATRepository, campRepo interfaces.CampaignRepository) *AccessTradeUsecase {
	return &AccessTradeUsecase{
		Repo:         repo,
		CampaignRepo: campRepo,
	}
}

func (u *AccessTradeUsecase) QueryAndSaveCampaigns(onlyApproval bool) (int, error) {
	// Then, query campaigns
	page := 1
	limit := 20
	totalPages := 1
	totalSync := 0

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
			if _, ok := savedCampaigns[atApproved.Id]; !ok {
				// Not yet save, insert this new campaign
				// and find merchant
				err := u.CampaignRepo.SaveATCampaign(&atApproved)
				if err != nil {
					log.LG.Errorf("create campaign error: %v", err)
					continue
				}
				totalSync += 1
			} else {
				// TODO: Update if campaign is changed
			}
		}

		if page >= totalPages {
			log.LG.Infof("done sync %d pages", totalPages)
			break
		}
		page += 1
	}

	return totalSync, nil
}

func (u *AccessTradeUsecase) CreateAndSaveLink() (int, error) {
	// TODO: Create link for campaign if not available
	return 0, nil
}
