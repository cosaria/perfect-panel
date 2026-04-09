package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/perfect-panel/server/internal/platform/persistence/identity"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"gorm.io/gorm"
)

func (m *customUserModel) FindOneDevice(ctx context.Context, id int64) (*Device, error) {
	if m.useIdentitySchema(nil) {
		data, err := m.identityRepo.FindUserDevice(ctx, id)
		if err != nil {
			return nil, err
		}
		return m.identityDeviceToLegacy(data), nil
	}
	deviceIdKey := fmt.Sprintf("%s%v", cacheUserDeviceIdPrefix, id)
	var resp Device
	err := m.QueryCtx(ctx, &resp, deviceIdKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Device{}).Where("`id` = ?", id).First(&resp).Error
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customUserModel) FindOneDeviceByIdentifier(ctx context.Context, id string) (*Device, error) {
	if m.useIdentitySchema(nil) {
		data, err := m.identityRepo.FindUserDeviceByIdentifier(ctx, id)
		if err != nil {
			return nil, err
		}
		return m.identityDeviceToLegacy(data), nil
	}
	deviceIdKey := fmt.Sprintf("%s%v", cacheUserDeviceNumberPrefix, id)
	var resp Device
	err := m.QueryCtx(ctx, &resp, deviceIdKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Device{}).Where("`identifier` = ?", id).First(&resp).Error
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// QueryDevicePageList  returns a list of records that meet the conditions.
func (m *customUserModel) QueryDevicePageList(ctx context.Context, userId, subscribeId int64, page, size int) ([]*Device, int64, error) {
	var list []*Device
	var total int64
	err := m.QueryNoCacheCtx(ctx, &list, func(conn *gorm.DB, v interface{}) error {
		if m.useIdentitySchema(conn) {
			query := conn.Model(&identity.UserDevice{}).Where("user_devices.user_id = ?", userId)
			if subscribeId > 0 {
				query = query.Joins("JOIN user_subscribe ON user_subscribe.user_id = user_devices.user_id").
					Where("user_subscribe.subscribe_id = ? AND user_subscribe.status IN ?", subscribeId, []int64{0, 1})
			}
			var rows []*identity.UserDevice
			if err := query.Count(&total).Limit(size).Offset((page - 1) * size).Find(&rows).Error; err != nil {
				return err
			}
			list = m.identityDevicesToLegacy(rows)
			return nil
		}
		return conn.Model(&Device{}).Where("`user_id` = ? and `subscribe_id` = ?", userId, subscribeId).Count(&total).Limit(size).Offset((page - 1) * size).Find(&list).Error
	})
	return list, total, err
}

// QueryDeviceList  returns a list of records that meet the conditions.
func (m *customUserModel) QueryDeviceList(ctx context.Context, userId int64) ([]*Device, int64, error) {
	var list []*Device
	var total int64
	err := m.QueryNoCacheCtx(ctx, &list, func(conn *gorm.DB, v interface{}) error {
		if m.useIdentitySchema(conn) {
			var rows []*identity.UserDevice
			if err := conn.Model(&identity.UserDevice{}).Where("user_id = ?", userId).Count(&total).Find(&rows).Error; err != nil {
				return err
			}
			list = m.identityDevicesToLegacy(rows)
			return nil
		}
		return conn.Model(&Device{}).Where("`user_id` = ?", userId).Count(&total).Find(&list).Error
	})
	return list, total, err
}

func (m *customUserModel) UpdateDevice(ctx context.Context, data *Device, tx ...*gorm.DB) error {
	var old *Device
	err := m.ExecNoCacheCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		existing, findErr := m.findDeviceWithConn(ctx, conn, data.Id)
		if findErr != nil {
			return findErr
		}
		old = existing
		if m.useIdentitySchema(conn) {
			return m.identityRepo.UpdateUserDevice(ctx, m.legacyDeviceToIdentity(data), conn)
		}
		return conn.Save(data).Error
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

func (m *customUserModel) DeleteDevice(ctx context.Context, id int64, tx ...*gorm.DB) error {
	var data *Device
	err := m.ExecNoCacheCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		existing, findErr := m.findDeviceWithConn(ctx, conn, id)
		if findErr != nil {
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return nil
			}
			return findErr
		}
		data = existing
		if m.useIdentitySchema(conn) {
			return m.identityRepo.DeleteUserDevice(ctx, id, conn)
		}
		return conn.Delete(&Device{}, id).Error
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

func (m *customUserModel) InsertDevice(ctx context.Context, data *Device, tx ...*gorm.DB) error {
	err := m.ExecNoCacheCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		if m.useIdentitySchema(conn) {
			row := m.legacyDeviceToIdentity(data)
			if err := m.identityRepo.InsertUserDevice(ctx, row, conn); err != nil {
				return err
			}
			data.Id = row.ID
			return nil
		}
		return conn.Create(data).Error
	})
	if err != nil {
		return err
	}
	var txConn *gorm.DB
	if len(tx) > 0 {
		txConn = tx[0]
	}
	if err := m.clearModelsCacheWithTx(ctx, txConn, data); err != nil {
		logger.Errorf("failed to clear device cache: %v", err)
	}
	return nil
}

func (m *customUserModel) identityDeviceToLegacy(item *identity.UserDevice) *Device {
	if item == nil {
		return nil
	}
	return &Device{
		Id:         item.ID,
		Ip:         item.IPAddress,
		UserId:     item.UserID,
		UserAgent:  item.UserAgent,
		Identifier: item.Identifier,
		Online:     item.Online,
		Enabled:    item.Enabled,
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
	}
}

func (m *customUserModel) identityDevicesToLegacy(items []*identity.UserDevice) []*Device {
	result := make([]*Device, 0, len(items))
	for _, item := range items {
		result = append(result, m.identityDeviceToLegacy(item))
	}
	return result
}

func (m *customUserModel) legacyDeviceToIdentity(item *Device) *identity.UserDevice {
	if item == nil {
		return nil
	}
	return &identity.UserDevice{
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
}

func (m *customUserModel) findDeviceWithConn(ctx context.Context, conn *gorm.DB, id int64) (*Device, error) {
	if m.useIdentitySchema(conn) {
		data, err := m.identityRepo.FindUserDevice(ctx, id, conn)
		if err != nil {
			return nil, err
		}
		return m.identityDeviceToLegacy(data), nil
	}

	var data Device
	if err := conn.WithContext(ctx).Model(&Device{}).Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
