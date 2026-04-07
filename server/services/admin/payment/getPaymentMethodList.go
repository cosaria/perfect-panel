package payment

import (
	"context"
	"encoding/json"
	"github.com/perfect-panel/server/models/payment"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	paymentPlatform "github.com/perfect-panel/server/modules/payment"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetPaymentMethodListInput struct {
	Body types.GetPaymentMethodListRequest
}

type GetPaymentMethodListOutput struct {
	Body *types.GetPaymentMethodListResponse
}

func GetPaymentMethodListHandler(deps Deps) func(context.Context, *GetPaymentMethodListInput) (*GetPaymentMethodListOutput, error) {
	return func(ctx context.Context, input *GetPaymentMethodListInput) (*GetPaymentMethodListOutput, error) {
		l := NewGetPaymentMethodListLogic(ctx, deps)
		resp, err := l.GetPaymentMethodList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetPaymentMethodListOutput{Body: resp}, nil
	}
}

type GetPaymentMethodListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetPaymentMethodListLogic Get Payment Method List
func NewGetPaymentMethodListLogic(ctx context.Context, deps Deps) *GetPaymentMethodListLogic {
	return &GetPaymentMethodListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetPaymentMethodListLogic) GetPaymentMethodList(req *types.GetPaymentMethodListRequest) (resp *types.GetPaymentMethodListResponse, err error) {
	total, list, err := l.deps.PaymentModel.FindListByPage(l.ctx, req.Page, req.Size, &payment.Filter{
		Search: req.Search,
		Mark:   req.Platform,
		Enable: req.Enable,
	})
	if err != nil {
		l.Errorw("find payment method list error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find payment method list error: %s", err.Error())
	}
	resp = &types.GetPaymentMethodListResponse{
		Total: total,
		List:  make([]types.PaymentMethodDetail, len(list)),
	}

	for i, v := range list {
		config := make(map[string]interface{})
		_ = json.Unmarshal([]byte(v.Config), &config)
		notifyUrl := ""

		if paymentPlatform.ParsePlatform(v.Platform) != paymentPlatform.Balance {
			if v.Domain != "" {
				notifyUrl = v.Domain + "/api/v1/notify/" + v.Platform + "/" + v.Token
			} else {
				notifyUrl = "https://" + l.deps.Config.Host + "/api/v1/notify/" + v.Platform + "/" + v.Token
			}
		}
		resp.List[i] = types.PaymentMethodDetail{
			Id:          v.Id,
			Name:        v.Name,
			Platform:    v.Platform,
			Icon:        v.Icon,
			Domain:      v.Domain,
			Config:      config,
			FeeMode:     v.FeeMode,
			FeePercent:  v.FeePercent,
			FeeAmount:   v.FeeAmount,
			Enable:      *v.Enable,
			NotifyURL:   notifyUrl,
			Description: v.Description,
		}
	}
	return
}
