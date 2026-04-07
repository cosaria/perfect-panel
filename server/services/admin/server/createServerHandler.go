// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateServerInput struct {
	Body types.CreateServerRequest
}

func CreateServerHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateServerInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateServerInput) (*struct{}, error) {
		l := NewCreateServerLogic(ctx, svcCtx)
		if err := l.CreateServer(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
