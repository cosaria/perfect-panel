package announcement

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
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

type UpdateAnnouncementLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update announcement
func NewUpdateAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAnnouncementLogic {
	return &UpdateAnnouncementLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAnnouncementLogic) UpdateAnnouncement(req *types.UpdateAnnouncementRequest) error {
	info, err := l.svcCtx.AnnouncementModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[UpdateAnnouncement] Query Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get announcement error: %v", err.Error())
	}
	info.Title = req.Title
	info.Content = req.Content
	if req.Show != nil {
		info.Show = req.Show
	}
	if req.Pinned != nil {
		info.Pinned = req.Pinned
	}
	if req.Popup != nil {
		info.Popup = req.Popup
	}
	err = l.svcCtx.AnnouncementModel.Update(l.ctx, info)
	if err != nil {
		l.Errorw("[UpdateAnnouncement] Update Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update announcement error: %v", err.Error())
	}
	return nil
}
