package order

import (
	"net/http"
	"strconv"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase interfaces.OrderUcase
}

func NewOrderHandler(usecase interfaces.OrderUcase) *OrderHandler {
	return &OrderHandler{
		usecase: usecase,
	}
}

// PostBackOrderHandle A callback to receive order from AccessTrade
// @Summary PostBackOrderHandle
// @Description A callback to receive order from AccessTrade
// @Tags 		redeem
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Param 	payload	body 			dto.ATPostBackRequest true "Request create link payload, required"
// @Success 200 		{object}	dto.ATPostBackResponse "when success, return the modified order for this post back"
// @Failure 424 		{object}	util.GeneralError
// @Failure 417 		{object}	util.GeneralError
// @Router 	/api/v1/order/post-back [post]
func (handler *OrderHandler) PostBackOrderHandle(c *gin.Context) {
	var req dto.ATPostBackRequest
	err := c.BindJSON(&req)
	if err != nil {
		// Log error and return failed
		log.LG.Errorf("parse access trade post back error: %v", err)
		c.JSON(http.StatusExpectationFailed, gin.H{
			"message": "parse access trade post back error",
			"error":   err.Error(),
		})
		return
	}

	m, err := handler.usecase.PostBackUpdateOrder(&req)
	if err != nil {
		// Log error and return failed
		log.LG.Errorf("save post back order error: %v", err)
		c.JSON(http.StatusFailedDependency, gin.H{
			"message": "save post back order error",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &dto.ATPostBackResponse{
		Success: true,
		OrderId: m.AccessTradeOrderId,
	})
}

// GetRewardHistory Get affiliate order details
// @Summary Get affiliate order details
// @Description Get affiliate order details - include status timeline and reward info
// @Tags 	reward
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 		{object}	dto.OrderDetailsDto
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/order/{id} [get]
func (handler *OrderHandler) GetOrderDetails(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	orderId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "id is required", err)
		return
	}

	// get reward
	res, err := handler.usecase.GetOrderDetails(ctx, user.ID, uint(orderId))
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get order details", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}

// GetRewardHistory Get order history
// @Summary Get order history
// @Description Get order history - include reward info
// @Tags 	reward
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 		{object}	dto.OrderHistoryResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/order [get]
func (handler *OrderHandler) GetOrderHistory(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	// get reward
	res, err := handler.usecase.GetOrderHistory(ctx, user.ID, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get order history", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}
