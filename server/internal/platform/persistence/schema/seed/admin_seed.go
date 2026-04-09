package seed

import (
	"time"

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
		var count int64
		if err := tx.Model(&user.User{}).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return nil
		}

		enable := true
		u := user.User{
			Password:  tool.EncodePassWord(password),
			IsAdmin:   &enable,
			ReferCode: uuidx.UserInviteCode(time.Now().Unix()),
		}
		if err := tx.Create(&u).Error; err != nil {
			return err
		}
		return tx.Create(&user.AuthMethods{
			UserId:         u.Id,
			AuthType:       "email",
			AuthIdentifier: email,
			Verified:       true,
		}).Error
	})
}
