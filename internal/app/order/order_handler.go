package order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

type OrderHandler struct {
	usecase interfaces.OrderUCase
}

func NewOrderHandler(usecase interfaces.OrderUCase) *OrderHandler {
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

// GetOrderDetails Get affiliate order details
// @Summary Get affiliate order details
// @Description Get affiliate order details - include status timeline and reward info
// @Tags 	order
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "order id to query"
// @Success 200 		{object}	dto.OrderDetailsDto
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/orders/{id} [get]
func (handler *OrderHandler) GetOrderDetails(ctx *gin.Context) {
	// First, take user from JWT
	// user, err := dto.GetUserInfo(ctx)
	// if err != nil {
	// 	util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
	// 	return
	// }
	user := dto.UserInfo{
		ID: 214,
	}

	orderId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "id is required", err)
		return
	}

	// get reward
	order, err := handler.usecase.GetOrderDetails(ctx, user.ID, uint(orderId))
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get order details", err)
		return
	}

	if order.ID == 0 {
		util.RespondError(ctx, http.StatusForbidden, "failed to get order details", errors.New("user do not have authorization on given order"))
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, order)
}

// GetOrderHistory Get order history
// @Summary Get order history
// @Description Get order history - include reward info
// @Tags 	order
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Param status query string false "order status to filter, valid query is 'wait_for_confirming','rewarding','complete','cancelled','rejected'"
// @Param page query string false "page to query, default is 1"
// @Success 200 		{object}	dto.OrderHistoryResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/orders [get]
func (handler *OrderHandler) GetOrderHistory(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	orderStatus := ctx.DefaultQuery("status", "")
	validStatus := []string{dto.OrderStatusWaitForConfirming, dto.OrderStatusRewarding,
		dto.OrderStatusComplete, dto.OrderStatusCancelled, dto.OrderStatusRejected}
	if orderStatus != "" && !slices.Contains(validStatus, orderStatus) {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", errors.New("invalid status query"))
		return
	}
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	// get reward
	res, err := handler.usecase.GetOrderHistory(ctx, user.ID, orderStatus, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "failed to get order history", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, res)
}
