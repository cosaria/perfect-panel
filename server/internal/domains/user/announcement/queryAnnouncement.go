package announcement

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/announcement"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type QueryAnnouncementInput struct {
	Body types.QueryAnnouncementRequest
}

type QueryAnnouncementOutput struct {
	Body *types.QueryAnnouncementResponse
}

func QueryAnnouncementHandler(deps Deps) func(context.Context, *QueryAnnouncementInput) (*QueryAnnouncementOutput, error) {
	return func(ctx context.Context, input *QueryAnnouncementInput) (*QueryAnnouncementOutput, error) {
		l := NewQueryAnnouncementLogic(ctx, deps)
		resp, err := l.QueryAnnouncement(&input.Body)
		if err != nil {
			return nil, err
		}
		return &QueryAnnouncementOutput{Body: resp}, nil
	}
}

type QueryAnnouncementLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Query announcement
func NewQueryAnnouncementLogic(ctx context.Context, deps Deps) *QueryAnnouncementLogic {
	return &QueryAnnouncementLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryAnnouncementLogic) QueryAnnouncement(req *types.QueryAnnouncementRequest) (resp *types.QueryAnnouncementResponse, err error) {
	enable := true
	total, list, err := l.deps.AnnouncementModel.GetAnnouncementListByPage(l.ctx, req.Page, req.Size, announcement.Filter{
		Show:   &enable,
		Pinned: req.Pinned,
		Popup:  req.Popup,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetAnnouncementListByPage error: %v", err.Error())
	}
	resp = &types.QueryAnnouncementResponse{}
	resp.Total = total
	resp.List = make([]types.Announcement, 0)
	tool.DeepCopy(&resp.List, list)
	return
}
