package coupon

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type GetCouponListInput struct {
	types.GetCouponListRequest
}

type GetCouponListOutput struct {
	Body *types.GetCouponListResponse
}

func GetCouponListHandler(deps Deps) func(context.Context, *GetCouponListInput) (*GetCouponListOutput, error) {
	return func(ctx context.Context, input *GetCouponListInput) (*GetCouponListOutput, error) {
		l := NewGetCouponListLogic(ctx, deps)
		resp, err := l.GetCouponList(&input.GetCouponListRequest)
		if err != nil {
			return nil, err
		}
		return &GetCouponListOutput{Body: resp}, nil
	}
}

type GetCouponListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get coupon list
func NewGetCouponListLogic(ctx context.Context, deps Deps) *GetCouponListLogic {
	return &GetCouponListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetCouponListLogic) GetCouponList(req *types.GetCouponListRequest) (resp *types.GetCouponListResponse, err error) {
	resp = &types.GetCouponListResponse{}
	// get coupon list from db
	total, list, err := l.deps.CouponModel.QueryCouponListByPage(l.ctx, int(req.Page), int(req.Size), req.Subscribe, req.Search)
	if err != nil {
		l.Errorw("[GetCouponList] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get coupon list error: %v", err.Error())
	}
	resp.Total = total
	resp.List = make([]types.Coupon, 0)
	for _, coupon := range list {
		couponInfo := types.Coupon{}
		tool.DeepCopy(&couponInfo, coupon)
		couponInfo.Subscribe = tool.StringToInt64Slice(coupon.Subscribe)
		resp.List = append(resp.List, couponInfo)
	}
	return
}
