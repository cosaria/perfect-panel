// huma:migrated
package authMethod

import (
	"context"
	"github.com/perfect-panel/server/services/admin/authMethod"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type TestSmsSendInput struct {
	Body types.TestSmsSendRequest
}

func TestSmsSendHandler(svcCtx *svc.ServiceContext) func(context.Context, *TestSmsSendInput) (*struct{}, error) {
	return func(ctx context.Context, input *TestSmsSendInput) (*struct{}, error) {
		l := authMethod.NewTestSmsSendLogic(ctx, svcCtx)
		if err := l.TestSmsSend(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
