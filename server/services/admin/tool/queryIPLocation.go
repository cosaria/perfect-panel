package tool

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"net"
)

type QueryIPLocationInput struct {
	types.QueryIPLocationRequest
}

type QueryIPLocationOutput struct {
	Body *types.QueryIPLocationResponse
}

func QueryIPLocationHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryIPLocationInput) (*QueryIPLocationOutput, error) {
	return func(ctx context.Context, input *QueryIPLocationInput) (*QueryIPLocationOutput, error) {
		l := NewQueryIPLocationLogic(ctx, svcCtx)
		resp, err := l.QueryIPLocation(&input.QueryIPLocationRequest)
		if err != nil {
			return nil, err
		}
		return &QueryIPLocationOutput{Body: resp}, nil
	}
}

type QueryIPLocationLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewQueryIPLocationLogic Query IP Location
func NewQueryIPLocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryIPLocationLogic {
	return &QueryIPLocationLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryIPLocationLogic) QueryIPLocation(req *types.QueryIPLocationRequest) (resp *types.QueryIPLocationResponse, err error) {
	if l.svcCtx.GeoIP == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), " GeoIP database not configured")
	}

	ip := net.ParseIP(req.IP)
	record, err := l.svcCtx.GeoIP.DB.City(ip)
	if err != nil {
		l.Errorf("Failed to query IP location: %v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Failed to query IP location")
	}

	var country, region, city string
	if record.Country.Names != nil {
		country = record.Country.Names["en"]
	}
	if len(record.Subdivisions) > 0 && record.Subdivisions[0].Names != nil {
		region = record.Subdivisions[0].Names["en"]
	}
	if record.City.Names != nil {
		city = record.City.Names["en"]
	}

	return &types.QueryIPLocationResponse{
		Country: country,
		Region:  region,
		City:    city,
	}, nil
}
