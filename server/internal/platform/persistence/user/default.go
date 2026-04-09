package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/perfect-panel/server/internal/platform/cache"
	"github.com/perfect-panel/server/internal/platform/persistence/identity"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	cacheUserIdPrefix    = "cache:user:id:"
	cacheUserEmailPrefix = "cache:user:email:"
)

type userDeferredCacheKeysContextKey struct{}

var _ Model = (*customUserModel)(nil)

type (
	Model interface {
		userModel
		customUserLogicModel
	}
	userModel interface {
		Insert(ctx context.Context, data *User, tx ...*gorm.DB) error
		FindOne(ctx context.Context, id int64) (*User, error)
		Update(ctx context.Context, data *User, tx ...*gorm.DB) error
		Delete(ctx context.Context, id int64, tx ...*gorm.DB) error
		Transaction(ctx context.Context, fn func(db *gorm.DB) error) error
	}

	customUserModel struct {
		*defaultUserModel
	}
	defaultUserModel struct {
		cache.CachedConn
		db           *gorm.DB
		table        string
		identityRepo *identity.Repository
	}
)

func newUserModel(db *gorm.DB, c *redis.Client) *defaultUserModel {
	return &defaultUserModel{
		CachedConn:  cache.NewConn(db, c),
		db:          db,
		table:       "`user`",
		identityRepo: identity.NewRepository(db),
	}
}

func (m *defaultUserModel) batchGetCacheKeys(users ...*User) []string {
	var keys []string
	for _, user := range users {
		keys = append(keys, user.GetCacheKeys()...)
	}
	return keys
}

func (m *defaultUserModel) getCacheKeys(data *User) []string {
	if data == nil {
		return []string{}
	}
	return data.GetCacheKeys()
}

func (m *defaultUserModel) FindOneByEmail(ctx context.Context, email string) (*User, error) {
	var userData User
	key := fmt.Sprintf("%s%v", cacheUserEmailPrefix, email)
	err := m.QueryCtx(ctx, &userData, key, func(conn *gorm.DB, v interface{}) error {
		if m.useIdentitySchema(conn) {
			data, err := m.identityRepo.FindUserByEmail(ctx, email, conn)
			if err != nil {
				return err
			}
			typed, ok := v.(*User)
			if !ok {
				return errors.New("invalid user destination")
			}
			*typed = *m.identityUserToLegacy(data)
			return nil
		}
		var data AuthMethods
		if err := conn.Model(&AuthMethods{}).Where("`auth_type` = 'email' AND `auth_identifier` = ?", email).First(&data).Error; err != nil {
			return err
		}
		return conn.Model(&User{}).Unscoped().Where("`id` = ?", data.UserId).Preload("UserDevices").Preload("AuthMethods").First(v).Error
	})
	return &userData, err
}

func (m *defaultUserModel) Insert(ctx context.Context, data *User, tx ...*gorm.DB) error {
	err := m.ExecCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		if m.useIdentitySchema(conn) {
			row := m.legacyUserToIdentity(data)
			if err := m.identityRepo.InsertUser(ctx, row, conn); err != nil {
				return err
			}
			data.Id = row.ID
			return nil
		}
		return conn.Create(&data).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *defaultUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryCtx(ctx, &resp, userIdKey, func(conn *gorm.DB, v interface{}) error {
		if m.useIdentitySchema(conn) {
			data, err := m.identityRepo.FindUserByID(ctx, id, conn)
			if err != nil {
				return err
			}
			typed, ok := v.(*User)
			if !ok {
				return errors.New("invalid user destination")
			}
			*typed = *m.identityUserToLegacy(data)
			return nil
		}
		return conn.Model(&User{}).Unscoped().Where("`id` = ?", id).Preload("UserDevices").Preload("AuthMethods").First(&resp).Error
	})
	return &resp, err
}

func (m *defaultUserModel) Update(ctx context.Context, data *User, tx ...*gorm.DB) error {
	old, err := m.FindOne(ctx, data.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		if m.useIdentitySchema(conn) {
			return m.identityRepo.UpdateUser(ctx, m.legacyUserToIdentity(data), conn)
		}
		return conn.Save(data).Error
	}, m.getCacheKeys(old)...)
	return err
}

func (m *defaultUserModel) Delete(ctx context.Context, id int64, tx ...*gorm.DB) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	// Use batch related cache cleaning, including a cache of all relevant data
	defer func() {
		if clearErr := m.BatchClearRelatedCache(ctx, data); clearErr != nil {
			// Record cache cleaning errors, but do not block deletion operations
			logger.Errorf("failed to clear related cache for user %d: %v", id, clearErr.Error())
		}
	}()

	return m.TransactCtx(ctx, func(db *gorm.DB) error {
		if len(tx) > 0 {
			db = tx[0]
		}
		if m.useIdentitySchema(db) {
			return m.identityRepo.DeleteUser(ctx, id, db)
		}
		// Soft deletion of user information without any processing of other information (Determine whether to allow login/subscription based on the user's deletion status)
		if err := db.Model(&User{}).Where("`id` = ?", id).Delete(&User{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (m *defaultUserModel) Transaction(ctx context.Context, fn func(db *gorm.DB) error) error {
	var deferredKeys []string
	err := m.TransactCtx(ctx, func(tx *gorm.DB) error {
		tx = tx.Session(&gorm.Session{
			Context: context.WithValue(tx.Statement.Context, userDeferredCacheKeysContextKey{}, &deferredKeys),
		})
		return fn(tx)
	})
	if err != nil {
		return err
	}
	return m.DelCacheCtx(ctx, uniqueStrings(deferredKeys)...)
}

func (m *defaultUserModel) useIdentitySchema(conn *gorm.DB) bool {
	if m.identityRepo == nil {
		return false
	}
	return m.identityRepo.Available(conn)
}

func (m *defaultUserModel) identityUserToLegacy(data *identity.User) *User {
	if data == nil {
		return nil
	}

	authMethods := make([]AuthMethods, 0, len(data.AuthIdentities))
	for _, item := range data.AuthIdentities {
		authMethods = append(authMethods, AuthMethods{
			Id:             item.ID,
			UserId:         item.UserID,
			AuthType:       item.AuthType,
			AuthIdentifier: item.AuthIdentifier,
			Verified:       item.Verified,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	userDevices := make([]Device, 0, len(data.UserDevices))
	for _, item := range data.UserDevices {
		userDevices = append(userDevices, Device{
			Id:         item.ID,
			Ip:         item.IPAddress,
			UserId:     item.UserID,
			UserAgent:  item.UserAgent,
			Identifier: item.Identifier,
			Online:     item.Online,
			Enabled:    item.Enabled,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		})
	}

	return &User{
		Id:                    data.ID,
		Password:              data.Password,
		Algo:                  data.Algo,
		Salt:                  data.Salt,
		Avatar:                data.Avatar,
		Balance:               data.Balance,
		ReferCode:             data.ReferCode,
		RefererId:             data.RefererID,
		Commission:            data.Commission,
		ReferralPercentage:    data.ReferralPercentage,
		OnlyFirstPurchase:     data.OnlyFirstPurchase,
		GiftAmount:            data.GiftAmount,
		Enable:                data.Enable,
		IsAdmin:               data.IsAdmin,
		EnableBalanceNotify:   data.EnableBalanceNotify,
		EnableLoginNotify:     data.EnableLoginNotify,
		EnableSubscribeNotify: data.EnableSubscribeNotify,
		EnableTradeNotify:     data.EnableTradeNotify,
		AuthMethods:           authMethods,
		UserDevices:           userDevices,
		Rules:                 data.Rules,
		CreatedAt:             data.CreatedAt,
		UpdatedAt:             data.UpdatedAt,
		DeletedAt:             data.DeletedAt,
	}
}

func (m *defaultUserModel) legacyUserToIdentity(data *User) *identity.User {
	if data == nil {
		return nil
	}
	return &identity.User{
		ID:                    data.Id,
		Password:              data.Password,
		Algo:                  data.Algo,
		Salt:                  data.Salt,
		Avatar:                data.Avatar,
		Balance:               data.Balance,
		ReferCode:             data.ReferCode,
		RefererID:             data.RefererId,
		Commission:            data.Commission,
		ReferralPercentage:    data.ReferralPercentage,
		OnlyFirstPurchase:     data.OnlyFirstPurchase,
		GiftAmount:            data.GiftAmount,
		Enable:                data.Enable,
		IsAdmin:               data.IsAdmin,
		EnableBalanceNotify:   data.EnableBalanceNotify,
		EnableLoginNotify:     data.EnableLoginNotify,
		EnableSubscribeNotify: data.EnableSubscribeNotify,
		EnableTradeNotify:     data.EnableTradeNotify,
		Rules:                 data.Rules,
		CreatedAt:             data.CreatedAt,
		UpdatedAt:             data.UpdatedAt,
		DeletedAt:             data.DeletedAt,
	}
}

func uniqueStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
