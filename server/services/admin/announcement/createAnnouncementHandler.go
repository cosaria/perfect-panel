// huma:migrated
package announcement

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateAnnouncementInput struct {
	Body types.CreateAnnouncementRequest
}

func CreateAnnouncementHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateAnnouncementInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateAnnouncementInput) (*struct{}, error) {
		l := NewCreateAnnouncementLogic(ctx, svcCtx)
		if err := l.CreateAnnouncement(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
