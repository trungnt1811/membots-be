package memeception

import (
	"fmt"
	"net/http"

	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/gin-gonic/gin"

	"github.com/flexstack.ai/membots-be/internal/dto"
	"github.com/flexstack.ai/membots-be/internal/interfaces"
	"github.com/flexstack.ai/membots-be/internal/util/log"
)

type MemeceptionHandler struct {
	UCase interfaces.MemeceptionUCase
}

func NewMemeceptionHandler(ucase interfaces.MemeceptionUCase) *MemeceptionHandler {
	return &MemeceptionHandler{
		UCase: ucase,
	}
}

// CreateMeme Create meme
// @Summary Create meme
// @Description Create meme
// @Tags 	memeception
// @Accept	json
// @Produce json
// @Param 	payload	body 			dto.CreateMemePayload true "Request create meme, required"
// @Success 200 		{object}	dto.CreateMemePayload "When success, return {"success": true}"
// @Failure 424 		{object}	util.GeneralError
// @Failure 417 		{object}	util.GeneralError
// @Router 	/api/v1/memes [post]
func (handler *MemeceptionHandler) CreateMeme(ctx *gin.Context) {
	var req dto.CreateMemePayload
	err := ctx.BindJSON(&req)
	if err != nil {
		log.LG.Errorf("parse create meme payload error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: add check duplicated logic later
	err = handler.UCase.CreateMeme(ctx, req)
	if err != nil {
		log.LG.Errorf("save meme error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// GetMemeDetail Get meme detail
// @Summary Get meme detail
// @Description Get meme detail
// @Tags 	memeception
// @Accept	json
// @Produce json
// @Param memeAddress query string false "memeAddress to query, default is "
// @Param symbol query string false "symbol to query, default is "
// @Success 200 		{object}	dto.MemeceptionDetailResp
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/meme [get]
func (handler *MemeceptionHandler) GetMemeDetail(ctx *gin.Context) {
	memeAddress := ctx.DefaultQuery("memeAddress", "")
	symbol := ctx.DefaultQuery("symbol", "")

	if memeAddress == "" && symbol == "" {
		util.RespondError(ctx, http.StatusInternalServerError, "Get memeception by meme address error: ", fmt.Errorf("meme address is empty"))
		return
	}
	if memeAddress != "" {
		response, err := handler.UCase.GetMemeceptionByContractAddress(ctx, memeAddress)
		if err != nil {
			util.RespondError(ctx, http.StatusInternalServerError, "Get memeception by meme address error: ", err)
			return
		}

		ctx.JSON(http.StatusOK, response)
	} else {
		response, err := handler.UCase.GetMemeceptionBySymbol(ctx, symbol)
		if err != nil {
			util.RespondError(ctx, http.StatusInternalServerError, "Get memeception by symbol error: ", err)
			return
		}

		ctx.JSON(http.StatusOK, response)
	}
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
