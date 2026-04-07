package subscribe

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetSubscribeDetailsLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get subscribe details
func NewGetSubscribeDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscribeDetailsLogic {
	return &GetSubscribeDetailsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSubscribeDetailsLogic) GetSubscribeDetails(req *types.GetSubscribeDetailsRequest) (resp *types.Subscribe, err error) {
	sub, err := l.svcCtx.SubscribeModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Logger.Error("[GetSubscribeDetailsLogic] get subscribe details failed: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get subscribe details failed: %v", err.Error())
	}
	resp = &types.Subscribe{}
	tool.DeepCopy(resp, sub)
	if sub.Discount != "" {
		err = json.Unmarshal([]byte(sub.Discount), &resp.Discount)
		if err != nil {
			l.Logger.Error("[GetSubscribeDetailsLogic] JSON unmarshal failed: ", logger.Field("error", err.Error()), logger.Field("discount", sub.Discount))
		}
	}
	resp.Nodes = tool.StringToInt64Slice(sub.Nodes)
	resp.NodeTags = strings.Split(sub.NodeTags, ",")
	return resp, nil
}
