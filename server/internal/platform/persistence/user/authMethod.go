package user

import (
	"context"
	"errors"

	"github.com/perfect-panel/server/internal/platform/persistence/identity"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"gorm.io/gorm"
)

func (m *defaultUserModel) FindUserAuthMethods(ctx context.Context, userId int64) ([]*AuthMethods, error) {
	if m.useIdentitySchema(nil) {
		rows, err := m.identityRepo.FindUserAuthIdentities(ctx, userId)
		if err != nil {
			return nil, err
		}
		return m.identityAuthMethodsToLegacy(rows), nil
	}
	var data []*AuthMethods
	err := m.QueryNoCacheCtx(ctx, &data, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&AuthMethods{}).Where("user_id = ?", userId).Find(&data).Error
	})
	return data, err
}

func (m *defaultUserModel) FindUserAuthMethodByOpenID(ctx context.Context, method, openID string) (*AuthMethods, error) {
	if m.useIdentitySchema(nil) {
		data, err := m.identityRepo.FindAuthIdentityByOpenID(ctx, method, openID)
		if err != nil {
			return nil, err
		}
		return m.identityAuthMethodToLegacy(data), nil
	}
	var data AuthMethods
	err := m.QueryNoCacheCtx(ctx, &data, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&AuthMethods{}).Where("auth_type = ? AND auth_identifier = ?", method, openID).First(&data).Error
	})
	return &data, err
}

func (m *defaultUserModel) FindUserAuthMethodByPlatform(ctx context.Context, userId int64, platform string) (*AuthMethods, error) {
	if m.useIdentitySchema(nil) {
		data, err := m.identityRepo.FindAuthIdentityByPlatform(ctx, userId, platform)
		if err != nil {
			return nil, err
		}
		return m.identityAuthMethodToLegacy(data), nil
	}
	var data AuthMethods
	err := m.QueryNoCacheCtx(ctx, &data, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&AuthMethods{}).Where("user_id = ? AND auth_type = ?", userId, platform).First(&data).Error
	})
	return &data, err
}

func (m *defaultUserModel) InsertUserAuthMethods(ctx context.Context, data *AuthMethods, tx ...*gorm.DB) error {
	err := m.ExecNoCacheCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		if m.useIdentitySchema(conn) {
			row := m.legacyAuthMethodToIdentity(data)
			if insertErr := m.identityRepo.InsertUserAuthIdentity(ctx, row, conn); insertErr != nil {
				return insertErr
			}
			data.Id = row.ID
			data.CreatedAt = row.CreatedAt
			data.UpdatedAt = row.UpdatedAt
			return nil
		}
		if createErr := conn.Model(&AuthMethods{}).Create(data).Error; createErr != nil {
			return createErr
		}
		return nil
	})
	if err != nil {
		return err
	}
	var txConn *gorm.DB
	if len(tx) > 0 {
		txConn = tx[0]
	}
	return m.clearModelsCacheWithTx(ctx, txConn, data)
}

func (m *defaultUserModel) UpdateUserAuthMethods(ctx context.Context, data *AuthMethods, tx ...*gorm.DB) error {
	var old *AuthMethods
	err := m.ExecNoCacheCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		existing, findErr := m.findUserAuthMethodWithConn(ctx, conn, data.Id, data.UserId, data.AuthType)
		if findErr != nil && !errors.Is(findErr, gorm.ErrRecordNotFound) {
			return findErr
		}
		old = existing
		if m.useIdentitySchema(conn) {
			if updateErr := m.identityRepo.UpdateUserAuthIdentity(ctx, m.legacyAuthMethodToIdentity(data), conn); updateErr != nil {
				return updateErr
			}
			return nil
		}
		if saveErr := conn.Model(&AuthMethods{}).Where("user_id = ? AND auth_type = ?", data.UserId, data.AuthType).Save(data).Error; saveErr != nil {
			return saveErr
		}
		return nil
	})
	if err != nil {
		return err
	}
	var txConn *gorm.DB
	if len(tx) > 0 {
		txConn = tx[0]
	}
	return m.clearModelsCacheWithTx(ctx, txConn, old, data)
}

func (m *defaultUserModel) DeleteUserAuthMethods(ctx context.Context, userId int64, platform string, tx ...*gorm.DB) error {
	var old *AuthMethods
	err := m.ExecNoCacheCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		existing, findErr := m.findUserAuthMethodWithConn(ctx, conn, 0, userId, platform)
		if findErr != nil && !errors.Is(findErr, gorm.ErrRecordNotFound) {
			return findErr
		}
		old = existing
		if m.useIdentitySchema(conn) {
			return m.identityRepo.DeleteUserAuthIdentity(ctx, userId, platform, conn)
		}
		return conn.Model(&AuthMethods{}).Where("user_id = ? AND auth_type = ?", userId, platform).Delete(&AuthMethods{}).Error
	})
	if err != nil {
		return err
	}
	var txConn *gorm.DB
	if len(tx) > 0 {
		txConn = tx[0]
	}
	if err = m.clearModelsCacheWithTx(context.Background(), txConn, old); err != nil {
		logger.Errorf("[UserModel] clear auth method cache failed: %v", err.Error())
	}
	return nil
}

func (m *defaultUserModel) FindUserAuthMethodByUserId(ctx context.Context, method string, userId int64) (*AuthMethods, error) {
	if m.useIdentitySchema(nil) {
		data, err := m.identityRepo.FindAuthIdentityByUserID(ctx, method, userId)
		if err != nil {
			return nil, err
		}
		return m.identityAuthMethodToLegacy(data), nil
	}
	var data AuthMethods
	err := m.QueryNoCacheCtx(ctx, &data, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&AuthMethods{}).Where("auth_type = ? AND user_id = ?", method, userId).First(&data).Error
	})
	return &data, err
}

func (m *defaultUserModel) identityAuthMethodsToLegacy(items []*identity.AuthIdentity) []*AuthMethods {
	result := make([]*AuthMethods, 0, len(items))
	for _, item := range items {
		result = append(result, m.identityAuthMethodToLegacy(item))
	}
	return result
}

func (m *defaultUserModel) identityAuthMethodToLegacy(item *identity.AuthIdentity) *AuthMethods {
	if item == nil {
		return nil
	}
	return &AuthMethods{
		Id:             item.ID,
		UserId:         item.UserID,
		AuthType:       item.AuthType,
		AuthIdentifier: item.AuthIdentifier,
		Verified:       item.Verified,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func (m *defaultUserModel) legacyAuthMethodToIdentity(item *AuthMethods) *identity.AuthIdentity {
	if item == nil {
		return nil
	}
	return &identity.AuthIdentity{
		ID:             item.Id,
		UserID:         item.UserId,
		AuthType:       item.AuthType,
		AuthIdentifier: item.AuthIdentifier,
		Verified:       item.Verified,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func (m *defaultUserModel) findUserAuthMethodWithConn(ctx context.Context, conn *gorm.DB, id, userId int64, authType string) (*AuthMethods, error) {
	if m.useIdentitySchema(conn) {
		var data identity.AuthIdentity
		query := conn.WithContext(ctx).Model(&identity.AuthIdentity{})
		if id != 0 {
			if err := query.Where("id = ?", id).First(&data).Error; err != nil {
				return nil, err
			}
			return m.identityAuthMethodToLegacy(&data), nil
		}
		if err := query.Where("user_id = ? AND auth_type = ?", userId, authType).First(&data).Error; err != nil {
			return nil, err
		}
		return m.identityAuthMethodToLegacy(&data), nil
	}

	var data AuthMethods
	query := conn.WithContext(ctx).Model(&AuthMethods{})
	if id != 0 {
		if err := query.Where("id = ?", id).First(&data).Error; err != nil {
			return nil, err
		}
		return &data, nil
	}
	if err := query.Where("user_id = ? AND auth_type = ?", userId, authType).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
