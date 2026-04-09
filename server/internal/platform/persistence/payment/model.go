package payment

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/persistence/billing"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type customPaymentLogicModel interface {
	FindOneByPaymentToken(ctx context.Context, token string) (*Payment, error)
	FindAll(ctx context.Context) ([]*Payment, error)
	FindListByPage(ctx context.Context, page, size int, req *Filter) (int64, []*Payment, error)
	FindAvailableMethods(ctx context.Context) ([]*Payment, error)
}

// NewModel returns a model for the database table.
func NewModel(conn *gorm.DB, c *redis.Client) Model {
	return &customPaymentModel{
		defaultPaymentModel: newPaymentModel(conn, c),
	}
}

func (m *customPaymentModel) FindOneByPaymentToken(ctx context.Context, token string) (*Payment, error) {
	if m.billing.Available() {
		record, err := m.billing.FindGatewayByToken(ctx, token)
		if err != nil {
			return nil, err
		}
		return gatewayRecordToPayment(record), nil
	}
	var resp *Payment
	key := cachePaymentTokenPrefix + token
	err := m.QueryCtx(ctx, &resp, key, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Payment{}).Where("token = ?", token).First(v).Error
	})
	return resp, err
}

func (m *customPaymentModel) FindAll(ctx context.Context) ([]*Payment, error) {
	if m.billing.Available() {
		rows, err := m.billing.ListGateways(ctx)
		if err != nil {
			return nil, err
		}
		result := make([]*Payment, 0, len(rows))
		for _, row := range rows {
			result = append(result, gatewayRecordToPayment(row))
		}
		return result, nil
	}
	var resp []*Payment
	err := m.QueryNoCacheCtx(ctx, &resp, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Payment{}).Find(v).Error
	})
	return resp, err
}

func (m *customPaymentModel) FindAvailableMethods(ctx context.Context) ([]*Payment, error) {
	if m.billing.Available() {
		rows, err := m.billing.ListAvailableGateways(ctx)
		if err != nil {
			return nil, err
		}
		result := make([]*Payment, 0, len(rows))
		for _, row := range rows {
			result = append(result, gatewayRecordToPayment(row))
		}
		return result, nil
	}
	var resp []*Payment
	err := m.QueryNoCacheCtx(ctx, &resp, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Payment{}).Where("enable = ?", true).Find(v).Error
	})
	return resp, err
}

func (m *customPaymentModel) FindListByPage(ctx context.Context, page, size int, req *Filter) (int64, []*Payment, error) {
	if m.billing.Available() {
		total, rows, err := m.billing.ListGatewaysByPage(ctx, page, size, &billing.GatewayFilter{
			Search: filterSearch(req),
			Mark:   filterMark(req),
			Enable: filterEnable(req),
		})
		if err != nil {
			return 0, nil, err
		}
		result := make([]*Payment, 0, len(rows))
		for _, row := range rows {
			result = append(result, gatewayRecordToPayment(row))
		}
		return total, result, nil
	}
	var resp []*Payment
	var total int64
	err := m.QueryNoCacheCtx(ctx, &resp, func(conn *gorm.DB, v interface{}) error {
		conn = conn.Model(&Payment{})
		if req != nil {
			if req.Enable != nil {
				conn = conn.Where("`enable` = ?", *req.Enable)
			}
			if req.Mark != "" {
				conn = conn.Where("`mark` = ?", req.Mark)
			}
			if req.Search != "" {
				conn = conn.Where("`name` LIKE ?", "%"+req.Search+"%")
			}
		}

		return conn.Count(&total).Offset((page - 1) * size).Limit(size).Find(v).Error
	})
	return total, resp, err
}

func filterSearch(req *Filter) string {
	if req == nil {
		return ""
	}
	return req.Search
}

func filterMark(req *Filter) string {
	if req == nil {
		return ""
	}
	return req.Mark
}

func filterEnable(req *Filter) *bool {
	if req == nil {
		return nil
	}
	return req.Enable
}
