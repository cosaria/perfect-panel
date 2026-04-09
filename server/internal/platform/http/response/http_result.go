package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HttpResult writes the standard success or error envelope.
func HttpResult(ctx *gin.Context, resp interface{}, err error) {
	if err == nil {
		ctx.JSON(http.StatusOK, Success(resp))
		return
	}

	WriteProblem(ctx, NewProblemFromError(err))
}

// ParamErrorResult writes the standard invalid-params envelope.
func ParamErrorResult(ctx *gin.Context, err error) {
	WriteProblem(ctx, NewValidationProblem(err))
}
