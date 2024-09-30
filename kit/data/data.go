package data

import (
	"net/http"

	"github.com/diycoder/elf/kit/errors"
	"github.com/gin-gonic/gin"
)

const (
	defaultError     = "unknown err"
	defaultMsg       = "ok"
	defaultErrorCode = 500
)

type JSONResult struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type PageResult struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

type PageInfo struct {
	PageIndex int64 `json:"page_index" form:"page_index" query:"page_index"`
	PageSize  int64 `json:"page_size" form:"page_size" query:"page_size"`
}

func Error(ctx *gin.Context, err error) {
	code, detail := http.StatusInternalServerError, defaultError
	if err != nil {
		e := errors.Parse((err).Error())
		if e.Code >= 0 {
			code = e.Code
		}
		detail = e.Detail
	}
	ctx.JSON(http.StatusOK, &JSONResult{
		Code:    code,
		Message: detail,
		Data:    "",
	})
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &JSONResult{
		Code:    http.StatusOK,
		Message: defaultMsg,
		Data:    data,
	})
}

func HTML(ctx *gin.Context, name string, data interface{}) {
	ctx.HTML(http.StatusOK, name, data)
}

func File(ctx *gin.Context, file, name string) {
	ctx.FileAttachment(file, name)
}
