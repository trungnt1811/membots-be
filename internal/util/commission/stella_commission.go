package commission

import (
	"fmt"
	"sort"
	"strings"

	"github.com/astraprotocol/affiliate-system/internal/model"
)

func GetStellaMaxCom(attributes []model.AffCampaignAttribute) string {
	if len(attributes) > 0 {
		sort.Slice(attributes, func(i, j int) bool {
			if attributes[i].AttributeType != attributes[j].AttributeType {
				return model.AttributeTypePriorityMapping[attributes[i].AttributeType] < model.AttributeTypePriorityMapping[attributes[j].AttributeType]
			}
			switch strings.Compare(attributes[i].AttributeValue, attributes[j].AttributeValue) {
			case 1:
				return true
			default:
				return false
			}
		})
		if attributes[0].AttributeType == "percent" {
			return fmt.Sprint(attributes[0].AttributeValue, "%")
		} else {
			return fmt.Sprint(attributes[0].AttributeValue, " VND")
		}
	} else {
		return ""
	}
}
