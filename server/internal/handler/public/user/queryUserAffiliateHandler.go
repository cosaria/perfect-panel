// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type QueryUserAffiliateOutput struct {
	Body *types.QueryUserAffiliateCountResponse
}

func QueryUserAffiliateHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryUserAffiliateOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserAffiliateOutput, error) {
		l := user.NewQueryUserAffiliateLogic(ctx, svcCtx)
		resp, err := l.QueryUserAffiliate()
		if err != nil {
			return nil, err
		}
		return &QueryUserAffiliateOutput{Body: resp}, nil
	}
}
