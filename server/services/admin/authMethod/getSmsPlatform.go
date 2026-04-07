package authMethod

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/notify/sms"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetSmsPlatformOutput struct {
	Body *types.PlatformResponse
}

func GetSmsPlatformHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetSmsPlatformOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSmsPlatformOutput, error) {
		l := NewGetSmsPlatformLogic(ctx, svcCtx)
		resp, err := l.GetSmsPlatform()
		if err != nil {
			return nil, err
		}
		return &GetSmsPlatformOutput{Body: resp}, nil
	}
}

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
