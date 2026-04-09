package server

import (
	"context"
	"strings"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type QueryNodeTagOutput struct {
	Body *types.QueryNodeTagResponse
}

func QueryNodeTagHandler(deps Deps) func(context.Context, *struct{}) (*QueryNodeTagOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryNodeTagOutput, error) {
		l := NewQueryNodeTagLogic(ctx, deps)
		resp, err := l.QueryNodeTag()
		if err != nil {
			return nil, err
		}
		return &QueryNodeTagOutput{Body: resp}, nil
	}
}

type QueryNodeTagLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewQueryNodeTagLogic Query all node tags
func NewQueryNodeTagLogic(ctx context.Context, deps Deps) *QueryNodeTagLogic {
	return &QueryNodeTagLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryNodeTagLogic) QueryNodeTag() (resp *types.QueryNodeTagResponse, err error) {

	var nodes []*node.Node
	if err = l.deps.DB.WithContext(l.ctx).Model(&node.Node{}).Find(&nodes).Error; err != nil {
		l.Errorw("[QueryNodeTag] Query Database Error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "[QueryNodeTag] Query Database Error")
	}
	var tags []string
	for _, item := range nodes {
		tags = append(tags, strings.Split(item.Tags, ",")...)
	}

	return &types.QueryNodeTagResponse{
		Tags: tool.RemoveDuplicateElements(tags...),
	}, nil
}
