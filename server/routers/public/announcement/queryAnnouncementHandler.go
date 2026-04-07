// huma:migrated
package announcement

import (
	"context"
	"github.com/perfect-panel/server/services/user/announcement"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryAnnouncementInput struct {
	Body types.QueryAnnouncementRequest
}

type QueryAnnouncementOutput struct {
	Body *types.QueryAnnouncementResponse
}

func QueryAnnouncementHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryAnnouncementInput) (*QueryAnnouncementOutput, error) {
	return func(ctx context.Context, input *QueryAnnouncementInput) (*QueryAnnouncementOutput, error) {
		l := announcement.NewQueryAnnouncementLogic(ctx, svcCtx)
		resp, err := l.QueryAnnouncement(&input.Body)
		if err != nil {
			return nil, err
		}
		return &QueryAnnouncementOutput{Body: resp}, nil
	}
}
