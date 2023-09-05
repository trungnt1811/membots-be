package types

type RedeemRewardResponse struct {
	WalletAddress string `json:"walletAddress"`
	HolderAddress string `json:"holderAddress"`
	Signature     string `json:"signature"`
	Deadline      int64  `json:"deadline"`
}

type CheckRedeemResponse struct {
	Valid          bool   `json:"valid"`
	Message        string `json:"message"`
	Amount         string `json:"amount"`
	TokenAddress   string `json:"tokenAddress"`
	SellerName     string `json:"sellerName"`
	SellerLogo     string `json:"sellerLogo"`
	ProgramAddress string `json:"programAddress"`
	HolderAddress  string `json:"holderAddress"`
}
