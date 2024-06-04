package dto

type Transaction struct {
	AmountETH     string `json:"amountETH"`
	AmountMeme    string `json:"amountMeme"`
	WalletAddress string `json:"walletAddress"`
	TxHash        string `json:"txHash"`
	TxType        string `json:"txType"`
	MemeID        string `json:"memeId"`
	Epoch         uint64 `json:"epoch"`
}

type LaunchpadInfo struct {
	Transactions []Transaction `json:"transactions"`
	Status       string        `json:"status"`
	TargetETH    string        `json:"targetETH"`
	CollectedETH string        `json:"collectedETH"`
}

type LaunchpadInfoRsp struct {
	LaunchpadInfo LaunchpadInfo `json:"launchpadInfo"`
}
