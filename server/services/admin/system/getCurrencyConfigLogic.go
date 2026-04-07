package system

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetCurrencyConfigLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Currency Config
func NewGetCurrencyConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrencyConfigLogic {
	return &GetCurrencyConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrencyConfigLogic) GetCurrencyConfig() (resp *types.CurrencyConfig, err error) {
	configs, err := l.svcCtx.SystemModel.GetCurrencyConfig(l.ctx)
	if err != nil {
		l.Errorw("[GetCurrencyConfigLogic] GetCurrencyConfig error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetCurrencyConfig error: %v", err.Error())
	}
	resp = &types.CurrencyConfig{}
	tool.SystemConfigSliceReflectToStruct(configs, resp)
	return
}
