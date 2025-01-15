package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	resp := Response{Code: errorCodeMap[ErrSuccess], Message: ErrSuccess.Error(), Data: data}
	if _, ok := errorCodeMap[ErrSuccess]; !ok {
		resp = Response{Code: 0, Message: "", Data: data}
	}
	ctx.JSON(http.StatusOK, resp)
}

func HandleError(ctx *gin.Context, httpCode int, err error, data interface{}) {
	// 优先级 httpCode > Error.HttpCode > 500
	if data == nil {
		data = map[string]string{}
	}
	if httpCode <= 0 {
		httpCode = errorHttpCodeMap[err]
	}
	if httpCode == 0 {
		httpCode = http.StatusInternalServerError
	}
	resp := Response{Code: errorCodeMap[err], Message: err.Error(), Data: data}
	if _, ok := errorCodeMap[err]; !ok {
		resp = Response{Code: 500, Message: "unknown error", Data: data}
	}
	ctx.JSON(httpCode, resp)
}

type Error struct {
	HttpCode int
	Code     int
	Message  string
}

var errorCodeMap = map[error]int{}
var errorHttpCodeMap = map[error]int{}

func newError(httpCode int, code int, msg string) error {
	err := errors.New(msg)
	errorCodeMap[err] = code
	errorHttpCodeMap[err] = httpCode
	return err
}
func (e Error) Error() string {
	return e.Message
}
