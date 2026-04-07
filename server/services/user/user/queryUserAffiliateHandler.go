// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
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
