package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/infra/exchange"
	"github.com/stretchr/testify/assert"
)

func TestGetASAPrice(t *testing.T) {
	config := conf.GetConfiguration().Tiki
	priceRepo := exchange.NewTikiClient(exchange.TikiClientConfig{BaseUrl: config.ApiUrl})

	currentPrice, err := priceRepo.GetAstraPriceFromExchange(context.Background(), exchange.CURRENT_PRICE)
	assert.Nil(t, err)
	fmt.Println("currentPrice", currentPrice)

	sendPrice, err := priceRepo.GetAstraPriceFromExchange(context.Background(), exchange.FOR_SEND_TO_USER)
	assert.Nil(t, err)
	fmt.Println("sendPrice", sendPrice)

	receivePrice, err := priceRepo.GetAstraPriceFromExchange(context.Background(), exchange.FOR_RECEIVE_FROM_USER)
	assert.Nil(t, err)
	fmt.Println("receivePrice", receivePrice)
}
