package revisions

import (
	"github.com/perfect-panel/server/internal/platform/persistence/auth"
	"github.com/perfect-panel/server/internal/platform/persistence/identity"
	"github.com/perfect-panel/server/internal/platform/persistence/schema"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type identitySystemRevision struct{}

func (identitySystemRevision) Name() string {
	return schema.RevisionName(2, "identity_system")
}

func (identitySystemRevision) Up(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&identity.User{},
		&identity.AuthIdentity{},
		&identity.UserSession{},
		&identity.UserDevice{},
		&identity.VerificationToken{},
		&identity.VerificationDelivery{},
		&identity.SecurityEvent{},
		&system.AuthProvider{},
		&system.AuthProviderConfig{},
		&system.VerificationPolicy{},
		&system.Setting{},
	); err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := backfillIdentityUsers(tx); err != nil {
			return err
		}
		if err := backfillIdentityAuthMethods(tx); err != nil {
			return err
		}
		if err := backfillIdentityDevices(tx); err != nil {
			return err
		}
		if err := backfillAuthProviders(tx); err != nil {
			return err
		}
		if err := backfillSystemSettings(tx); err != nil {
			return err
		}
		return nil
	})
}

func backfillIdentityUsers(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&user.User{}) {
		return nil
	}
	var legacyUsers []user.User
	if err := tx.Unscoped().Find(&legacyUsers).Error; err != nil {
		return err
	}
	for _, item := range legacyUsers {
		row := identity.User{
			ID:                    item.Id,
			Password:              item.Password,
			Algo:                  item.Algo,
			Salt:                  item.Salt,
			Avatar:                item.Avatar,
			Balance:               item.Balance,
			ReferCode:             item.ReferCode,
			RefererID:             item.RefererId,
			Commission:            item.Commission,
			ReferralPercentage:    item.ReferralPercentage,
			OnlyFirstPurchase:     item.OnlyFirstPurchase,
			GiftAmount:            item.GiftAmount,
			Enable:                item.Enable,
			IsAdmin:               item.IsAdmin,
			EnableBalanceNotify:   item.EnableBalanceNotify,
			EnableLoginNotify:     item.EnableLoginNotify,
			EnableSubscribeNotify: item.EnableSubscribeNotify,
			EnableTradeNotify:     item.EnableTradeNotify,
			Rules:                 item.Rules,
			CreatedAt:             item.CreatedAt,
			UpdatedAt:             item.UpdatedAt,
			DeletedAt:             item.DeletedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func backfillIdentityAuthMethods(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&user.AuthMethods{}) {
		return nil
	}
	var legacyAuths []user.AuthMethods
	if err := tx.Find(&legacyAuths).Error; err != nil {
		return err
	}
	for _, item := range legacyAuths {
		row := identity.AuthIdentity{
			ID:             item.Id,
			UserID:         item.UserId,
			AuthType:       item.AuthType,
			AuthIdentifier: item.AuthIdentifier,
			Verified:       item.Verified,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func backfillIdentityDevices(tx *gorm.DB) error {
	if !tx.Migrator().HasTable("user_device") {
		return nil
	}
	var legacyDevices []user.Device
	if err := tx.Table("user_device").Find(&legacyDevices).Error; err != nil {
		return err
	}
	for _, item := range legacyDevices {
		row := identity.UserDevice{
			ID:         item.Id,
			IPAddress:  item.Ip,
			UserID:     item.UserId,
			UserAgent:  item.UserAgent,
			Identifier: item.Identifier,
			Online:     item.Online,
			Enabled:    item.Enabled,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func backfillAuthProviders(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&auth.Auth{}) {
		return nil
	}
	var legacyProviders []auth.Auth
	if err := tx.Find(&legacyProviders).Error; err != nil {
		return err
	}
	for _, item := range legacyProviders {
		provider := system.AuthProvider{
			ID:      item.Id,
			Method:  item.Method,
			Enabled: item.Enabled,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&provider).Error; err != nil {
			return err
		}
		config := system.AuthProviderConfig{
			ProviderID: provider.ID,
			Config:     item.Config,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "provider_id"}},
			UpdateAll: true,
		}).Create(&config).Error; err != nil {
			return err
		}
	}
	return nil
}

func backfillSystemSettings(tx *gorm.DB) error {
	if !tx.Migrator().HasTable(&system.System{}) {
		return nil
	}
	var legacyRows []system.System
	if err := tx.Find(&legacyRows).Error; err != nil {
		return err
	}
	for _, item := range legacyRows {
		switch item.Category {
		case "verify", "verify_code":
			row := system.VerificationPolicy{
				ID:       item.Id,
				Category: item.Category,
				Key:      item.Key,
				Value:    item.Value,
				Type:     item.Type,
				Desc:     item.Desc,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&row).Error; err != nil {
				return err
			}
		default:
			row := system.Setting{
				ID:       item.Id,
				Category: item.Category,
				Key:      item.Key,
				Value:    item.Value,
				Type:     item.Type,
				Desc:     item.Desc,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&row).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
