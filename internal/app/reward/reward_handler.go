package reward

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/dto"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	usecase interfaces.RewardUCase
}

func NewRewardHandler(usecase interfaces.RewardUCase) *RewardHandler {
	return &RewardHandler{
		usecase: usecase,
	}
}

// GetAllReward Get reward summary of an account
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
	res, err := handler.usecase.GetRewardSummary(ctx, user.ID)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get reward summary", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}

// GetRewardHistory Get reward withdraw history records
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
	res, err := handler.usecase.GetWithdrawHistory(ctx, user.ID, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get withdraw history", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}

// GetAllReward Claim reward of all orders
// @Summary Claim reward of all orders
// @Description Claim reward of all orders
// @Tags 	reward
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 		{object}	dto.RewardSummary
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/rewards/withdraw [post]
func (handler *RewardHandler) WithdrawReward(ctx *gin.Context) {
	// First, take user from JWT
	// user, err := dto.GetUserInfo(ctx)
	// if err != nil {
	// 	util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
	// 	return
	// }
	user := dto.UserInfo{
		ID:            208,
		WalletAddress: "0x1f6183b9b06a90c4df957bcb4cbfa8a1d72d82e3",
	}

	// get reward
	res, err := handler.usecase.WithdrawReward(ctx, user.ID, user.WalletAddress)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to withdraw reward", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}
