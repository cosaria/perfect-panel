package authMethod

import (
	"context"
	"fmt"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/notify/email"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type TestEmailSendInput struct {
	Body types.TestEmailSendRequest
}

func TestEmailSendHandler(deps Deps) func(context.Context, *TestEmailSendInput) (*struct{}, error) {
	return func(ctx context.Context, input *TestEmailSendInput) (*struct{}, error) {
		l := NewTestEmailSendLogic(ctx, deps)
		if err := l.TestEmailSend(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type TestEmailSendLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Test email send
func NewTestEmailSendLogic(ctx context.Context, deps Deps) *TestEmailSendLogic {
	return &TestEmailSendLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *TestEmailSendLogic) TestEmailSend(req *types.TestEmailSendRequest) error {
	cfg := l.deps.currentConfig()
	client, err := email.NewSender(cfg.Email.Platform, cfg.Email.PlatformConfig, cfg.Site.SiteName)
	if err != nil {
		l.Errorw("new email sender err", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "new email sender err: %v", err.Error())
	}
	err = client.Send([]string{req.Email}, "Test Email Send", "this a test email send by ppanel")
	if err != nil {
		return errors.Wrapf(xerr.NewErrCodeMsg(500, fmt.Sprintf("send email err: %v", err.Error())), "send email err: %v", err.Error())
	}
	return nil
}
