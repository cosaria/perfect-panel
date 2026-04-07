package oauth

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"net/http"
	"net/url"
)

func AppleLoginCallbackHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.AppleLoginCallbackRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}
		l := NewAppleLoginCallbackLogic(c.Request.Context(), svcCtx)
		err := l.AppleLoginCallback(&req, c.Request, c.Writer)
		if err != nil {
			response.HttpResult(c, nil, err)
		}
	}
}

type AppleLoginCallbackLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Apple Login Callback
func NewAppleLoginCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppleLoginCallbackLogic {
	return &AppleLoginCallbackLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppleLoginCallbackLogic) AppleLoginCallback(req *types.AppleLoginCallbackRequest, r *http.Request, w http.ResponseWriter) error {
	// validate the state code
	result, err := l.svcCtx.Redis.Get(l.ctx, fmt.Sprintf("apple:%s", req.State)).Result()
	if err != nil {
		l.Errorw("get apple state code from redis failed", logger.Field("error", err.Error()), logger.Field("code", req.State))
		http.Redirect(w, r, l.svcCtx.Config.Site.Host, http.StatusTemporaryRedirect)
		return nil
	}
	http.Redirect(w, r, fmt.Sprintf("%s?method=apple&code=%s&state=%s", result, req.Code, req.State), http.StatusFound)
	l.Infow("redirect to apple login page", logger.Field("url", fmt.Sprintf("%s?method=apple&code=%s&state=%s", result, url.QueryEscape(req.Code), req.State)))
	return nil
}
