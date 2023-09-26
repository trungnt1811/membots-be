package reward

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/dto"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	uCase interfaces.RewardUCase
}

func NewRewardHandler(uCase interfaces.RewardUCase) *RewardHandler {
	return &RewardHandler{
		uCase: uCase,
	}
}

// GetRewardSummary Get reward summary of an account
// @Summary Get reward summary of an account
// @Description Get reward summary of an account
// @Tags 	reward
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 		{object}	dto.RewardSummary
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/rewards/summary [get]
func (handler *RewardHandler) GetRewardSummary(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	// get reward
	res, err := handler.uCase.GetRewardSummary(ctx, user.ID)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get reward summary", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}

// GetWithdrawHistory Get reward withdraw history records
// @Summary Get reward withdraw history records
// @Description Get reward withdraw history records
// @Tags 	reward
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 		{object}	dto.RewardWithdrawResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/rewards/withdraw [get]
func (handler *RewardHandler) GetWithdrawHistory(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	// get reward
	res, err := handler.uCase.GetWithdrawHistory(ctx, user.ID, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get withdraw history", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}

// WithdrawReward Claim reward of all orders
// @Summary Claim reward of all orders
// @Description Claim reward of all orders
// @Tags 	reward
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 		{object}	dto.RewardSummary
// @Failure 424 		{object}	util.GeneralError
// @Failure 429 		{object}	string "reach withdraw limit, max one time each three seconds"
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/rewards/withdraw [post]
func (handler *RewardHandler) WithdrawReward(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	// get reward
	res, err := handler.uCase.WithdrawReward(ctx, user.ID, user.WalletAddress)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to withdraw reward", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}
