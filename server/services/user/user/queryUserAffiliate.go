package user

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type QueryUserAffiliateOutput struct {
	Body *types.QueryUserAffiliateCountResponse
}

func QueryUserAffiliateHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryUserAffiliateOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserAffiliateOutput, error) {
		l := NewQueryUserAffiliateLogic(ctx, svcCtx)
		resp, err := l.QueryUserAffiliate()
		if err != nil {
			return nil, err
		}
		return &QueryUserAffiliateOutput{Body: resp}, nil
	}
}

type QueryUserAffiliateLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Query User Balance Log
func NewQueryUserAffiliateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserAffiliateLogic {
	return &QueryUserAffiliateLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryUserAffiliateLogic) QueryUserAffiliate() (resp *types.QueryUserAffiliateCountResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	var sum int64
	var total int64
	err = l.svcCtx.UserModel.Transaction(l.ctx, func(db *gorm.DB) error {
		return db.Model(&user.User{}).Where("referer_id = ?", u.Id).Count(&total).Find(&user.User{}).Error
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query User Affiliate failed: %v", err)
	}
	data, _, err := l.svcCtx.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     1,
		Size:     99999,
		Type:     log.TypeCommission.Uint8(),
		ObjectID: u.Id,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query User Affiliate logs failed: %v", err)
	}

	for _, datum := range data {
		content := log.Commission{}
		if err = content.Unmarshal([]byte(datum.Content)); err != nil {
			l.Errorf("[QueryUserAffiliate] unmarshal comission log failed: %v", err.Error())
			continue
		}
		sum += content.Amount
	}

	return &types.QueryUserAffiliateCountResponse{
		Registers:       total,
		TotalCommission: sum,
	}, nil
}
