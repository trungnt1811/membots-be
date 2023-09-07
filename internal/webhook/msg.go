package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/astraprotocol/affiliate-system/internal/util/log"

	"github.com/astraprotocol/affiliate-system/conf"
)

type TxFailMsg struct {
	TxHash   string `json:"TxHash"`
	TxAction string `json:"TxAction"`
}

func (msg TxFailMsg) String() string {
	conf := conf.GetConfiguration()
	msgFormatter := new(MsgFormatter).
		FormatTitle("Affiliate System - Tx Fail").
		FormatKeyValueMsg("Action", msg.TxAction).
		FormatKeyValueMsg("URL", fmt.Sprintf("%s/tx/%s", conf.EvmRpc.ExplorerURL, msg.TxHash))

	return msgFormatter.String()
}

type LowBalanceMsg struct {
	AddrType   string  `json:"addr_type"`
	Address    string  `json:"address"`
	Balance    float64 `json:"balance"`
	MinRequire float64 `json:"min_require"`
}

func (msg LowBalanceMsg) String() string {

	title := fmt.Sprintf("%v - Low Balance", msg.AddrType)

	msgFormatter := new(MsgFormatter).
		FormatTitle(title).
		FormatKeyValueMsg("Worker", msg.Address).
		FormatKeyValueMsg("Balance", fmt.Sprintf("%v ASA", msg.Balance)).
		FormatKeyValueMsg("MinRequire", fmt.Sprintf("%v ASA", msg.MinRequire))

	return msgFormatter.String()
}

type LowFundBalanceMsg struct {
	SellerID    uint    `json:"id"`
	FundAddress string  `json:"fund_address"`
	Balance     float64 `json:"balance"`
	MinRequire  float64 `json:"min_require"`
}

func (msg LowFundBalanceMsg) String() string {
	msgFormatter := new(MsgFormatter).
		FormatTitle("Funding account - Low Balance").
		FormatKeyValueMsg("SellerID", msg.SellerID).
		FormatKeyValueMsg("FundAddress", msg.FundAddress).
		FormatKeyValueMsg("Balance", fmt.Sprintf("%v ASA", msg.Balance)).
		FormatKeyValueMsg("MinRequire", fmt.Sprintf("%v ASA", msg.MinRequire))

	return msgFormatter.String()
}

type EstimateGasFailMsg struct {
	Contract string      `json:"contract"`
	Action   string      `json:"action"`
	Params   interface{} `json:"params"`
}

func (msg EstimateGasFailMsg) String() string {
	params, err := json.Marshal(msg.Params)
	if err != nil {
		log.LG.Errorf("EstimateGasFailMsg - Failed to stringify tx params: %v", params)
	}

	msgFormatter := new(MsgFormatter).
		FormatTitle("Affiliate System - Estimate Gas Fail").
		FormatKeyValueMsg("Contract", msg.Contract).
		FormatKeyValueMsg("Action", msg.Action).
		FormatKeyValueMsg("Parameters", string(params))

	return msgFormatter.String()
}
