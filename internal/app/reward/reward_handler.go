package reward

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/astraprotocol/affiliate-system/internal/dto"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RewardHandler struct {
	usecase interfaces.RewardUCase
}

func NewRewardHandler(usecase interfaces.RewardUCase) *RewardHandler {
	return &RewardHandler{
		usecase: usecase,
	}
}

// GetRewardByOrderId Get reward by affiliate order
// @Summary PostRequestRedeem Get reward by affiliate order
// @Description Get reward by affiliate order
// @Tags 	reward
// @Accept	json
// @Produce json
// @Param affOrderId query number false "affiliate order id to query"
// @Success 200 		{object}	dto.RewardDto
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/rewards/by-order-id [get]
func (handler *RewardHandler) GetRewardByOrderId(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	// Then verify payload data
	orderIdQuery := ctx.DefaultQuery("affOrderId", "")
	orderId, err := strconv.Atoi(orderIdQuery)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "invalid affOrderId", err)
		return
	}

	// get reward
	res, err := handler.usecase.GetRewardByOrderId(ctx, user.ID, uint(orderId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.RespondError(ctx, http.StatusNotFound, "failed to get reward", errors.New("reward not found"))
			return
		}
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get reward", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}

// GetAllReward Get reward summary of an account
// @Summary Get reward summary of an account
// @Description Get reward summary of an account
// @Tags 	reward
// @Accept	json
// @Produce json
// @Success 200 		{object}	dto.RewardSummary
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/rewards/summary [get]
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

// GetAllReward Get all rewards
// @Summary Get all rewards
// @Description Get all rewards
// @Tags 	reward
// @Accept	json
// @Produce json
// @Success 200 		{object}	dto.RewardResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/rewards [get]
func (handler *RewardHandler) GetAllReward(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	// get reward
	res, err := handler.usecase.GetAllReward(ctx, user.ID, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get rewards", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}

// GetRewardHistory Get reward history records
// @Summary Get reward history records
// @Description Get reward history records
// @Tags 	reward
// @Accept	json
// @Produce json
// @Success 200 		{object}	dto.RewardHistoryResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/rewards/claims [get]
func (handler *RewardHandler) GetClaimHistory(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	// get reward
	res, err := handler.usecase.GetClaimHistory(ctx, user.ID, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get reward history", err)
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
// @Success 200 		{object}	dto.RewardSummary
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/rewards/claims [post]
func (handler *RewardHandler) ClaimReward(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	// get reward
	res, err := handler.usecase.ClaimReward(ctx, user.ID, user.WalletAddress)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to claim reward", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}
