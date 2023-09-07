package order

import (
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/interfaces"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase interfaces.OrderUsecase
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
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/order/postback [post]
func (handler *OrderHandler) PostBackOrderHandle(c *gin.Context) {

}
