package dto

type Nft struct {
	ID              string `json:"id"`
	MemeID          string `json:"memeId"`
	BaseUrl         string `json:"baseUrl"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	ContractAddress string `json:"contractAddress"`
	IsFungible      bool   `json:"isFungible"`
	Tiers           []Tier `json:"tiers"`
}

type Tier struct {
	MemeId          string   `json:"memeId"`
	NftId           string   `json:"nftId"`
	Rank            uint64   `json:"rank"`
	AmountThreshold string   `json:"amountThreshold"`
	TokenIdLower    uint64   `json:"tokenIdLower"`
	TokenIdUpper    string   `json:"tokenIdUpper"`
	TokenIdURLs     []string `json:"tokenIdURLs"`
}
