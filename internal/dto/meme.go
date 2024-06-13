package dto

type MemeDetail struct {
	ID              uint64            `json:"id"`
	Name            string            `json:"name"`
	Symbol          string            `json:"symbol"`
	Description     string            `json:"description"`
	TotalSupply     string            `json:"totalSupply"`
	Decimals        uint64            `json:"decimals"`
	LogoUrl         string            `json:"logoUrl"`
	BannerUrl       string            `json:"bannerUrl"`
	CreatorAddress  string            `json:"creatorAddress"`
	ContractAddress string            `json:"contractAddress"`
	SwapFeeBps      uint64            `json:"swapFeeBps"`
	VestingAllocBps uint64            `json:"vestingAllocBps"`
	Meta            string            `json:"meta"`
	Live            bool              `json:"live"`
	NetworkID       uint64            `json:"networkId"`
	Website         string            `json:"website"`
	Memeception     MemeceptionCommon `json:"memeception"`
	Socials         map[string]Social `json:"socials"`
}

type MemeCommon struct {
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	LogoUrl         string `json:"logoUrl"`
	BannerUrl       string `json:"bannerUrl"`
	Description     string `json:"description"`
	ContractAddress string `json:"contract_address"`
	Meta            string `json:"meta"`
}
