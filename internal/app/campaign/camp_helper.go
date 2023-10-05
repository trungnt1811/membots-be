package campaign

import (
	"github.com/astraprotocol/affiliate-system/internal/infra/accesstrade/types"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

func CheckCampaignChanged(old *model.AffCampaign, synced *types.ATCampaign) (map[string]any, map[string]any) {
	campChanges := map[string]any{}
	descriptionChanges := map[string]any{}

	if old.Status != synced.Status {
		// Campaign status changed
		campChanges["status"] = synced.Status
		if synced.Status != 1 {
			// Campaign not running anymore
			campChanges["stella_status"] = model.StellaStatusEnded
		}
	}

	if old.MaxCom != synced.MaxCom {
		campChanges["max_com"] = synced.MaxCom
	}
	if old.Logo != synced.Logo {
		campChanges["logo"] = synced.Logo
	}
	if old.Name != synced.Name {
		campChanges["name"] = synced.Name
	}
	if old.Scope != synced.Scope {
		campChanges["scope"] = synced.Scope
	}
	if old.Type != synced.Type {
		campChanges["scope"] = synced.Type
	}
	if old.Url != synced.Url {
		campChanges["url"] = synced.Url
	}
	if old.CookieDuration != synced.CookieDuration {
		campChanges["cookie_duration"] = synced.CookieDuration
	}
	if old.CookiePolicy != synced.CookiePolicy {
		campChanges["cookie_policy"] = synced.CookiePolicy
	}

	if old.StartTime != nil && old.StartTime.Compare(synced.StartTime.Time) > 1000 {
		// When start time is changed
		campChanges["start_time"] = synced.StartTime.Time
	}
	if !synced.EndTime.IsZero() && old.EndTime == nil {
		// When end time is set
		campChanges["end_time"] = synced.EndTime.Time
	}

	if old.Description.ActionPoint != synced.Description.ActionPoint {
		descriptionChanges["action_point"] = synced.Description.ActionPoint
	}
	if old.Description.CommissionPolicy != synced.Description.CommissionPolicy {
		descriptionChanges["commission_policy"] = synced.Description.CommissionPolicy
	}
	if old.Description.CookiePolicy != synced.Description.CookiePolicy {
		descriptionChanges["cookie_policy"] = synced.Description.CookiePolicy
	}
	if old.Description.Introduction != synced.Description.Introduction {
		descriptionChanges["introduction"] = synced.Description.Introduction
	}
	if old.Description.OtherNotice != synced.Description.OtherNotice {
		descriptionChanges["other_notice"] = synced.Description.OtherNotice
	}
	if old.Description.RejectedReason != synced.Description.RejectedReason {
		descriptionChanges["rejected_reason"] = synced.Description.RejectedReason
	}
	if old.Description.TrafficBuildingPolicy != synced.Description.TrafficBuildingPolicy {
		descriptionChanges["traffic_building_policy"] = synced.Description.TrafficBuildingPolicy
	}

	return campChanges, descriptionChanges
}
