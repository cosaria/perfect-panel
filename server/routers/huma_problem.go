package handler

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/svc"
)

var defaultOperationErrors = []int{
	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusForbidden,
	http.StatusNotFound,
	http.StatusTooManyRequests,
	http.StatusBadGateway,
}

func configureHumaAPI(api huma.API, compatibilityMode bool) {
	if !compatibilityMode {
		return
	}

	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		next(huma.WithContext(ctx, response.WithHumaCompatibilityMode(ctx.Context(), true)))
	})
}

func compatibilityEnabled(serverCtx *svc.ServiceContext, specOnly bool) bool {
	return !specOnly && serverCtx != nil && serverCtx.Config.ErrorCompatibilityMode
}

func registerOperation[I, O any](api huma.API, op huma.Operation, handler func(context.Context, *I) (*O, error)) {
	if len(op.Errors) == 0 {
		op.Errors = append([]int(nil), defaultOperationErrors...)
	}

	huma.Register(api, op, func(ctx context.Context, input *I) (*O, error) {
		output, err := handler(ctx, input)
		if err != nil {
			return nil, response.AsHumaStatusError(ctx, err)
		}
		return output, nil
	})
}
