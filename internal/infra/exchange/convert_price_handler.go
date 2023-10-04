package exchange

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/model"
)

type convertPriceHandler struct {
	tokenPriceRepo interfaces.TokenPriceRepo
}

func (c *convertPriceHandler) ConvertVndPriceToAstra(ctx context.Context, attributes []model.AffCampaignAttribute) string {
	astraPrice, err := c.tokenPriceRepo.GetAstraPrice(ctx)
	if err != nil {
		return ""
	}
	if len(attributes) > 0 {
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
		if attributes[0].AttributeType == "percent" {
			return fmt.Sprint(attributes[0].AttributeValue, "%")
		} else {
			s, err := strconv.ParseFloat(attributes[0].AttributeValue, 64)
			if err != nil {
				return ""
			}
			tmp := s / float64(astraPrice)
			return fmt.Sprintf("%.2f ASA", tmp)
		}
	} else {
		return ""
	}
}

func NewConvertPriceHandler(tokenPriceRepo interfaces.TokenPriceRepo) interfaces.ConvertPriceHandler {
	return &convertPriceHandler{tokenPriceRepo: tokenPriceRepo}
}
