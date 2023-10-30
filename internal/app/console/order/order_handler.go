package consoleOrder

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type ConsoleOrderHandler struct {
	usecase    interfaces.ConsoleOrderUcase
	orderUCase interfaces.OrderUCase
}

func NewConsoleOrderHandler(usecase interfaces.ConsoleOrderUcase, orderUCase interfaces.OrderUCase) *ConsoleOrderHandler {
	return &ConsoleOrderHandler{
		usecase:    usecase,
		orderUCase: orderUCase,
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

// GetPostBackList Get post back order list
// @Summary Get post back order list
// @Description Get post back order list by time range and other filter
// @Tags 	console
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Param 	payload	query 		dto.PostBackListQuery false "Post Back list query"
// @Success 200 		{object}	dto.PostBackListResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/console/orders/logs [get]
func (handler *ConsoleOrderHandler) GetPostBackList(ctx *gin.Context) {
	var q dto.PostBackListQuery
	err := ctx.BindQuery(&q)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "parse query error", err)
		return
	}

	resp, err := handler.usecase.GetPostBackList(&q)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "get list error", err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// RetryErrorPostBack Retry error post back by log
// @Summary Retry error post back by log
// @Description Retry error post back order by using sent data
// @Tags 	console
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Param 	pbId	path 	string true "Post back id to retry"
// @Success 200 		{object}	dto.ATPostBackResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/console/orders/retry/:pbId [post]
func (handler *ConsoleOrderHandler) RetryErrorPostBack(ctx *gin.Context) {
	pbIdParam, ok := ctx.Params.Get("pbId")
	if !ok {
		util.RespondError(ctx, http.StatusBadRequest, "pbId param is required")
		return
	}
	pbId, err := strconv.ParseUint(pbIdParam, 10, 64)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "pbId param is not a number")
		return
	}

	data, err := handler.usecase.GetPostBackList(&dto.PostBackListQuery{
		PostBackId: uint(pbId),
		IsError:    true,
	})
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "find post back error", err)
		return
	}
	if len(data.Data) == 0 {
		util.RespondError(ctx, http.StatusBadRequest, "no error post back found", err)
		return
	}

	item := data.Data[0]
	var pb dto.ATPostBackRequest
	err = json.Unmarshal([]byte(item.Data), &pb)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "parse pb data error", err)
		return
	}

	order, err := handler.orderUCase.PostBackUpdateOrder(&pb)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "post-back retry error", err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}
