package authMethod

import (
	"context"

	"github.com/perfect-panel/server/pkg/sms"

	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetSmsPlatformLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get sms support platform
func NewGetSmsPlatformLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSmsPlatformLogic {
	return &GetSmsPlatformLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSmsPlatformLogic) GetSmsPlatform() (resp *types.PlatformResponse, err error) {
	return &types.PlatformResponse{
		List: sms.GetSupportedPlatforms(),
	}, nil
}
