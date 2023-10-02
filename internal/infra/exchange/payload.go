package exchange

import (
	"strconv"

	"github.com/astraprotocol/affiliate-system/internal/util"
)

// Helper Types

type ErrorResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Errors  []interface{} `json:"errors"`
}

type ErrorResponseData struct {
	Error ErrorResponse `json:"error"`
}

type ExchangeErrorResponse struct {
	Errors []string `json:"errors"`
}

// DTOs

type ExchangePriceResponse struct {
	Timestamp uint64     `json:"timestamp"`
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
}

type ExchangePrice struct {
	Timestamp uint64      `json:"timestamp"`
	Asks      [][]float64 `json:"asks"`
	Bids      [][]float64 `json:"bids"`
}

func (o *ExchangePriceResponse) ToExchangePrice() ExchangePrice {
	asks := [][]float64{}
	bids := [][]float64{}

	for _, item := range o.Asks {
		if len(item) < 2 {
			continue
		}
		price, err := strconv.ParseFloat(item[0], 64)
		if err != nil {
			continue
		}
		amount, err := strconv.ParseFloat(item[1], 64)
		if err != nil {
			continue
		}
		asks = append(asks, []float64{price, amount})
	}

	for _, item := range o.Bids {
		if len(item) < 2 {
			continue
		}
		price, err := strconv.ParseFloat(item[0], 64)
		if err != nil {
			continue
		}
		amount, err := strconv.ParseFloat(item[1], 64)
		if err != nil {
			continue
		}
		bids = append(bids, []float64{price, amount})
	}

	return ExchangePrice{
		Timestamp: o.Timestamp,
		Asks:      asks,
		Bids:      bids,
	}
}

func (o *ExchangePrice) GetPrice(calculateMode int, thresholdVol float64) float64 {
	var totalXu float64 = 0
	var totalASA float64 = 0

	var pricePairs [][]float64
	if calculateMode == FOR_SEND_TO_USER {
		// We send ASA to user, need to prevent user push the price down (because the want to get more ASA)
		// When user want to push the price down, they sell ASA to exchange
		// => We calculate price base on the 'asks' price so that the price is higher that current price
		// If user actually sell ASA to manipulate exchange, the price is not going down too much
		pricePairs = o.Asks
	} else if calculateMode == FOR_RECEIVE_FROM_USER {
		// vice versa
		pricePairs = o.Bids
	} else {
		// get current price
		if len(o.Asks) < 1 || len(o.Asks[0]) < 1 {
			return 0
		}
		return o.Asks[0][0]
	}

	for _, pair := range pricePairs {
		price := pair[0]
		amount := pair[1]
		if amount > thresholdVol {
			// When current amount is larger than remain volume
			// take the remain volume as calculated amount
			totalASA += thresholdVol
			totalXu += thresholdVol * price
			thresholdVol = 0
			break
		} else {
			// When current amount is less than remain volume
			// reduce the amount from remain volume and take
			// order amount as calculated one
			totalASA += amount
			totalXu += amount * price
			thresholdVol -= amount
		}
	}

	avgSellPrice := totalXu / totalASA // price that user sell to exchange
	return util.RoundFloat(avgSellPrice, 2)
}
