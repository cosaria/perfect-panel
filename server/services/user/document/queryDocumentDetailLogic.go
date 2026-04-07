package document

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type QueryDocumentDetailLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get document detail
func NewQueryDocumentDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryDocumentDetailLogic {
	return &QueryDocumentDetailLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryDocumentDetailLogic) QueryDocumentDetail(req *types.QueryDocumentDetailRequest) (resp *types.Document, err error) {
	// find document
	data, err := l.svcCtx.DocumentModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[QueryDocumentDetailLogic] FindOne error", logger.Field("id", req.Id), logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOne error: %s", err.Error())
	}
	resp = &types.Document{}
	tool.DeepCopy(resp, data)
	return
}
