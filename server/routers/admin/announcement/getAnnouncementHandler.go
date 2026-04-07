// huma:migrated
package announcement

import (
	"context"
	"github.com/perfect-panel/server/services/admin/announcement"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetAnnouncementInput struct {
	types.GetAnnouncementRequest
}

type GetAnnouncementOutput struct {
	Body *types.Announcement
}

func GetAnnouncementHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetAnnouncementInput) (*GetAnnouncementOutput, error) {
	return func(ctx context.Context, input *GetAnnouncementInput) (*GetAnnouncementOutput, error) {
		l := announcement.NewGetAnnouncementLogic(ctx, svcCtx)
		resp, err := l.GetAnnouncement(&input.GetAnnouncementRequest)
		if err != nil {
			return nil, err
		}
		return &GetAnnouncementOutput{Body: resp}, nil
	}
}
