package announcement

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetAnnouncementInput struct {
	types.GetAnnouncementRequest
}

type GetAnnouncementOutput struct {
	Body *types.Announcement
}

func GetAnnouncementHandler(deps Deps) func(context.Context, *GetAnnouncementInput) (*GetAnnouncementOutput, error) {
	return func(ctx context.Context, input *GetAnnouncementInput) (*GetAnnouncementOutput, error) {
		l := NewGetAnnouncementLogic(ctx, deps)
		resp, err := l.GetAnnouncement(&input.GetAnnouncementRequest)
		if err != nil {
			return nil, err
		}
		return &GetAnnouncementOutput{Body: resp}, nil
	}
}

type GetAnnouncementLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get announcement
func NewGetAnnouncementLogic(ctx context.Context, deps Deps) *GetAnnouncementLogic {
	return &GetAnnouncementLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetAnnouncementLogic) GetAnnouncement(req *types.GetAnnouncementRequest) (resp *types.Announcement, err error) {
	info, err := l.deps.AnnouncementModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[GetAnnouncement] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get announcement error: %v", err.Error())
	}
	resp = &types.Announcement{}
	tool.DeepCopy(resp, info)
	return
}
