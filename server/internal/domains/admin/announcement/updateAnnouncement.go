package announcement

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
)

type UpdateAnnouncementInput struct {
	Body types.UpdateAnnouncementRequest
}

func UpdateAnnouncementHandler(deps Deps) func(context.Context, *UpdateAnnouncementInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateAnnouncementInput) (*struct{}, error) {
		l := NewUpdateAnnouncementLogic(ctx, deps)
		if err := l.UpdateAnnouncement(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateAnnouncementLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update announcement
func NewUpdateAnnouncementLogic(ctx context.Context, deps Deps) *UpdateAnnouncementLogic {
	return &UpdateAnnouncementLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateAnnouncementLogic) UpdateAnnouncement(req *types.UpdateAnnouncementRequest) error {
	info, err := l.deps.AnnouncementModel.FindOne(l.ctx, req.Id)
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
	err = l.deps.AnnouncementModel.Update(l.ctx, info)
	if err != nil {
		l.Errorw("[UpdateAnnouncement] Update Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update announcement error: %v", err.Error())
	}
	return nil
}
