package announcement

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
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

type DeleteAnnouncementLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete announcement
func NewDeleteAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAnnouncementLogic {
	return &DeleteAnnouncementLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAnnouncementLogic) DeleteAnnouncement(req *types.DeleteAnnouncementRequest) error {
	if err := l.svcCtx.AnnouncementModel.Delete(l.ctx, req.Id); err != nil {
		l.Errorw("[DeleteAnnouncement] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete announcement failed: %v", err.Error())
	}
	return nil
}
