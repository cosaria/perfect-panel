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

func (billingSubscriptionRevision) Name() string {
	return schema.RevisionName(4, "billing_subscription")
}

func (billingSubscriptionRevision) Up(db *gorm.DB) error {
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
