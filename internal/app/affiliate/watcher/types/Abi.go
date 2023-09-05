package types

type TxAbiResp struct {
	Message string    `json:"message"`
	Result  AbiResult `json:"result"`
	Status  string    `json:"status"`
}

type AbiResult struct {
	Abi      AbiItem `json:"abi"`
	Verified bool    `json:"verified"`
}

type AbiItem struct {
	Name string `json: "name`
}
