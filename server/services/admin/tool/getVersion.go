package tool

import (
	"context"
	"fmt"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetVersionOutput struct {
	Body *types.VersionResponse
}

func GetVersionHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetVersionOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetVersionOutput, error) {
		l := NewGetVersionLogic(ctx, svcCtx)
		resp, err := l.GetVersion()
		if err != nil {
			return nil, err
		}
		return &GetVersionOutput{Body: resp}, nil
	}
}

type GetVersionLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetVersionLogic Get Version
func NewGetVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVersionLogic {
	return &GetVersionLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVersionLogic) GetVersion() (resp *types.VersionResponse, err error) {
	version := config.Version
	buildTime := config.BuildTime

	// Normalize unknown values
	if version == "unknown version" {
		version = "unknown"
	}
	if buildTime == "unknown time" {
		buildTime = "unknown"
	}

	// Format version based on whether it starts with 'v'
	var formattedVersion string
	if len(version) > 0 && version[0] == 'v' {
		formattedVersion = fmt.Sprintf("%s(%s)", version[1:], buildTime)
	} else {
		formattedVersion = fmt.Sprintf("%s(%s) Develop", version, buildTime)
	}

	return &types.VersionResponse{
		Version: formattedVersion,
	}, nil
}
