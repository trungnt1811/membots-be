package reward

import (
	"errors"
	"net/http"
	"strconv"

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

// GetWithdrawHistory Get reward withdraw details
// @Summary Get reward withdraw details
// @Description Get reward withdraw details
// @Tags 	reward
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "withdraw id to query"
// @Success 200 		{object}	dto.RewardWithdrawDto
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/rewards/withdraw/{id} [get]
func (handler *RewardHandler) GetWithdrawDetails(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	withdrawId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "failed to get withdraw details", errors.New("withdraw id is required"))
		return
	}

	// get reward
	res, err := handler.uCase.GetWithdrawDetails(ctx, user.ID, uint(withdrawId))
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get withdraw details", err)
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
// @Success 200 		{object}	dto.RewardWithdrawDto
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
