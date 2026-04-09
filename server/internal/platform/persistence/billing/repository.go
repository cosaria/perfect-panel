package billing

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	billingSubscriptionRevisionName = "0004_billing_subscription"
	billingSchemaRegistryTable      = "schema_registry"
	billingRevisionStateApplied     = "applied"
)

type Repository struct {
	db *gorm.DB
}

type PaymentCallbackAttempt struct {
	PaymentID       int64
	CallbackType    string
	IdempotencyKey  string
	RawPayload      string
	AuthStatus      string
	ProcessingState string
}

type PaymentCallbackDecision struct {
	Accepted        bool
	CallbackID      string
	ProcessingState string
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) RecordPaymentCallbackAttempt(ctx context.Context, data *PaymentCallbackAttempt, tx ...*gorm.DB) (*PaymentCallbackDecision, error) {
	if data == nil {
		return nil, nil
	}
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil, nil
	}
	callback := PaymentCallback{
		PaymentID:    data.PaymentID,
		CallbackType: data.CallbackType,
		CallbackID:   data.IdempotencyKey,
		RawPayload:   data.RawPayload,
		Status:       data.ProcessingState,
	}
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "callback_id"}},
		DoNothing: true,
	}).Create(&callback)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return &PaymentCallbackDecision{
			Accepted:        true,
			CallbackID:      callback.CallbackID,
			ProcessingState: callback.Status,
		}, nil
	}
	var existing PaymentCallback
	if err := db.Where("callback_id = ?", data.IdempotencyKey).First(&existing).Error; err != nil {
		return nil, err
	}
	return &PaymentCallbackDecision{
		Accepted:        false,
		CallbackID:      existing.CallbackID,
		ProcessingState: existing.Status,
	}, nil
}

func (r *Repository) MarkPaymentCallbackProcessed(ctx context.Context, callbackID string, state string, tx ...*gorm.DB) error {
	if callbackID == "" {
		return nil
	}
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil
	}
	return db.Model(&PaymentCallback{}).
		Where("callback_id = ?", callbackID).
		Update("status", state).Error
}

func (r *Repository) Available(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil || !r.revisionApplied(db) {
		return false
	}
	return r.Installed(db)
}

func (r *Repository) Installed(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil {
		return false
	}
	return db.Migrator().HasTable(&Order{}) &&
		db.Migrator().HasTable(&OrderItem{}) &&
		db.Migrator().HasTable(&PaymentGateway{}) &&
		db.Migrator().HasTable(&PaymentGatewaySecret{})
}

func (r *Repository) FindGatewayByID(ctx context.Context, id int64, tx ...*gorm.DB) (*GatewayRecord, error) {
	return r.findGateway(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("payment_gateways.id = ?", id)
	}, tx...)
}

func (r *Repository) FindGatewayByToken(ctx context.Context, token string, tx ...*gorm.DB) (*GatewayRecord, error) {
	return r.findGateway(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("payment_gateways.token = ?", token)
	}, tx...)
}

func (r *Repository) ListGateways(ctx context.Context, tx ...*gorm.DB) ([]*GatewayRecord, error) {
	return r.listGateways(ctx, nil, tx...)
}

func (r *Repository) ListAvailableGateways(ctx context.Context, tx ...*gorm.DB) ([]*GatewayRecord, error) {
	return r.listGateways(ctx, &GatewayFilter{Enable: boolPtr(true)}, tx...)
}

func (r *Repository) ListGatewaysByPage(ctx context.Context, page, size int, filter *GatewayFilter, tx ...*gorm.DB) (int64, []*GatewayRecord, error) {
	query := r.gatewayQuery(ctx, tx...)
	query = applyGatewayFilter(query, filter)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	var rows []gatewayRow
	err := query.Order("payment_gateways.id desc").Offset((page - 1) * size).Limit(size).Scan(&rows).Error
	if err != nil {
		return 0, nil, err
	}
	return total, rowsToGatewayRecords(rows), nil
}

func (r *Repository) UpsertGateway(ctx context.Context, data *GatewayRecord, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil
	}
	gateway := PaymentGateway{
		ID:           data.ID,
		Name:         data.Name,
		Platform:     data.Platform,
		Icon:         data.Icon,
		Domain:       data.Domain,
		PublicConfig: data.Config,
		Description:  data.Description,
		FeeMode:      data.FeeMode,
		FeePercent:   data.FeePercent,
		FeeAmount:    data.FeeAmount,
		Enable:       data.Enable,
		Token:        data.Token,
	}
	secret := PaymentGatewaySecret{
		ID:               data.ID,
		PaymentGatewayID: data.ID,
		SecretConfig:     data.Config,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&gateway).Error; err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "payment_gateway_id"}},
			UpdateAll: true,
		}).Create(&secret).Error
	})
}

func (r *Repository) DeleteGateway(ctx context.Context, id int64, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("payment_gateway_id = ?", id).Delete(&PaymentGatewaySecret{}).Error; err != nil {
			return err
		}
		return tx.Delete(&PaymentGateway{}, id).Error
	})
}

func (r *Repository) FindOrderByID(ctx context.Context, id int64, tx ...*gorm.DB) (*OrderRecord, error) {
	return r.findOrder(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("orders.id = ?", id)
	}, tx...)
}

func (r *Repository) FindOrderByOrderNo(ctx context.Context, orderNo string, tx ...*gorm.DB) (*OrderRecord, error) {
	return r.findOrder(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("orders.order_no = ?", orderNo)
	}, tx...)
}

func (r *Repository) ListOrdersByPage(ctx context.Context, page, size int, status uint8, userID, subscribeID int64, search string, tx ...*gorm.DB) (int64, []*OrderRecord, error) {
	query := r.orderQuery(ctx, tx...)
	if status > 0 {
		query = query.Where("orders.status = ?", status)
	}
	if userID > 0 {
		query = query.Where("orders.user_id = ?", userID)
	}
	if subscribeID > 0 {
		query = query.Where("orders.subscribe_id = ?", subscribeID)
	}
	if search != "" {
		query = query.Where("orders.order_no LIKE ? OR orders.trade_no LIKE ? OR orders.coupon LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	var rows []orderRow
	err := query.Select(orderSelectColumns()).
		Order("orders.id desc").
		Offset((page - 1) * size).
		Limit(size).
		Scan(&rows).Error
	if err != nil {
		return 0, nil, err
	}
	return total, rowsToOrderRecords(rows), nil
}

func (r *Repository) UpsertOrder(ctx context.Context, data *OrderRecord, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil
	}
	order := Order{
		ID:               data.ID,
		ParentOrderID:    data.ParentID,
		UserID:           data.UserID,
		OrderNo:          data.OrderNo,
		Type:             data.Type,
		Quantity:         data.Quantity,
		Price:            data.Price,
		Amount:           data.Amount,
		Discount:         data.Discount,
		Coupon:           data.Coupon,
		CouponDiscount:   data.CouponDiscount,
		PaymentGatewayID: data.PaymentID,
		Method:           data.Method,
		FeeAmount:        data.FeeAmount,
		TradeNo:          data.TradeNo,
		GiftAmount:       data.GiftAmount,
		Commission:       data.Commission,
		Status:           data.Status,
		SubscribeID:      data.SubscribeID,
		SubscribeToken:   data.SubscribeToken,
		IsNew:            data.IsNew,
	}
	item := OrderItem{
		ID:          data.ID,
		OrderID:     data.ID,
		SubscribeID: data.SubscribeID,
		Quantity:    data.Quantity,
		UnitPrice:   data.Price,
		Amount:      data.Amount,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&order).Error; err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&item).Error
	})
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderNo string, status uint8, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil
	}
	return db.Model(&Order{}).Where("order_no = ?", orderNo).Update("status", status).Error
}

func (r *Repository) DeleteOrder(ctx context.Context, id int64, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", id).Delete(&OrderItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Order{}, id).Error
	})
}

type gatewayRow struct {
	ID           int64
	Name         string
	Platform     string
	Icon         string
	Domain       string
	PublicConfig string
	SecretConfig string
	Description  string
	FeeMode      uint
	FeePercent   int64
	FeeAmount    int64
	Enable       *bool
	Token        string
}

type orderRow struct {
	ID                  int64
	ParentOrderID       int64
	UserID              int64
	OrderNo             string
	Type                uint8
	Quantity            int64
	Price               int64
	Amount              int64
	Discount            int64
	Coupon              string
	CouponDiscount      int64
	PaymentGatewayID    int64
	Method              string
	FeeAmount           int64
	TradeNo             string
	GiftAmount          int64
	Commission          int64
	Status              uint8
	SubscribeID         int64
	SubscribeToken      string
	IsNew               bool
	CreatedAt           string
	UpdatedAt           string
	GatewayID           int64
	GatewayName         string
	GatewayPlatform     string
	GatewayIcon         string
	GatewayDomain       string
	GatewayPublicConfig string
	GatewaySecretConfig string
	GatewayDescription  string
	GatewayFeeMode      uint
	GatewayFeePercent   int64
	GatewayFeeAmount    int64
	GatewayEnable       *bool
	GatewayToken        string
}

func (r *Repository) gatewayQuery(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	return r.conn(ctx, tx...).Model(&PaymentGateway{}).
		Joins("LEFT JOIN payment_gateway_secrets ON payment_gateway_secrets.payment_gateway_id = payment_gateways.id")
}

func (r *Repository) orderQuery(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	return r.conn(ctx, tx...).Model(&Order{}).
		Joins("LEFT JOIN payment_gateways ON payment_gateways.id = orders.payment_gateway_id").
		Joins("LEFT JOIN payment_gateway_secrets ON payment_gateway_secrets.payment_gateway_id = payment_gateways.id")
}

func (r *Repository) findGateway(ctx context.Context, scope func(*gorm.DB) *gorm.DB, tx ...*gorm.DB) (*GatewayRecord, error) {
	var row gatewayRow
	query := r.gatewayQuery(ctx, tx...)
	if scope != nil {
		query = scope(query)
	}
	err := query.Select("payment_gateways.*, payment_gateway_secrets.secret_config").Take(&row).Error
	if err != nil {
		return nil, err
	}
	return rowToGatewayRecord(row), nil
}

func (r *Repository) listGateways(ctx context.Context, filter *GatewayFilter, tx ...*gorm.DB) ([]*GatewayRecord, error) {
	query := applyGatewayFilter(r.gatewayQuery(ctx, tx...), filter)
	var rows []gatewayRow
	err := query.Order("payment_gateways.id desc").Select("payment_gateways.*, payment_gateway_secrets.secret_config").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rowsToGatewayRecords(rows), nil
}

func (r *Repository) findOrder(ctx context.Context, scope func(*gorm.DB) *gorm.DB, tx ...*gorm.DB) (*OrderRecord, error) {
	var row orderRow
	query := r.orderQuery(ctx, tx...)
	if scope != nil {
		query = scope(query)
	}
	err := query.Select(orderSelectColumns()).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return rowToOrderRecord(row), nil
}

func applyGatewayFilter(query *gorm.DB, filter *GatewayFilter) *gorm.DB {
	if filter == nil {
		return query
	}
	if filter.Enable != nil {
		query = query.Where("payment_gateways.enable = ?", *filter.Enable)
	}
	if filter.Mark != "" {
		query = query.Where("payment_gateways.platform = ?", filter.Mark)
	}
	if filter.Search != "" {
		query = query.Where("payment_gateways.name LIKE ?", "%"+filter.Search+"%")
	}
	return query
}

func rowToGatewayRecord(row gatewayRow) *GatewayRecord {
	return &GatewayRecord{
		ID:          row.ID,
		Name:        row.Name,
		Platform:    row.Platform,
		Icon:        row.Icon,
		Domain:      row.Domain,
		Config:      coalesceConfig(row.SecretConfig, row.PublicConfig),
		Description: row.Description,
		FeeMode:     row.FeeMode,
		FeePercent:  row.FeePercent,
		FeeAmount:   row.FeeAmount,
		Enable:      row.Enable,
		Token:       row.Token,
	}
}

func rowsToGatewayRecords(rows []gatewayRow) []*GatewayRecord {
	result := make([]*GatewayRecord, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToGatewayRecord(row))
	}
	return result
}

func rowToOrderRecord(row orderRow) *OrderRecord {
	record := &OrderRecord{
		ID:             row.ID,
		ParentID:       row.ParentOrderID,
		UserID:         row.UserID,
		OrderNo:        row.OrderNo,
		Type:           row.Type,
		Quantity:       row.Quantity,
		Price:          row.Price,
		Amount:         row.Amount,
		Discount:       row.Discount,
		Coupon:         row.Coupon,
		CouponDiscount: row.CouponDiscount,
		PaymentID:      row.PaymentGatewayID,
		Method:         row.Method,
		FeeAmount:      row.FeeAmount,
		TradeNo:        row.TradeNo,
		GiftAmount:     row.GiftAmount,
		Commission:     row.Commission,
		Status:         row.Status,
		SubscribeID:    row.SubscribeID,
		SubscribeToken: row.SubscribeToken,
		IsNew:          row.IsNew,
	}
	if row.GatewayID != 0 {
		record.Payment = &GatewayRecord{
			ID:          row.GatewayID,
			Name:        row.GatewayName,
			Platform:    row.GatewayPlatform,
			Icon:        row.GatewayIcon,
			Domain:      row.GatewayDomain,
			Config:      coalesceConfig(row.GatewaySecretConfig, row.GatewayPublicConfig),
			Description: row.GatewayDescription,
			FeeMode:     row.GatewayFeeMode,
			FeePercent:  row.GatewayFeePercent,
			FeeAmount:   row.GatewayFeeAmount,
			Enable:      row.GatewayEnable,
			Token:       row.GatewayToken,
		}
	}
	return record
}

func rowsToOrderRecords(rows []orderRow) []*OrderRecord {
	result := make([]*OrderRecord, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToOrderRecord(row))
	}
	return result
}

func coalesceConfig(primary, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return primary
	}
	return fallback
}

func orderSelectColumns() string {
	return `
		orders.*,
		payment_gateways.id AS gateway_id,
		payment_gateways.name AS gateway_name,
		payment_gateways.platform AS gateway_platform,
		payment_gateways.icon AS gateway_icon,
		payment_gateways.domain AS gateway_domain,
		payment_gateways.public_config AS gateway_public_config,
		payment_gateway_secrets.secret_config AS gateway_secret_config,
		payment_gateways.description AS gateway_description,
		payment_gateways.fee_mode AS gateway_fee_mode,
		payment_gateways.fee_percent AS gateway_fee_percent,
		payment_gateways.fee_amount AS gateway_fee_amount,
		payment_gateways.enable AS gateway_enable,
		payment_gateways.token AS gateway_token
	`
}

func boolPtr(v bool) *bool {
	return &v
}

func (r *Repository) conn(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		if ctx != nil {
			return tx[0].WithContext(ctx)
		}
		return tx[0]
	}
	if r.db == nil {
		return nil
	}
	if ctx != nil {
		return r.db.WithContext(ctx)
	}
	return r.db
}

func (r *Repository) revisionApplied(db *gorm.DB) bool {
	if db == nil || !db.Migrator().HasTable(billingSchemaRegistryTable) {
		return false
	}
	var count int64
	if err := db.Table(billingSchemaRegistryTable).
		Where("id = ? AND state = ?", billingSubscriptionRevisionName, billingRevisionStateApplied).
		Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
