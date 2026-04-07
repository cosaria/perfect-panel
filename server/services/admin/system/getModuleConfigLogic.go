package system

import (
	"context"
	"os"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetModuleConfigLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Module Config
func NewGetModuleConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetModuleConfigLogic {
	return &GetModuleConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetModuleConfigLogic) GetModuleConfig() (resp *types.ModuleConfig, err error) {
	value, exists := os.LookupEnv("SECRET_KEY")
	if !exists {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), " SECRET_KEY not set in environment variables")
	}

	return &types.ModuleConfig{
		Secret:         value,
		ServiceName:    config.ServiceName,
		ServiceVersion: config.Version,
	}, nil
}
