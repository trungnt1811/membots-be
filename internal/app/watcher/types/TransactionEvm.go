package types

import "time"

type UTCTime struct {
	unixNano int64
}

type Log struct {
	Address     string   `json:"address"`
	AddressName string   `json:"addressName"`
	Data        string   `json:"data"`
	Index       string   `json:"index"`
	Topics      []string `json:"topics"`
}

type TokenTransfer struct {
	Amount               string `json:"amount"`
	Decimals             string `json:"decimals"`
	FromAddress          string `json:"fromAddress"`
	FromAddressName      string `json:"fromAddressName"`
	LogIndex             string `json:"logIndex"`
	ToAddress            string `json:"toAddress"`
	ToAddressName        string `json:"toAddressName"`
	TokenContractAddress string `json:"tokenContractAddress"`
	TokenName            string `json:"tokenName"`
	TokenSymbol          string `json:"tokenSymbol"`
	TokenId              string `json:"tokenId"`
	TokenType            string `json:"tokenType"`
}

type TransactionEvm struct {
	BlockHeight                  int64     `json:"blockHeight"`
	BlockHash                    string    `json:"blockHash"`
	BlockTime                    time.Time `json:"blockTime"`
	Confirmations                int64     `json:"confirmations"`
	Hash                         string    `json:"hash"`
	CosmosHash                   string    `json:"cosmosHash"`
	Index                        int       `json:"index"`
	Success                      bool      `json:"success"`
	Error                        string    `json:"error"`
	RevertReason                 string    `json:"revertReason"`
	CreatedContractAddressHash   string    `json:"createdContractAddressHash"`
	CreatedContractAddressName   string    `json:"createdContractAddressName"`
	CreatedContractCodeIndexedAt time.Time `json:"createdContractCodeIndexedAt"`
	From                         string    `json:"from"`
	FromAddressName              string    `json:"fromAddressName"`
	To                           string    `json:"to"`
	ToAddressName                string    `json:"toAddressName"`
	IsInteractWithContract       bool      `json:"isInteractWithContract"`
	Value                        string    `json:"value"`

	// If return cosmos tx, we got
	Status string `json:"status"`
}

type TxResp struct {
	Message string         `json:"message"`
	Result  TransactionEvm `json:"result"`
	Status  string         `json:"status"`
}
