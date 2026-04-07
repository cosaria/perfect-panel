package announcement

import (
	"context"
	"github.com/perfect-panel/server/models/announcement"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
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

type CreateAnnouncementLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create announcement
func NewCreateAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAnnouncementLogic {
	return &CreateAnnouncementLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAnnouncementLogic) CreateAnnouncement(req *types.CreateAnnouncementRequest) error {

	if err := l.svcCtx.AnnouncementModel.Insert(l.ctx, &announcement.Announcement{
		Title:   req.Title,
		Content: req.Content,
	}); err != nil {
		l.Errorw("[CreateAnnouncement] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "create announcement failed: %v", err.Error())
	}

	return nil
}
