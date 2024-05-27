package fair_launch

import (
	"net/http"

	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/gin-gonic/gin"
)

type FairLaunchHandler struct {
	UCase interfaces.FairLauchUCase
}

func NewFairLaunchHandler(u interfaces.FairLauchUCase) *FairLaunchHandler {
	return &FairLaunchHandler{
		UCase: u,
	}
}

// GetMeme20MetaByTicker Get Meme20 metadata by ticker
// @Summary Get Meme20 metadata by ticker
// @Description Get Meme20 metadata by ticker
// @Tags 	category
// @Accept	json
// @Produce json
// @Param ticker query string false "ticker to query"
// @Success 200 		{object}	dto.Meme20MetaDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/module/fair_lauch/meme/{ticker} [get]
func (handler *FairLaunchHandler) GetMeme20MetaByTicker(ctx *gin.Context) {
	ticker := ctx.Param("ticker")
	response, err := handler.UCase.GetMeme20MetaByTicker(ctx, ticker)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get top favorited aff brands error: ", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
