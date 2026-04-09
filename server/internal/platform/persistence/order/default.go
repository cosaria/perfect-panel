package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/perfect-panel/server/internal/platform/cache"
	"github.com/perfect-panel/server/internal/platform/persistence/billing"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ Model = (*customOrderModel)(nil)
var (
	cacheOrderIdPrefix = "cache:order:id:"
	cacheOrderNoPrefix = "cache:order:no:"
)

type (
	Model interface {
		orderModel
		customOrderLogicModel
	}
	orderModel interface {
		Insert(ctx context.Context, data *Order, tx ...*gorm.DB) error
		FindOne(ctx context.Context, id int64) (*Order, error)
		FindOneByOrderNo(ctx context.Context, orderNo string) (*Order, error)
		Update(ctx context.Context, data *Order, tx ...*gorm.DB) error
		Delete(ctx context.Context, id int64, tx ...*gorm.DB) error
		Transaction(ctx context.Context, fn func(db *gorm.DB) error) error
	}

	customOrderModel struct {
		*defaultOrderModel
	}
	defaultOrderModel struct {
		cache.CachedConn
		db      *gorm.DB
		table   string
		billing *billing.Repository
	}
)

func newOrderModel(db *gorm.DB, c *redis.Client) *defaultOrderModel {
	return &defaultOrderModel{
		CachedConn: cache.NewConn(db, c),
		db:         db,
		table:      "`order`",
		billing:    billing.NewRepository(db),
	}
}

//nolint:unused
func (m *defaultOrderModel) batchGetCacheKeys(Orders ...*Order) []string {
	var keys []string
	for _, order := range Orders {
		keys = append(keys, m.getCacheKeys(order)...)
	}
	return keys

}
func (m *defaultOrderModel) getCacheKeys(data *Order) []string {
	if data == nil {
		return []string{}
	}
	orderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, data.Id)
	orderNoKey := fmt.Sprintf("%s%v", cacheOrderNoPrefix, data.OrderNo)
	cacheKeys := []string{
		orderIdKey,
		orderNoKey,
	}
	return cacheKeys
}

func (m *defaultOrderModel) Insert(ctx context.Context, data *Order, tx ...*gorm.DB) error {
	if m.billing.Available(firstOrderTx(tx)...) {
		return m.ExecCtx(ctx, func(conn *gorm.DB) error {
			return m.billing.UpsertOrder(ctx, orderToBillingRecord(data), orderWithConn(conn, tx)...)
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

func (m *defaultOrderModel) FindOne(ctx context.Context, id int64) (*Order, error) {
	if m.billing.Available() {
		record, err := m.billing.FindOrderByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return billingRecordToOrder(record), nil
	}
	OrderIdKey := fmt.Sprintf("%s%v", cacheOrderIdPrefix, id)
	var resp Order
	err := m.QueryCtx(ctx, &resp, OrderIdKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Order{}).Where("`id` = ?", id).First(&resp).Error
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *defaultOrderModel) FindOneByOrderNo(ctx context.Context, orderNo string) (*Order, error) {
	if m.billing.Available() {
		record, err := m.billing.FindOrderByOrderNo(ctx, orderNo)
		if err != nil {
			return nil, err
		}
		return billingRecordToOrder(record), nil
	}
	OrderNoKey := fmt.Sprintf("%s%v", cacheOrderNoPrefix, orderNo)
	var resp Order
	err := m.QueryCtx(ctx, &resp, OrderNoKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Order{}).Where("`order_no` = ?", orderNo).First(&resp).Error
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *defaultOrderModel) Update(ctx context.Context, data *Order, tx ...*gorm.DB) error {
	old, err := m.FindOne(ctx, data.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if m.billing.Available(firstOrderTx(tx)...) {
		return m.ExecCtx(ctx, func(conn *gorm.DB) error {
			return m.billing.UpsertOrder(ctx, orderToBillingRecord(data), orderWithConn(conn, tx)...)
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

func (m *defaultOrderModel) Delete(ctx context.Context, id int64, tx ...*gorm.DB) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if m.billing.Available(firstOrderTx(tx)...) {
		return m.ExecCtx(ctx, func(conn *gorm.DB) error {
			return m.billing.DeleteOrder(ctx, id, orderWithConn(conn, tx)...)
		}, m.getCacheKeys(data)...)
	}
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		return conn.Delete(&Order{}, id).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *defaultOrderModel) Transaction(ctx context.Context, fn func(db *gorm.DB) error) error {
	return m.TransactCtx(ctx, fn)
}

func orderToBillingRecord(data *Order) *billing.OrderRecord {
	if data == nil {
		return nil
	}
	return &billing.OrderRecord{
		ID:             data.Id,
		ParentID:       data.ParentId,
		UserID:         data.UserId,
		OrderNo:        data.OrderNo,
		Type:           data.Type,
		Quantity:       data.Quantity,
		Price:          data.Price,
		Amount:         data.Amount,
		Discount:       data.Discount,
		Coupon:         data.Coupon,
		CouponDiscount: data.CouponDiscount,
		PaymentID:      data.PaymentId,
		Method:         data.Method,
		FeeAmount:      data.FeeAmount,
		TradeNo:        data.TradeNo,
		GiftAmount:     data.GiftAmount,
		Commission:     data.Commission,
		Status:         data.Status,
		SubscribeID:    data.SubscribeId,
		SubscribeToken: data.SubscribeToken,
		IsNew:          data.IsNew,
	}
}

func billingRecordToOrder(data *billing.OrderRecord) *Order {
	if data == nil {
		return nil
	}
	return &Order{
		Id:             data.ID,
		ParentId:       data.ParentID,
		UserId:         data.UserID,
		OrderNo:        data.OrderNo,
		Type:           data.Type,
		Quantity:       data.Quantity,
		Price:          data.Price,
		Amount:         data.Amount,
		GiftAmount:     data.GiftAmount,
		Discount:       data.Discount,
		Coupon:         data.Coupon,
		CouponDiscount: data.CouponDiscount,
		Commission:     data.Commission,
		PaymentId:      data.PaymentID,
		Method:         data.Method,
		FeeAmount:      data.FeeAmount,
		TradeNo:        data.TradeNo,
		Status:         data.Status,
		SubscribeId:    data.SubscribeID,
		SubscribeToken: data.SubscribeToken,
		IsNew:          data.IsNew,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}
}

func firstOrderTx(tx []*gorm.DB) []*gorm.DB {
	if len(tx) == 0 || tx[0] == nil {
		return nil
	}
	return []*gorm.DB{tx[0]}
}

func orderWithConn(conn *gorm.DB, tx []*gorm.DB) []*gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return []*gorm.DB{tx[0]}
	}
	return []*gorm.DB{conn}
}
