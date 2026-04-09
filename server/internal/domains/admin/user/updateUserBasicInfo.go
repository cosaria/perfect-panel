package user

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type UpdateUserBasicInfoInput struct {
	Body types.UpdateUserBasiceInfoRequest
}

func UpdateUserBasicInfoHandler(deps Deps) func(context.Context, *UpdateUserBasicInfoInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserBasicInfoInput) (*struct{}, error) {
		l := NewUpdateUserBasicInfoLogic(ctx, deps)
		if err := l.UpdateUserBasicInfo(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateUserBasicInfoLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewUpdateUserBasicInfoLogic Update user basic info
func NewUpdateUserBasicInfoLogic(ctx context.Context, deps Deps) *UpdateUserBasicInfoLogic {
	return &UpdateUserBasicInfoLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateUserBasicInfoLogic) UpdateUserBasicInfo(req *types.UpdateUserBasiceInfoRequest) error {
	userInfo, err := l.deps.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		l.Errorw("[UpdateUserBasicInfoLogic] Find User Error:", logger.Field("err", err.Error()), logger.Field("userId", req.UserId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Find User Error")
	}

	isDemo := strings.ToLower(os.Getenv("PPANEL_MODE")) == "demo"

	if req.Avatar != "" && !tool.IsValidImageSize(req.Avatar, 1024) {
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Invalid Image Size")
	}

	if userInfo.Balance != req.Balance {
		change := req.Balance - userInfo.Balance
		balanceLog := log.Balance{
			Type:      log.BalanceTypeAdjust,
			Amount:    change,
			OrderNo:   "",
			Balance:   req.Balance,
			Timestamp: time.Now().UnixMilli(),
		}
		content, _ := balanceLog.Marshal()

		err = l.deps.LogModel.Insert(l.ctx, &log.SystemLog{
			Type:     log.TypeBalance.Uint8(),
			Date:     time.Now().Format(time.DateOnly),
			ObjectID: userInfo.Id,
			Content:  string(content),
		})
		if err != nil {
			l.Errorw("[UpdateUserBasicInfoLogic] Insert Balance Log Error:", logger.Field("err", err.Error()), logger.Field("userId", req.UserId))
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "Insert Balance Log Error")
		}
		userInfo.Balance = req.Balance
	}

	if userInfo.GiftAmount != req.GiftAmount {
		change := req.GiftAmount - userInfo.GiftAmount
		if change != 0 {
			var changeType uint16
			if userInfo.GiftAmount < req.GiftAmount {
				changeType = log.GiftTypeIncrease
			} else {
				changeType = log.GiftTypeReduce
			}
			giftLog := log.Gift{
				Type:      changeType,
				Amount:    change,
				Balance:   req.GiftAmount,
				Remark:    "Admin adjustment",
				Timestamp: time.Now().UnixMilli(),
			}
			content, _ := giftLog.Marshal()
			// Add gift amount change log
			err = l.deps.LogModel.Insert(l.ctx, &log.SystemLog{
				Type:     log.TypeGift.Uint8(),
				Date:     time.Now().Format(time.DateOnly),
				ObjectID: userInfo.Id,
				Content:  string(content),
			})
			if err != nil {
				l.Errorw("[UpdateUserBasicInfoLogic] Insert Balance Log Error:", logger.Field("err", err.Error()), logger.Field("userId", req.UserId))
				return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "Insert Balance Log Error")
			}
			userInfo.GiftAmount = req.GiftAmount
		}
	}

	if req.Commission != userInfo.Commission {

		commentLog := log.Commission{
			Type:      log.CommissionTypeAdjust,
			Amount:    req.Commission - userInfo.Commission,
			Timestamp: time.Now().UnixMilli(),
		}

		content, _ := commentLog.Marshal()
		err = l.deps.LogModel.Insert(l.ctx, &log.SystemLog{
			Type:     log.TypeCommission.Uint8(),
			Date:     time.Now().Format(time.DateOnly),
			ObjectID: userInfo.Id,
			Content:  string(content),
		})
		if err != nil {
			l.Errorw("[UpdateUserBasicInfoLogic] Insert Commission Log Error:", logger.Field("err", err.Error()), logger.Field("userId", req.UserId))
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "Insert Commission Log Error")
		}
		userInfo.Commission = req.Commission
	}
	tool.DeepCopy(userInfo, req)
	userInfo.OnlyFirstPurchase = &req.OnlyFirstPurchase
	userInfo.ReferralPercentage = req.ReferralPercentage

	if req.Password != "" {
		if userInfo.Id == 2 && isDemo {
			return errors.Wrapf(xerr.NewErrCodeMsg(503, "Demo mode does not allow modification of the admin user password"), "UpdateUserBasicInfo failed: cannot update admin user password in demo mode")
		}
		userInfo.Password = tool.EncodePassWord(req.Password)
		userInfo.Algo = "default"
	}

	err = l.deps.UserModel.Update(l.ctx, userInfo)
	if err != nil {
		l.Errorw("[UpdateUserBasicInfoLogic] Update User Error:", logger.Field("err", err.Error()), logger.Field("userId", req.UserId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "Update User Error")
	}

	return nil
}
