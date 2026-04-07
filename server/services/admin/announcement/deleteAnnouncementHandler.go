// huma:migrated
package announcement

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteAnnouncementInput struct {
	Body types.DeleteAnnouncementRequest
}

func DeleteAnnouncementHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteAnnouncementInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteAnnouncementInput) (*struct{}, error) {
		l := NewDeleteAnnouncementLogic(ctx, svcCtx)
		if err := l.DeleteAnnouncement(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
