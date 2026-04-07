// huma:migrated
package announcement

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetAnnouncementListInput struct {
	Body types.GetAnnouncementListRequest
}

type GetAnnouncementListOutput struct {
	Body *types.GetAnnouncementListResponse
}

func GetAnnouncementListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetAnnouncementListInput) (*GetAnnouncementListOutput, error) {
	return func(ctx context.Context, input *GetAnnouncementListInput) (*GetAnnouncementListOutput, error) {
		l := NewGetAnnouncementListLogic(ctx, svcCtx)
		resp, err := l.GetAnnouncementList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetAnnouncementListOutput{Body: resp}, nil
	}
}
