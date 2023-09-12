package redeem

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/app/redeem/types"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"

	"github.com/gin-gonic/gin"
)

type RedeemHandler struct {
	usecase interfaces.RedeemUCase
}

func NewRedeemHandler(usecase interfaces.RedeemUCase) *RedeemHandler {
	return &RedeemHandler{
		usecase: usecase,
	}
}

// PostRequestRedeem Send cashback to customer wallet
// @Summary PostRequestRedeem
// @Description Send cashback to customer wallet by redeem code
// @Tags 		redeem
// @Accept	json
// @Produce json
// @Param 	payload	body 			types.RedeemRequestPayload true "Redeem payload, required"
// @Success 200 		{object}	types.RedeemRewardResponse "when redeem code is available, only valid if not claimed yet"
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router 	/api/v1/redeem/request [post]
func (handler *RedeemHandler) PostRequestRedeem(c *gin.Context) {
	// First, take user from JWT
	_, err := dto.GetUserInfo(c)
	if err != nil {
		util.RespondError(c, http.StatusBadRequest, "logged in user required", err)
		return
	}
	// Then verify payload data
	var payload types.RedeemRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		util.RespondError(c, http.StatusBadRequest, "request payload required", err)
		return
	}
	if err := payload.Valid(); err != nil {
		util.RespondError(c, http.StatusBadRequest, "wrong payload format", err)
		return
	}
	// send reward with redeem code
	resp, err := handler.usecase.RedeemCashback(payload)
	if err != nil {
		util.RespondError(c, http.StatusFailedDependency, "failed to redeem", err)
		return
	}

	// Response transaction status
	c.JSON(http.StatusOK, resp)
}
