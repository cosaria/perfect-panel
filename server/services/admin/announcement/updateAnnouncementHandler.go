// huma:migrated
package announcement

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateAnnouncementInput struct {
	Body types.UpdateAnnouncementRequest
}

func UpdateAnnouncementHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateAnnouncementInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateAnnouncementInput) (*struct{}, error) {
		l := NewUpdateAnnouncementLogic(ctx, svcCtx)
		if err := l.UpdateAnnouncement(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
