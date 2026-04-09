package authMethod

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type GetAuthMethodConfigInput struct {
	types.GetAuthMethodConfigRequest
}

type GetAuthMethodConfigOutput struct {
	Body *types.AuthMethodConfig
}

func GetAuthMethodConfigHandler(deps Deps) func(context.Context, *GetAuthMethodConfigInput) (*GetAuthMethodConfigOutput, error) {
	return func(ctx context.Context, input *GetAuthMethodConfigInput) (*GetAuthMethodConfigOutput, error) {
		l := NewGetAuthMethodConfigLogic(ctx, deps)
		resp, err := l.GetAuthMethodConfig(&input.GetAuthMethodConfigRequest)
		if err != nil {
			return nil, err
		}
		return &GetAuthMethodConfigOutput{Body: resp}, nil
	}
}

type GetAuthMethodConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetAuthMethodConfigLogic Get auth method config
func NewGetAuthMethodConfigLogic(ctx context.Context, deps Deps) *GetAuthMethodConfigLogic {
	return &GetAuthMethodConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetAuthMethodConfigLogic) GetAuthMethodConfig(req *types.GetAuthMethodConfigRequest) (resp *types.AuthMethodConfig, err error) {
	method, err := l.deps.AuthModel.FindOneByMethod(l.ctx, req.Method)
	if err != nil {
		l.Errorw("find one by method failed", logger.Field("method", req.Method), logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find one by method failed: %v", err.Error())
	}

	resp = new(types.AuthMethodConfig)
	tool.DeepCopy(resp, method)
	if method.Config != "" {
		if err := json.Unmarshal([]byte(method.Config), &resp.Config); err != nil {
			l.Errorw("unmarshal config failed", logger.Field("config", method.Config), logger.Field("error", err.Error()))
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "unmarshal apple config failed: %v", err.Error())
		}
	}
	return
}
