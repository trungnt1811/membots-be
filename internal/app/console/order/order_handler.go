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
		util.RespondError(ctx, http.StatusBadGateway, "parse query error", err)
		return
	}

	resp, err := handler.usecase.GetOrderList(&q)
	if err != nil {
		util.RespondError(ctx, http.StatusBadGateway, "get list error", err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
