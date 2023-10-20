package alert

import (
	"fmt"
)

type SuspiciousOrderMsg struct {
	OrderId   string  `json:"order_id"`
	Billing   uint    `json:"billing"`
	Reward    float64 `json:"reward"`
	Threshold float64 `json:"threshold"`
}

func (msg SuspiciousOrderMsg) String() string {
	msgFormatter := new(MsgFormatter).
		FormatTitle("Aff - Big Reward Order").
		FormatKeyValueMsg("Billing", fmt.Sprintf("%v VND", msg.Billing)).
		FormatKeyValueMsg("Reward", fmt.Sprintf("%v ASA", msg.Reward)).
		FormatKeyValueMsg("Threshold", fmt.Sprintf("%v ASA", msg.Threshold))

	return msgFormatter.String()
}
