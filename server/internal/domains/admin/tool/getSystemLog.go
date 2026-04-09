package tool

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
)

type GetSystemLogOutput struct {
	Body *types.LogResponse
}

func GetSystemLogHandler(deps Deps) func(context.Context, *struct{}) (*GetSystemLogOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSystemLogOutput, error) {
		l := NewGetSystemLogLogic(ctx, deps)
		resp, err := l.GetSystemLog()
		if err != nil {
			return nil, err
		}
		return &GetSystemLogOutput{Body: resp}, nil
	}
}

type GetSystemLogLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetSystemLogLogic Get System Log
func NewGetSystemLogLogic(ctx context.Context, deps Deps) *GetSystemLogLogic {
	return &GetSystemLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetSystemLogLogic) GetSystemLog() (resp *types.LogResponse, err error) {
	if l.deps.Config == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "system log path is not configured")
	}

	lines, err := logger.ReadLastNLines(l.deps.Config.Logger.Path, 50)
	if err != nil {
		l.Error(err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "get system log error: %v", err.Error())
	}
	var list []map[string]interface{}
	for _, line := range lines {
		var log map[string]interface{}
		if err = json.Unmarshal([]byte(line), &log); err != nil {
			l.Error(err)
			continue
		}
		list = append(list, log)
	}

	return &types.LogResponse{
		List: list,
	}, nil
}
