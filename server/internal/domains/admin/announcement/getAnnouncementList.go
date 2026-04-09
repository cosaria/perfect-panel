package announcement

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/announcement"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type GetAnnouncementListInput struct {
	Body types.GetAnnouncementListRequest
}

type GetAnnouncementListOutput struct {
	Body *types.GetAnnouncementListResponse
}

func GetAnnouncementListHandler(deps Deps) func(context.Context, *GetAnnouncementListInput) (*GetAnnouncementListOutput, error) {
	return func(ctx context.Context, input *GetAnnouncementListInput) (*GetAnnouncementListOutput, error) {
		l := NewGetAnnouncementListLogic(ctx, deps)
		resp, err := l.GetAnnouncementList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetAnnouncementListOutput{Body: resp}, nil
	}
}

type GetAnnouncementListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get announcement list
func NewGetAnnouncementListLogic(ctx context.Context, deps Deps) *GetAnnouncementListLogic {
	return &GetAnnouncementListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetAnnouncementListLogic) GetAnnouncementList(req *types.GetAnnouncementListRequest) (resp *types.GetAnnouncementListResponse, err error) {
	total, list, err := l.deps.AnnouncementModel.GetAnnouncementListByPage(l.ctx, int(req.Page), int(req.Size), announcement.Filter{
		Show:   req.Show,
		Pinned: req.Pinned,
		Popup:  req.Popup,
		Search: req.Search,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetAnnouncementListByPage error: %v", err.Error())
	}
	resp = &types.GetAnnouncementListResponse{}
	resp.Total = total
	resp.List = make([]types.Announcement, 0)
	tool.DeepCopy(&resp.List, list)
	return
}
