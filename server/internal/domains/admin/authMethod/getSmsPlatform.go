package authMethod

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/notify/sms"
)

type GetSmsPlatformOutput struct {
	Body *types.PlatformResponse
}

func GetSmsPlatformHandler(deps Deps) func(context.Context, *struct{}) (*GetSmsPlatformOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSmsPlatformOutput, error) {
		l := NewGetSmsPlatformLogic(ctx, deps)
		resp, err := l.GetSmsPlatform()
		if err != nil {
			return nil, err
		}
		return &GetSmsPlatformOutput{Body: resp}, nil
	}
}

type GetSmsPlatformLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get sms support platform
func NewGetSmsPlatformLogic(ctx context.Context, deps Deps) *GetSmsPlatformLogic {
	return &GetSmsPlatformLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetSmsPlatformLogic) GetSmsPlatform() (resp *types.PlatformResponse, err error) {
	return &types.PlatformResponse{
		List: sms.GetSupportedPlatforms(),
	}, nil
}
