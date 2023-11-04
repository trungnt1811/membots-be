package accesstrade

import (
	"fmt"
	"strings"

	"github.com/astraprotocol/affiliate-system/internal/app/campaign"
	"github.com/astraprotocol/affiliate-system/internal/infra/discord"
	interfaces2 "github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
)

type accessTradeUCase struct {
	Repo           interfaces2.ATRepository
	CampaignRepo   interfaces2.CampaignRepository
	DiscordWebhook string
	DiscordSender  discord.DiscordSender
}

func NewAccessTradeUCase(repo interfaces2.ATRepository, campRepo interfaces2.CampaignRepository, discordWebhook string) interfaces2.ATUCase {
	return &accessTradeUCase{
		Repo:         repo,
		CampaignRepo: campRepo,
		DiscordSender: discord.DiscordSender{
			DiscordWebhook: discordWebhook,
		},
	}
}

func (u *accessTradeUCase) QueryAndSaveCampaigns(onlyApproval bool) (int, int, error) {
	// Then, query campaigns
	page := 1
	limit := 20
	totalPages := 1
	totalSynced := 0
	totalUpdated := 0

	// First query all active campaign ids
	ids, err := u.CampaignRepo.QueryActiveIds()
	if err != nil {
		log.LG.Errorf("query active campaigns error: %v", err)
		return 0, 0, err
	}
	checkedCampaignIds := map[uint]bool{}
	for _, id := range ids {
		checkedCampaignIds[id] = false
	}

	for true {
		_, atResp, err := u.Repo.QueryCampaigns(true, page, limit)
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
			// Take the saved campaign for update
			oldCampaign, ok := savedCampaigns[atApproved.Id]
			if !ok {
				// Not yet save, insert this new campaign
				// and find merchant
				created, err := u.CampaignRepo.SaveATCampaign(&atApproved)
				if err != nil {
					log.LG.Errorf("create campaign error: %v", err)
					continue
				}
				totalSynced += 1
				err = u.sendMsgNewCampaignSync(created)
				if err != nil {
					log.LG.Errorf("send msg discord error: %v", err)
				}
			} else {
				// Mark this campaign as checked
				checkedCampaignIds[oldCampaign.ID] = true

				// Update campaign if changed
				campChanges, descriptionChanges := campaign.CheckCampaignChanged(&oldCampaign, &atApproved)
				err := u.CampaignRepo.UpdateCampaignByID(oldCampaign.ID, campChanges, descriptionChanges)
				if err != nil {
					log.LG.Errorf("update campaign '%d' error: %v", oldCampaign.ID, err)
					continue
				}
				if len(campChanges) != 0 || len(descriptionChanges) != 0 {
					totalUpdated += 1

					err = u.sendMsgCampaignUpdated(&oldCampaign, campChanges, descriptionChanges)
					if err != nil {
						log.LG.Errorf("send msg discord error: %v", err)
					}
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
	// Mark remain unchecked campaigns as ended
	for id, checked := range checkedCampaignIds {
		if !checked {
			// Mark campaign as ended
			campaign, err := u.CampaignRepo.DeactivateCampaign(id)
			if err != nil {
				log.LG.Errorf("deactivate campaign '%d' error: %v", id, err)
				continue
			}

			// Update aff_link active status when campaign ended
			err = u.CampaignRepo.DeactivateCampaignLinks(id)
			if err != nil {
				log.LG.Errorf("deactivate campaign '%d' links error: %v", id, err)
				continue
			}
			_ = u.sendMsgCampaignUpdated(campaign, map[string]any{"stella_status": model.StellaStatusEnded}, map[string]any{})
		}
	}

	return totalSynced, totalUpdated, nil
}

func (u *accessTradeUCase) sendMsgNewCampaignSync(camp *model.AffCampaign) error {
	msg := fmt.Sprintf("Action: New Campaign Synced\nCampaign ID: %d\nName: %s", camp.ID, camp.Name)
	return u.sendDiscordMsg(msg)
}

func (u *accessTradeUCase) sendMsgCampaignUpdated(camp *model.AffCampaign, changes map[string]any, description map[string]any) error {
	fields := make([]string, 0, len(changes)+len(description))
	for k := range changes {
		if k != "updated_at" {
			fields = append(fields, k)
		}
	}
	for k := range description {
		if k != "updated_at" {
			fields = append(fields, k)
		}
	}

	msg := fmt.Sprintf("Action: Campaign Updated\nCampaign ID: %d\nName: %s\nFields Changed: %s", camp.ID, camp.Name, strings.Join(fields, ","))
	return u.sendDiscordMsg(msg)
}

func (u *accessTradeUCase) sendDiscordMsg(msg string) error {
	return u.DiscordSender.SendMsg("CAMPAIGN SYNC", msg)
}
