package dto

type Swap struct {
	Buy           bool   `json:"buy"`
	Amount0       string `json:"amount0"`
	Amount1       string `json:"amount1"`
	AmountUSD     string `json:"amountUSD"`
	PriceETH      string `json:"priceETH"`
	PriceUSD      string `json:"priceUSD"`
	WalletAddress string `json:"walletAddress"`
	TxHash        string `json:"txHash"`
	SwapAt        uint64 `json:"swapAt"`
	Token1IsMeme  bool   `json:"token1IsMeme"`
}
