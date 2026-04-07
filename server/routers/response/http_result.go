package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/perfect-panel/server/modules/infra/xerr"
)

// HttpResult writes the standard success or error envelope.
func HttpResult(ctx *gin.Context, resp interface{}, err error) {
	if err == nil {
		ctx.JSON(http.StatusOK, Success(resp))
		return
	}

	code := xerr.ERROR
	msg := "Internal Server Error"

	var e *xerr.CodeError
	if errors.As(errors.Cause(err), &e) {
		code = e.GetErrCode()
		msg = e.GetErrMsg()
	}
	ctx.JSON(http.StatusOK, Error(code, msg))
}

// ParamErrorResult writes the standard invalid-params envelope.
func ParamErrorResult(ctx *gin.Context, err error) {
	errMsg := err.Error()
	_ = ctx.Error(errors.New(errMsg))
	ctx.JSON(http.StatusOK, Error(xerr.InvalidParams, errMsg))
}
