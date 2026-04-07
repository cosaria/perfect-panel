package authMethod

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/notify/email"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetEmailPlatformOutput struct {
	Body *types.PlatformResponse
}

func GetEmailPlatformHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetEmailPlatformOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetEmailPlatformOutput, error) {
		l := NewGetEmailPlatformLogic(ctx, svcCtx)
		resp, err := l.GetEmailPlatform()
		if err != nil {
			return nil, err
		}
		return &GetEmailPlatformOutput{Body: resp}, nil
	}
}

type GetEmailPlatformLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get email support platform
func NewGetEmailPlatformLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmailPlatformLogic {
	return &GetEmailPlatformLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmailPlatformLogic) GetEmailPlatform() (resp *types.PlatformResponse, err error) {
	return &types.PlatformResponse{
		List: email.GetSupportedPlatforms(),
	}, nil
}
