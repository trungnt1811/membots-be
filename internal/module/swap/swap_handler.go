package swap

import (
	"fmt"
	"net/http"

	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/gin-gonic/gin"

	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

type SwapHandler struct {
	UCase interfaces.SwapUCase
}

func NewSwapHandler(ucase interfaces.SwapUCase) *SwapHandler {
	return &SwapHandler{
		UCase: ucase,
	}
}

// GetSwapHistoryByAddress Get swaps history by memeId
// @Summary Get swaps history by memeId
// @Description Get swaps history by memeId
// @Tags 	swap
// @Accept	json
// @Produce json
// @Param memeId query string false "memeId to query, default is "
// @Success 200 		{object}	dto.SwapHistoryByAddressResp
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/truglymeme/swaps [get]
func (handler *SwapHandler) GetSwapHistoryByAddress(ctx *gin.Context) {
	memeId := ctx.DefaultQuery("memeId", "")

	if memeId == "" {
		util.RespondError(ctx, http.StatusInternalServerError, "Get swaps by memeId error: ", fmt.Errorf("memeId is empty"))
		return
	}

	response, err := handler.UCase.GetSwaps(ctx, memeId)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get swaps by memeId error: ", err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
