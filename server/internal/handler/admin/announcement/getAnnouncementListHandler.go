// huma:migrated
package announcement

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/announcement"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetAnnouncementListInput struct {
	Body types.GetAnnouncementListRequest
}

type GetAnnouncementListOutput struct {
	Body *types.GetAnnouncementListResponse
}

func GetAnnouncementListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetAnnouncementListInput) (*GetAnnouncementListOutput, error) {
	return func(ctx context.Context, input *GetAnnouncementListInput) (*GetAnnouncementListOutput, error) {
		l := announcement.NewGetAnnouncementListLogic(ctx, svcCtx)
		resp, err := l.GetAnnouncementList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetAnnouncementListOutput{Body: resp}, nil
	}
}
