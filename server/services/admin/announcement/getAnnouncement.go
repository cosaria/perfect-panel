package announcement

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetAnnouncementInput struct {
	types.GetAnnouncementRequest
}

type GetAnnouncementOutput struct {
	Body *types.Announcement
}

func GetAnnouncementHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetAnnouncementInput) (*GetAnnouncementOutput, error) {
	return func(ctx context.Context, input *GetAnnouncementInput) (*GetAnnouncementOutput, error) {
		l := NewGetAnnouncementLogic(ctx, svcCtx)
		resp, err := l.GetAnnouncement(&input.GetAnnouncementRequest)
		if err != nil {
			return nil, err
		}
		return &GetAnnouncementOutput{Body: resp}, nil
	}
}

type GetAnnouncementLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get announcement
func NewGetAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAnnouncementLogic {
	return &GetAnnouncementLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAnnouncementLogic) GetAnnouncement(req *types.GetAnnouncementRequest) (resp *types.Announcement, err error) {
	info, err := l.svcCtx.AnnouncementModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[GetAnnouncement] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get announcement error: %v", err.Error())
	}
	resp = &types.Announcement{}
	tool.DeepCopy(resp, info)
	return
}
