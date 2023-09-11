package order

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
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
		})
		return
	}

	m, err := handler.usecase.PostBackUpdateOrder(&req)
	if err != nil {
		// Log error and return failed
		log.LG.Errorf("save post back order error: %v", err)
		c.JSON(http.StatusFailedDependency, gin.H{
			"message": "save post back order error",
		})
		return
	}
	c.JSON(http.StatusOK, &dto.ATPostBackResponse{
		Success: true,
		OrderId: m.AccessTradeOrderId,
	})
}
