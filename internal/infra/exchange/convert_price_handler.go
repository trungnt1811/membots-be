package exchange

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type convertPriceHandler struct {
	tokenPriceRepo interfaces.TokenPriceRepo
}

func (c *convertPriceHandler) GetStellaMaxCommission(ctx context.Context, attributes []model.AffCampaignAttribute) string {
	if len(attributes) == 0 {
		return ""
	}

	sort.Slice(attributes, func(i, j int) bool {
		if attributes[i].AttributeType != attributes[j].AttributeType {
			return model.AttributeTypePriorityMapping[attributes[i].AttributeType] < model.AttributeTypePriorityMapping[attributes[j].AttributeType]
		}
		value1, _ := strconv.ParseFloat(strings.TrimSpace(attributes[i].AttributeValue), 64)
		value2, _ := strconv.ParseFloat(strings.TrimSpace(attributes[j].AttributeValue), 64)
		if value1 > value2 {
			return true
		} else {
			return false
		}
	})

	maxCommission := c.ConvertVndPriceToAstra(ctx, attributes[0])
	if attributes[0].AttributeType == "percent" {
		maxCommission += "%"
	}
	return maxCommission
}

func (c *convertPriceHandler) ConvertVndPriceToAstra(ctx context.Context, attribute model.AffCampaignAttribute) string {
	astraPrice, err := c.tokenPriceRepo.GetAstraPrice(ctx)
	if err != nil {
		return ""
	}

	stellaCommission := conf.GetConfiguration().Aff.StellaCommission
	value, err := strconv.ParseFloat(attribute.AttributeValue, 64)
	if err != nil {
		return ""
	}
	netValue := value - value*stellaCommission/100

	if attribute.AttributeType == "percent" {
		// for percentage attribute, remove unit
		return fmt.Sprintf("%.2f", netValue)
	} else {
		// for vnd attribute, we convert to asa
		asaValue := netValue / float64(astraPrice)
		return fmt.Sprintf("%.2f ASA", asaValue)
	}

}

func NewConvertPriceHandler(tokenPriceRepo interfaces.TokenPriceRepo) interfaces.ConvertPriceHandler {
	return &convertPriceHandler{tokenPriceRepo: tokenPriceRepo}
}
