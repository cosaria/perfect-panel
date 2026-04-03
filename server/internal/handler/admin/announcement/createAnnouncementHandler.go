// huma:migrated
package announcement

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/announcement"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateAnnouncementInput struct {
	Body types.CreateAnnouncementRequest
}

func CreateAnnouncementHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateAnnouncementInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateAnnouncementInput) (*struct{}, error) {
		l := announcement.NewCreateAnnouncementLogic(ctx, svcCtx)
		if err := l.CreateAnnouncement(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
