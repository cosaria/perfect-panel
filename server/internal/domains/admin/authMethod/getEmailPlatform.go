package authMethod

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/notify/email"
)

type GetEmailPlatformOutput struct {
	Body *types.PlatformResponse
}

func GetEmailPlatformHandler(deps Deps) func(context.Context, *struct{}) (*GetEmailPlatformOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetEmailPlatformOutput, error) {
		l := NewGetEmailPlatformLogic(ctx, deps)
		resp, err := l.GetEmailPlatform()
		if err != nil {
			return nil, err
		}
		return &GetEmailPlatformOutput{Body: resp}, nil
	}
}

type GetEmailPlatformLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get email support platform
func NewGetEmailPlatformLogic(ctx context.Context, deps Deps) *GetEmailPlatformLogic {
	return &GetEmailPlatformLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetEmailPlatformLogic) GetEmailPlatform() (resp *types.PlatformResponse, err error) {
	return &types.PlatformResponse{
		List: email.GetSupportedPlatforms(),
	}, nil
}
