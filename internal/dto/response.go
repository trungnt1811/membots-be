package dto

import (
	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/gin-gonic/gin"
)

type ErrorMap struct {
	Errors map[string]interface{} `json:"errors"`
}

type ResponseData struct {
	Status   int         `json:"status"`
	Data     interface{} `json:"data"`
	IsCached bool        `json:"isCached"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Status int         `json:"status"`
	Errors interface{} `json:"errors"`
}

type GeneralError struct {
	// HTTP error code, or custom error code
	Code int `json:"code"`
	// Friendly error message
	Message string `json:"string"`
	// List of error send server 2 server
	Errors []string `json:"errors"`
}

type ServiceError struct {
	// Custom error code
	Code int `json:"code"`
	// HTTP error code
	Status int `json:"status"`
	// error throw by service
	Err error `json:"error"`
}

func RespondJSON(w *gin.Context, status int, payload interface{}, isCached bool) {
	var res ResponseData
	res.Status = status
	res.Data = payload
	res.IsCached = isCached
	w.Abort()
	w.JSON(status, res)
}

func RespondErrors(w *gin.Context, status int, payload interface{}) {
	var res ResponseError
	res.Status = status
	res.Errors = payload
	w.Abort()
	w.JSON(status, res)
}

func RespondErrorWithCode(w *gin.Context, status int, code int, msg string, errors ...error) {
	var errResp util.GeneralError
	errResp.Code = code
	errResp.Message = msg
	errResp.Errors = []string{}
	for _, err := range errors {
		if err != nil {
			errResp.Errors = append(errResp.Errors, err.Error())
		}
	}

	w.Abort()
	w.JSON(status, errResp)
}

func NewError(key string, err error) map[string]interface{} {
	res := ErrorMap{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res.Errors
}
