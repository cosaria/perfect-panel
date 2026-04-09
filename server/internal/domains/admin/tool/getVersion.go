package tool

import (
	"context"
	"fmt"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
)

type GetVersionOutput struct {
	Body *types.VersionResponse
}

func GetVersionHandler(deps Deps) func(context.Context, *struct{}) (*GetVersionOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetVersionOutput, error) {
		l := NewGetVersionLogic(ctx, deps)
		resp, err := l.GetVersion()
		if err != nil {
			return nil, err
		}
		return &GetVersionOutput{Body: resp}, nil
	}
}

type GetVersionLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetVersionLogic Get Version
func NewGetVersionLogic(ctx context.Context, deps Deps) *GetVersionLogic {
	return &GetVersionLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
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
