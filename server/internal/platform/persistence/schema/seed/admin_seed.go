package seed

import (
	"errors"
	"time"

	"github.com/perfect-panel/server/internal/platform/persistence/identity"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/uuidx"
	"gorm.io/gorm"
)

func Admin(tx *gorm.DB, email, password string) error {
	if tx == nil {
		return gorm.ErrInvalidDB
	}
	return tx.Transaction(func(tx *gorm.DB) error {
		var adminUser user.User
		err := tx.Where("is_admin = ?", true).First(&adminUser).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			enable := true
			adminUser = user.User{
				Password:  tool.EncodePassWord(password),
				IsAdmin:   &enable,
				ReferCode: uuidx.UserInviteCode(time.Now().Unix()),
			}
			if err := tx.Create(&adminUser).Error; err != nil {
				return err
			}
		case err != nil:
			return err
		}

		var adminAuth user.AuthMethods
		err = tx.Where("user_id = ? AND auth_type = ?", adminUser.Id, "email").First(&adminAuth).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			adminAuth = user.AuthMethods{
				UserId:         adminUser.Id,
				AuthType:       "email",
				AuthIdentifier: email,
				Verified:       true,
			}
			if err := tx.Create(&adminAuth).Error; err != nil {
				return err
			}
		case err != nil:
			return err
		}

		return seedNormalizedAdmin(tx, &adminUser, &adminAuth)
	})
}

func seedNormalizedAdmin(tx *gorm.DB, adminUser *user.User, adminAuth *user.AuthMethods) error {
	if adminUser == nil || adminAuth == nil {
		return nil
	}
	if !tx.Migrator().HasTable(&identity.User{}) || !tx.Migrator().HasTable(&identity.AuthIdentity{}) {
		return nil
	}

	var identityUser identity.User
	err := tx.Where("id = ?", adminUser.Id).First(&identityUser).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		identityUser = identity.User{
			ID:                    adminUser.Id,
			Password:              adminUser.Password,
			Algo:                  adminUser.Algo,
			Salt:                  adminUser.Salt,
			Avatar:                adminUser.Avatar,
			Balance:               adminUser.Balance,
			ReferCode:             adminUser.ReferCode,
			RefererID:             adminUser.RefererId,
			Commission:            adminUser.Commission,
			ReferralPercentage:    adminUser.ReferralPercentage,
			OnlyFirstPurchase:     adminUser.OnlyFirstPurchase,
			GiftAmount:            adminUser.GiftAmount,
			Enable:                adminUser.Enable,
			IsAdmin:               adminUser.IsAdmin,
			EnableBalanceNotify:   adminUser.EnableBalanceNotify,
			EnableLoginNotify:     adminUser.EnableLoginNotify,
			EnableSubscribeNotify: adminUser.EnableSubscribeNotify,
			EnableTradeNotify:     adminUser.EnableTradeNotify,
			Rules:                 adminUser.Rules,
			CreatedAt:             adminUser.CreatedAt,
			UpdatedAt:             adminUser.UpdatedAt,
			DeletedAt:             adminUser.DeletedAt,
		}
		if err := tx.Create(&identityUser).Error; err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		if err := tx.Model(&identityUser).Updates(map[string]any{
			"password":                adminUser.Password,
			"algo":                    adminUser.Algo,
			"salt":                    adminUser.Salt,
			"avatar":                  adminUser.Avatar,
			"balance":                 adminUser.Balance,
			"refer_code":              adminUser.ReferCode,
			"referer_id":              adminUser.RefererId,
			"commission":              adminUser.Commission,
			"referral_percentage":     adminUser.ReferralPercentage,
			"only_first_purchase":     adminUser.OnlyFirstPurchase,
			"gift_amount":             adminUser.GiftAmount,
			"enable":                  adminUser.Enable,
			"is_admin":                adminUser.IsAdmin,
			"enable_balance_notify":   adminUser.EnableBalanceNotify,
			"enable_login_notify":     adminUser.EnableLoginNotify,
			"enable_subscribe_notify": adminUser.EnableSubscribeNotify,
			"enable_trade_notify":     adminUser.EnableTradeNotify,
			"rules":                   adminUser.Rules,
			"deleted_at":              adminUser.DeletedAt,
		}).Error; err != nil {
			return err
		}
	}

	var authIdentity identity.AuthIdentity
	err = tx.Where("user_id = ? AND auth_type = ?", adminAuth.UserId, adminAuth.AuthType).First(&authIdentity).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		authIdentity = identity.AuthIdentity{
			ID:             adminAuth.Id,
			UserID:         adminAuth.UserId,
			AuthType:       adminAuth.AuthType,
			AuthIdentifier: adminAuth.AuthIdentifier,
			Verified:       adminAuth.Verified,
			CreatedAt:      adminAuth.CreatedAt,
			UpdatedAt:      adminAuth.UpdatedAt,
		}
		if err := tx.Create(&authIdentity).Error; err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		if err := tx.Model(&authIdentity).Updates(map[string]any{
			"auth_identifier": adminAuth.AuthIdentifier,
			"verified":        adminAuth.Verified,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}
