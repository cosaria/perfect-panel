// huma:migrated
package authMethod

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/authMethod"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type TestEmailSendInput struct {
	Body types.TestEmailSendRequest
}

func TestEmailSendHandler(svcCtx *svc.ServiceContext) func(context.Context, *TestEmailSendInput) (*struct{}, error) {
	return func(ctx context.Context, input *TestEmailSendInput) (*struct{}, error) {
		l := authMethod.NewTestEmailSendLogic(ctx, svcCtx)
		if err := l.TestEmailSend(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
