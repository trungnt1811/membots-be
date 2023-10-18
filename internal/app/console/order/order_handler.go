package consoleOrder

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type ConsoleOrderHandler struct {
	usecase interfaces.ConsoleOrderUcase
}

func NewConsoleOrderHandler(usecase interfaces.ConsoleOrderUcase) *ConsoleOrderHandler {
	return &ConsoleOrderHandler{
		usecase: usecase,
	}
}

// GetOrderList Get affiliate order list
// @Summary Get affiliate order list
// @Description Get affiliate order list by time range and other filter
// @Tags 	console
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Param 	payload	query 		dto.OrderListQuery false "Order list query"
// @Success 200 		{object}	dto.OrderListResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/console/orders [get]
func (handler *ConsoleOrderHandler) GetOrderList(ctx *gin.Context) {
	var q dto.OrderListQuery
	err := ctx.BindQuery(&q)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "parse query error", err)
		return
	}

	resp, err := handler.usecase.GetOrderList(&q)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "get list error", err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetOrderByOrderId Get affiliate order by order id
// @Summary Get affiliate order by order id
// @Description Get affiliate order by order id
// @Tags 	console
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Param 	orderId	path 	string true "Order id param"
// @Success 200 		{object}	dto.AffOrder
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/console/orders/:orderId [get]
func (handler *ConsoleOrderHandler) GetOrderByOrderId(ctx *gin.Context) {
	orderId, ok := ctx.Params.Get("orderId")
	if !ok {
		util.RespondError(ctx, http.StatusBadRequest, "order_id param is required")
		return
	}

	resp, err := handler.usecase.GetOrderByOrderId(orderId)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "get order error", err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// SyncOrderReward Recalculate order reward
// @Summary Recalculate order reward
// @Description Recalculate order reward (only if order is not yet rewarding, meaning that order status must be one of 'initial','pending','approved')
// @Tags 	console
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Param 	payload	body 			dto.SyncOrderRewardPayload true "accesstrade order id to sync, required"
// @Success 200 		{object}	dto.ResponseMessage
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/console/orders/sync-reward [post]
func (handler *ConsoleOrderHandler) SyncOrderReward(ctx *gin.Context) {
	var payload dto.SyncOrderRewardPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "send payload is required", err)
		return
	}

	err := handler.usecase.SyncOrderReward(payload.AccessTradeOrderId)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "get order error", err)
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseMessage{
		Message: "success",
	})
}
