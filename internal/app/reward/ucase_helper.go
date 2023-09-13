package reward

import (
	"context"
	"encoding/json"

	"github.com/astraprotocol/affiliate-system/internal/infra/msgqueue"
)

func (u *RewardUsecase) getNewApprovedAtOrderId(ctx context.Context) (string, error) {
	var msg msgqueue.MsgOrderApproved
	m, err := u.approveQ.FetchMessage(ctx)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(m.Value, &msg)
	if err != nil {
		return "", err
	}

	return msg.AtOrderID, nil
}
