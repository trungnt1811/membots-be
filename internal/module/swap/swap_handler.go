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

// GetSwapHistoryByAddress Get swaps history by memeAddress
// @Summary Get swaps history by memeAddress
// @Description Get swaps history by memeAddress
// @Tags 	swap
// @Accept	json
// @Produce json
// @Param memeAddress query string false "memeAddress to query, default is "
// @Success 200 		{object}	dto.SwapHistoryByAddressResp
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/swaps [get]
func (handler *SwapHandler) GetSwapHistoryByAddress(ctx *gin.Context) {
	memeAddress := ctx.DefaultQuery("memeAddress", "")

	if memeAddress == "" {
		util.RespondError(ctx, http.StatusInternalServerError, "Get swaps by memeAddress error: ", fmt.Errorf("memeAddress is empty"))
		return
	}

	response, err := handler.UCase.GetSwaps(ctx, memeAddress)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get swaps by memeAddress error: ", err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetSwapRouter Get Swap router
// @Summary Swap router
// @Description Swap router
// @Tags 	swap
// @Accept	json
// @Produce json
// @Param protocols query string false "protocols to query, default is v3"
// @Param tokenInAddress query string false "tokenInAddress to query"
// @Param tokenInChainId query string false "tokenInChainId to query"
// @Param tokenOutAddress query string false "tokenOutAddress to query"
// @Param tokenOutChainId query string false "tokenOutChainId to query"
// @Param amount query string false "amount to query"
// @Param type query string false "type to query"
// @Success 200 		{object}	interface{}
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/quote [get]
func (handler *SwapHandler) GetSwapRouter(ctx *gin.Context) {
	// protocols := ctx.DefaultQuery("protocols", "v3")
	// TODO: maybe use another version in the future
	protocols := "v3"
	tokenInAddress := ctx.DefaultQuery("tokenInAddress", "")
	tokenInChainId := ctx.DefaultQuery("tokenInChainId", "")
	tokenOutAddress := ctx.DefaultQuery("tokenOutAddress", "")
	tokenOutChainId := ctx.DefaultQuery("tokenOutChainId", "")
	amount := ctx.DefaultQuery("amount", "")
	typeQuery := ctx.DefaultQuery("type", "")
	url := fmt.Sprintf("https://www.trugly.meme/api/v1/quote?protocols=%s&tokenInAddress=%s&tokenInChainId=%s&tokenOutAddress=%s&tokenOutChainId=%s&amount=%s&type=%s",
		protocols, tokenInAddress, tokenInChainId, tokenOutAddress, tokenOutChainId, amount, typeQuery)

	response, err := handler.UCase.GetQuote(ctx, url)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get quote error: ", err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
