package reward

import (
	"errors"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"net/http"
	"strconv"

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
	res, err := handler.usecase.GetRewardByOrderId(ctx, user.UserInfo.ID, uint(orderId))
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
	res, err := handler.usecase.GetAllReward(ctx, user.UserInfo.ID, page, size)
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
// @Router 	/api/v1/rewards/history [get]
func (handler *RewardHandler) GetRewardHistory(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	// get reward
	res, err := handler.usecase.GetRewardHistory(ctx, user.UserInfo.ID, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get reward history", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}
