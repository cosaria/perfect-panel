// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/services/admin/marketing"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateQuotaTaskInput struct {
	Body types.CreateQuotaTaskRequest
}

func CreateQuotaTaskHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateQuotaTaskInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateQuotaTaskInput) (*struct{}, error) {
		l := marketing.NewCreateQuotaTaskLogic(ctx, svcCtx)
		if err := l.CreateQuotaTask(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
