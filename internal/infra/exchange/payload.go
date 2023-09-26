package exchange

import "strconv"

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
