package payment

import (
	"context"
	"errors"
	"fmt"

	"github.com/perfect-panel/server/internal/platform/cache"
	"github.com/perfect-panel/server/internal/platform/persistence/billing"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ Model = (*customPaymentModel)(nil)
var (
	cachePaymentIdPrefix    = "cache:payment:id:"
	cachePaymentTokenPrefix = "cache:payment:token:"
)

type (
	Model interface {
		paymentModel
		customPaymentLogicModel
	}
	paymentModel interface {
		Insert(ctx context.Context, data *Payment, tx ...*gorm.DB) error
		FindOne(ctx context.Context, id int64) (*Payment, error)
		Update(ctx context.Context, data *Payment, tx ...*gorm.DB) error
		Delete(ctx context.Context, id int64, tx ...*gorm.DB) error
		Transaction(ctx context.Context, fn func(db *gorm.DB) error) error
	}

	customPaymentModel struct {
		*defaultPaymentModel
	}
	defaultPaymentModel struct {
		cache.CachedConn
		db      *gorm.DB
		table   string
		billing *billing.Repository
	}
)

func newPaymentModel(db *gorm.DB, c *redis.Client) *defaultPaymentModel {
	return &defaultPaymentModel{
		CachedConn: cache.NewConn(db, c),
		db:         db,
		table:      "`Payment`",
		billing:    billing.NewRepository(db),
	}
}

//nolint:unused
func (m *defaultPaymentModel) batchGetCacheKeys(Payments ...*Payment) []string {
	var keys []string
	for _, payment := range Payments {
		keys = append(keys, m.getCacheKeys(payment)...)
	}
	return keys

}
func (m *defaultPaymentModel) getCacheKeys(data *Payment) []string {
	if data == nil {
		return []string{}
	}
	paymentIdKey := fmt.Sprintf("%s%v", cachePaymentIdPrefix, data.Id)
	paymentNameKey := fmt.Sprintf("%s%v", cachePaymentTokenPrefix, data.Token)
	cacheKeys := []string{
		paymentIdKey,
		paymentNameKey,
	}
	return cacheKeys
}

func (m *defaultPaymentModel) Insert(ctx context.Context, data *Payment, tx ...*gorm.DB) error {
	if m.billing.Available(firstTx(tx)...) {
		return m.ExecCtx(ctx, func(conn *gorm.DB) error {
			return m.billing.UpsertGateway(ctx, paymentToGatewayRecord(data), withConn(conn, tx)...)
		}, m.getCacheKeys(data)...)
	}
	err := m.ExecCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		return conn.Create(&data).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *defaultPaymentModel) FindOne(ctx context.Context, id int64) (*Payment, error) {
	if m.billing.Available() {
		record, err := m.billing.FindGatewayByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return gatewayRecordToPayment(record), nil
	}
	PaymentIdKey := fmt.Sprintf("%s%v", cachePaymentIdPrefix, id)
	var resp Payment
	err := m.QueryCtx(ctx, &resp, PaymentIdKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Payment{}).Where("`id` = ?", id).First(&resp).Error
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *defaultPaymentModel) Update(ctx context.Context, data *Payment, tx ...*gorm.DB) error {
	old, err := m.FindOne(ctx, data.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if m.billing.Available(firstTx(tx)...) {
		return m.ExecCtx(ctx, func(conn *gorm.DB) error {
			return m.billing.UpsertGateway(ctx, paymentToGatewayRecord(data), withConn(conn, tx)...)
		}, m.getCacheKeys(old)...)
	}
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		return conn.Save(data).Error
	}, m.getCacheKeys(old)...)
	return err
}

func (m *defaultPaymentModel) Delete(ctx context.Context, id int64, tx ...*gorm.DB) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if m.billing.Available(firstTx(tx)...) {
		return m.ExecCtx(ctx, func(conn *gorm.DB) error {
			return m.billing.DeleteGateway(ctx, id, withConn(conn, tx)...)
		}, m.getCacheKeys(data)...)
	}
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		return conn.Delete(&Payment{}, id).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *defaultPaymentModel) Transaction(ctx context.Context, fn func(db *gorm.DB) error) error {
	return m.TransactCtx(ctx, fn)
}

func paymentToGatewayRecord(data *Payment) *billing.GatewayRecord {
	if data == nil {
		return nil
	}
	return &billing.GatewayRecord{
		ID:          data.Id,
		Name:        data.Name,
		Platform:    data.Platform,
		Icon:        data.Icon,
		Domain:      data.Domain,
		Config:      data.Config,
		Description: data.Description,
		FeeMode:     data.FeeMode,
		FeePercent:  data.FeePercent,
		FeeAmount:   data.FeeAmount,
		Enable:      data.Enable,
		Token:       data.Token,
	}
}

func gatewayRecordToPayment(data *billing.GatewayRecord) *Payment {
	if data == nil {
		return nil
	}
	return &Payment{
		Id:          data.ID,
		Name:        data.Name,
		Platform:    data.Platform,
		Icon:        data.Icon,
		Domain:      data.Domain,
		Config:      data.Config,
		Description: data.Description,
		FeeMode:     data.FeeMode,
		FeePercent:  data.FeePercent,
		FeeAmount:   data.FeeAmount,
		Enable:      data.Enable,
		Token:       data.Token,
	}
}

func firstTx(tx []*gorm.DB) []*gorm.DB {
	if len(tx) == 0 || tx[0] == nil {
		return nil
	}
	return []*gorm.DB{tx[0]}
}

func withConn(conn *gorm.DB, tx []*gorm.DB) []*gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return []*gorm.DB{tx[0]}
	}
	return []*gorm.DB{conn}
}
