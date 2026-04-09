package authMethod

import (
	"context"
	"fmt"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/notify/sms"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type TestSmsSendInput struct {
	Body types.TestSmsSendRequest
}

func TestSmsSendHandler(deps Deps) func(context.Context, *TestSmsSendInput) (*struct{}, error) {
	return func(ctx context.Context, input *TestSmsSendInput) (*struct{}, error) {
		l := NewTestSmsSendLogic(ctx, deps)
		if err := l.TestSmsSend(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type TestSmsSendLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Test sms send
func NewTestSmsSendLogic(ctx context.Context, deps Deps) *TestSmsSendLogic {
	return &TestSmsSendLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *TestSmsSendLogic) TestSmsSend(req *types.TestSmsSendRequest) error {
	cfg := l.deps.currentConfig()
	client, err := sms.NewSender(cfg.Mobile.Platform, cfg.Mobile.PlatformConfig)
	if err != nil {
		l.Errorw("new sms sender err", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "new sms sender err: %v", err.Error())
	}
	err = client.SendCode(req.AreaCode, req.Telephone, "123456")
	if err != nil {
		l.Errorw("send sms err", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCodeMsg(500, fmt.Sprintf("send sms err: %v", err.Error())), "send sms err: %v", err.Error())
	}
	return nil
}
