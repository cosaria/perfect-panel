package system

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetSiteConfigOutput struct {
	Body *types.SiteConfig
}

func GetSiteConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetSiteConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSiteConfigOutput, error) {
		l := NewGetSiteConfigLogic(ctx, deps)
		resp, err := l.GetSiteConfig()
		if err != nil {
			return nil, err
		}
		return &GetSiteConfigOutput{Body: resp}, nil
	}
}

type GetSiteConfigLogic struct {
	ctx  context.Context
	deps Deps
	logger.Logger
}

func NewGetSiteConfigLogic(ctx context.Context, deps Deps) *GetSiteConfigLogic {
	return &GetSiteConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetSiteConfigLogic) GetSiteConfig() (resp *types.SiteConfig, err error) {
	resp = &types.SiteConfig{}
	// get site config from db
	siteConfigs, err := l.deps.SystemModel.GetSiteConfig(l.ctx)
	if err != nil {
		l.Error("[GetSiteConfig] Database query error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get site config failed: %v", err.Error())
	}
	// reflect to response
	tool.SystemConfigSliceReflectToStruct(siteConfigs, resp)
	return resp, nil
}
