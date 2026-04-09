package tool

import (
	"context"
	"net"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
)

type QueryIPLocationInput struct {
	types.QueryIPLocationRequest
}

type QueryIPLocationOutput struct {
	Body *types.QueryIPLocationResponse
}

func QueryIPLocationHandler(deps Deps) func(context.Context, *QueryIPLocationInput) (*QueryIPLocationOutput, error) {
	return func(ctx context.Context, input *QueryIPLocationInput) (*QueryIPLocationOutput, error) {
		l := NewQueryIPLocationLogic(ctx, deps)
		resp, err := l.QueryIPLocation(&input.QueryIPLocationRequest)
		if err != nil {
			return nil, err
		}
		return &QueryIPLocationOutput{Body: resp}, nil
	}
}

type QueryIPLocationLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewQueryIPLocationLogic Query IP Location
func NewQueryIPLocationLogic(ctx context.Context, deps Deps) *QueryIPLocationLogic {
	return &QueryIPLocationLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryIPLocationLogic) QueryIPLocation(req *types.QueryIPLocationRequest) (resp *types.QueryIPLocationResponse, err error) {
	if l.deps.GeoIPDB == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ERROR), " GeoIP database not configured")
	}

	ip := net.ParseIP(req.IP)
	record, err := l.deps.GeoIPDB.City(ip)
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
