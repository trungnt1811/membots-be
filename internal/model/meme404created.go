package model

type Meme404Created struct {
	ID           string `json:"id"`
	MemeToken    string `json:"memeToken"`
	ParamsSymbol string `json:"params_symbol"`
	ParamsName   string `json:"params_name"`
	Tiers        Tier   `json:"tiers"`
}

type Tier struct {
	ID              string `json:"id"`
	NftID           string `json:"nftId"`
	LowerId         string `json:"lowerId"`
	UpperId         string `json:"upperId"`
	NftSymbol       string `json:"nftSymbol"`
	NftName         string `json:"nftName"`
	AmountThreshold string `json:"amountThreshold"`
	BaseURL         string `json:"baseURL"`
	IsFungible      bool   `json:"isFungible"`
}
