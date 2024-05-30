package memeception

import (
	"fmt"
	"net/http"

	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/gin-gonic/gin"

	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

type MemeceptionHandler struct {
	UCase interfaces.MemeceptionUCase
}

func NewMemeceptionHandler(ucase interfaces.MemeceptionUCase) *MemeceptionHandler {
	return &MemeceptionHandler{
		UCase: ucase,
	}
}

// GetMemeceptionBySymbol Get memeception by symbol
// @Summary Get memeception by symbol
// @Description Get memeception by symbol
// @Tags 	memeception
// @Accept	json
// @Produce json
// @Param symbol query string false "symbol to query, default is "
// @Success 200 		{object}	dto.MemeceptionDetailResp
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/memeception [get]
func (handler *MemeceptionHandler) GetMemeceptionBySymbol(ctx *gin.Context) {
	symbol := ctx.DefaultQuery("symbol", "")

	if symbol == "" {
		util.RespondError(ctx, http.StatusInternalServerError, "Get memeception by symbol error: ", fmt.Errorf("symbol is empty"))
		return
	}

	response, err := handler.UCase.GetMemeceptionBySymbol(ctx, symbol)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get memeception by symbol error: ", err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetMemeceptions Get memeceptions
// @Summary Get memeceptions
// @Description Get memeceptions
// @Tags 	memeception
// @Accept	json
// @Produce json
// @Success 200 		{object}	dto.MemeceptionsResp
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/memeceptions [get]
func (handler *MemeceptionHandler) GetMemeceptions(ctx *gin.Context) {
	response, err := handler.UCase.GetMemeceptions(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get memeceptions error: ", err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
