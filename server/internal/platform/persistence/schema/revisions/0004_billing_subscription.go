package revisions

import (
	"time"

	"github.com/perfect-panel/server/internal/platform/persistence/billing"
	legacyorder "github.com/perfect-panel/server/internal/platform/persistence/order"
	legacypayment "github.com/perfect-panel/server/internal/platform/persistence/payment"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	newsubscription "github.com/perfect-panel/server/internal/platform/persistence/subscription"
	legacysubscription "github.com/perfect-panel/server/internal/platform/persistence/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type billingSubscriptionRevision struct{}

type legacyUserSubscriptionTable struct {
	ID          int64      `gorm:"primaryKey"`
	UserID      int64      `gorm:"index:idx_user_id;not null;comment:User ID"`
	OrderID     int64      `gorm:"index:idx_order_id;not null;comment:Order ID"`
	SubscribeID int64      `gorm:"index:idx_subscribe_id;not null;comment:Subscription ID"`
	StartTime   time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null;comment:Subscription Start Time"`
	ExpireTime  time.Time  `gorm:"default:NULL;comment:Subscription Expire Time"`
	FinishedAt  *time.Time `gorm:"default:NULL;comment:Finished Time"`
	Traffic     int64      `gorm:"default:0;comment:Traffic"`
	Download    int64      `gorm:"default:0;comment:Download Traffic"`
	Upload      int64      `gorm:"default:0;comment:Upload Traffic"`
	Token       string     `gorm:"index:idx_token;unique;type:varchar(255);default:'';comment:Token"`
	UUID        string     `gorm:"type:varchar(255);unique;index:idx_uuid;default:'';comment:UUID"`
	Status      uint8      `gorm:"type:tinyint(1);default:0;comment:Subscription Status: 0: Pending 1: Active 2: Finished 3: Expired 4: Deducted 5: stopped"`
	Note        string     `gorm:"type:varchar(500);default:'';comment:User note for subscription"`
	CreatedAt   time.Time  `gorm:"<-:create;comment:Creation Time"`
	UpdatedAt   time.Time  `gorm:"comment:Update Time"`
}

func (legacyUserSubscriptionTable) TableName() string {
	return "user_subscribe"
}

func (billingSubscriptionRevision) Name() string {
	return schema.RevisionName(4, "billing_subscription")
}

func (billingSubscriptionRevision) Up(db *gorm.DB) error {
	if db.Dialector.Name() == "mysql" {
		if err := db.AutoMigrate(&legacyUserSubscriptionTable{}); err != nil {
			return err
		}
	}
	if err := db.AutoMigrate(
		&billing.PaymentGateway{},
		&billing.PaymentGatewaySecret{},
		&billing.Order{},
		&billing.OrderItem{},
		&billing.Payment{},
		&billing.PaymentCallback{},
		&billing.Refund{},
		&billing.BillingLedger{},
		&newsubscription.Subscription{},
		&newsubscription.SubscriptionPeriod{},
		&newsubscription.SubscriptionToken{},
		&newsubscription.SubscriptionUsageSnapshot{},
		&newsubscription.SubscriptionEvent{},
	); err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := backfillPaymentGateways(tx); err != nil {
			return err
		}
		if err := backfillOrders(tx); err != nil {
			return err
		}
		if err := backfillSubscriptions(tx); err != nil {
			return err
		}
		return nil
	})
}

func backfillPaymentGateways(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&legacypayment.Payment{}) {
		return nil
	}
	var legacyRows []legacypayment.Payment
	if err := tx.Find(&legacyRows).Error; err != nil {
		return err
	}
	for _, item := range legacyRows {
		gateway := billing.PaymentGateway{
			ID:           item.Id,
			Name:         item.Name,
			Platform:     item.Platform,
			Icon:         item.Icon,
			Domain:       item.Domain,
			PublicConfig: item.Config,
			Description:  item.Description,
			FeeMode:      item.FeeMode,
			FeePercent:   item.FeePercent,
			FeeAmount:    item.FeeAmount,
			Enable:       item.Enable,
			Token:        item.Token,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&gateway).Error; err != nil {
			return err
		}
		secret := billing.PaymentGatewaySecret{
			ID:               item.Id,
			PaymentGatewayID: item.Id,
			SecretConfig:     item.Config,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "payment_gateway_id"}},
			UpdateAll: true,
		}).Create(&secret).Error; err != nil {
			return err
		}
	}
	return nil
}

func backfillOrders(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&legacyorder.Order{}) {
		return nil
	}
	var legacyRows []legacyorder.Order
	if err := tx.Find(&legacyRows).Error; err != nil {
		return err
	}
	for _, item := range legacyRows {
		order := billing.Order{
			ID:               item.Id,
			ParentOrderID:    item.ParentId,
			UserID:           item.UserId,
			OrderNo:          item.OrderNo,
			Type:             item.Type,
			Quantity:         item.Quantity,
			Price:            item.Price,
			Amount:           item.Amount,
			Discount:         item.Discount,
			Coupon:           item.Coupon,
			CouponDiscount:   item.CouponDiscount,
			PaymentGatewayID: item.PaymentId,
			Method:           item.Method,
			FeeAmount:        item.FeeAmount,
			TradeNo:          item.TradeNo,
			GiftAmount:       item.GiftAmount,
			Commission:       item.Commission,
			Status:           item.Status,
			SubscribeID:      item.SubscribeId,
			SubscribeToken:   item.SubscribeToken,
			IsNew:            item.IsNew,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&order).Error; err != nil {
			return err
		}
		itemRow := billing.OrderItem{
			ID:          item.Id,
			OrderID:     item.Id,
			SubscribeID: item.SubscribeId,
			Quantity:    item.Quantity,
			UnitPrice:   item.Price,
			Amount:      item.Amount,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&itemRow).Error; err != nil {
			return err
		}
	}
	return nil
}

func backfillSubscriptions(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&legacysubscription.Subscribe{}) {
		return nil
	}
	var legacyRows []legacysubscription.Subscribe
	if err := tx.Find(&legacyRows).Error; err != nil {
		return err
	}
	for _, item := range legacyRows {
		expireTime := item.ExpireTime
		if expireTime.IsZero() {
			expireTime = time.UnixMilli(0)
		}
		subscription := newsubscription.Subscription{
			ID:          item.Id,
			UserID:      item.UserId,
			OrderID:     item.OrderId,
			SubscribeID: item.SubscribeId,
			StartTime:   item.StartTime,
			ExpireTime:  expireTime,
			FinishedAt:  item.FinishedAt,
			Traffic:     item.Traffic,
			Download:    item.Download,
			Upload:      item.Upload,
			Status:      item.Status,
			Note:        item.Note,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&subscription).Error; err != nil {
			return err
		}
		period := newsubscription.SubscriptionPeriod{
			ID:             item.Id,
			SubscriptionID: item.Id,
			StartTime:      item.StartTime,
			ExpireTime:     expireTime,
			FinishedAt:     item.FinishedAt,
			Status:         item.Status,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&period).Error; err != nil {
			return err
		}
		token := newsubscription.SubscriptionToken{
			ID:             item.Id,
			SubscriptionID: item.Id,
			Token:          item.Token,
			UUID:           item.UUID,
			IsPrimary:      true,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&token).Error; err != nil {
			return err
		}
		usage := newsubscription.SubscriptionUsageSnapshot{
			ID:             item.Id,
			SubscriptionID: item.Id,
			Traffic:        item.Traffic,
			Download:       item.Download,
			Upload:         item.Upload,
			CapturedAt:     item.UpdatedAt,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&usage).Error; err != nil {
			return err
		}
	}
	return nil
}
