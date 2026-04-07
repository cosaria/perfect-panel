// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateNodeInput struct {
	Body types.CreateNodeRequest
}

func CreateNodeHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateNodeInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateNodeInput) (*struct{}, error) {
		l := NewCreateNodeLogic(ctx, svcCtx)
		if err := l.CreateNode(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
