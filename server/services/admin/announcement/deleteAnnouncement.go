package announcement

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type DeleteAnnouncementInput struct {
	Body types.DeleteAnnouncementRequest
}

func DeleteAnnouncementHandler(deps Deps) func(context.Context, *DeleteAnnouncementInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteAnnouncementInput) (*struct{}, error) {
		l := NewDeleteAnnouncementLogic(ctx, deps)
		if err := l.DeleteAnnouncement(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteAnnouncementLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Delete announcement
func NewDeleteAnnouncementLogic(ctx context.Context, deps Deps) *DeleteAnnouncementLogic {
	return &DeleteAnnouncementLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeleteAnnouncementLogic) DeleteAnnouncement(req *types.DeleteAnnouncementRequest) error {
	if err := l.deps.AnnouncementModel.Delete(l.ctx, req.Id); err != nil {
		l.Errorw("[DeleteAnnouncement] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete announcement failed: %v", err.Error())
	}
	return nil
}
