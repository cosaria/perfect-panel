// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/user/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryUserAffiliateListInput struct {
	types.QueryUserAffiliateListRequest
}

type QueryUserAffiliateListOutput struct {
	Body *types.QueryUserAffiliateListResponse
}

func QueryUserAffiliateListHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryUserAffiliateListInput) (*QueryUserAffiliateListOutput, error) {
	return func(ctx context.Context, input *QueryUserAffiliateListInput) (*QueryUserAffiliateListOutput, error) {
		l := user.NewQueryUserAffiliateListLogic(ctx, svcCtx)
		resp, err := l.QueryUserAffiliateList(&input.QueryUserAffiliateListRequest)
		if err != nil {
			return nil, err
		}
		return &QueryUserAffiliateListOutput{Body: resp}, nil
	}
}
