package user

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/perfect-panel/server/internal/platform/persistence/identity"
	"github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modelsubscription "github.com/perfect-panel/server/internal/platform/persistence/subscription"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"gorm.io/gorm"
)

type subscriptionAssignmentSyncer interface {
	SyncUserSubscription(ctx context.Context, userSubscribeId, subscribeId int64, status uint8, tx ...*gorm.DB) error
	DeleteUserSubscription(ctx context.Context, userSubscribeId int64, tx ...*gorm.DB) error
}

func (m *defaultUserModel) UpdateUserSubscribeCache(ctx context.Context, data *Subscribe) error {
	return m.ClearSubscribeCacheByModels(ctx, data)
}

// QueryActiveSubscriptions returns the number of active subscriptions.
func (m *defaultUserModel) QueryActiveSubscriptions(ctx context.Context, subscribeId ...int64) (map[int64]int64, error) {
	if m.subscriptionRepo.Available() {
		return m.subscriptionRepo.CountActiveBySubscribeID(ctx, subscribeId)
	}
	type SubscriptionCount struct {
		SubscribeId int64
		Total       int64
	}
	var result []SubscriptionCount
	err := m.QueryNoCacheCtx(ctx, &result, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Subscribe{}).
			Where("subscribe_id IN ? AND `status` IN ?", subscribeId, []int64{1, 0}).
			Select("subscribe_id, COUNT(id) as total").
			Group("subscribe_id").
			Scan(&result).
			Error
	})

	if err != nil {
		return nil, err
	}

	resultMap := make(map[int64]int64)
	for _, item := range result {
		resultMap[item.SubscribeId] = item.Total
	}

	return resultMap, nil
}

func (m *defaultUserModel) FindOneSubscribeByOrderId(ctx context.Context, orderId int64) (*Subscribe, error) {
	if m.subscriptionRepo.Available() {
		record, err := m.subscriptionRepo.FindByOrderID(ctx, orderId)
		if err != nil {
			return nil, err
		}
		return subscriptionRecordToLegacy(record), nil
	}
	var data Subscribe
	err := m.QueryNoCacheCtx(ctx, &data, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Subscribe{}).Where("order_id = ?", orderId).First(&data).Error
	})
	return &data, err
}

func (m *defaultUserModel) FindOneSubscribe(ctx context.Context, id int64) (*Subscribe, error) {
	if m.subscriptionRepo.Available() {
		record, err := m.subscriptionRepo.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return subscriptionRecordToLegacy(record), nil
	}
	var data Subscribe
	key := fmt.Sprintf("%s%d", cacheUserSubscribeIdPrefix, id)
	err := m.QueryCtx(ctx, &data, key, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Subscribe{}).Where("id = ?", id).First(&data).Error
	})
	return &data, err
}

func (m *defaultUserModel) FindUsersSubscribeBySubscribeId(ctx context.Context, subscribeId int64) ([]*Subscribe, error) {
	if m.subscriptionRepo.Available() {
		rows, err := m.subscriptionRepo.FindBySubscribeID(ctx, subscribeId, []int64{1, 0})
		if err != nil {
			return nil, err
		}
		result := make([]*Subscribe, 0, len(rows))
		for _, row := range rows {
			result = append(result, subscriptionRecordToLegacy(row))
		}
		return result, nil
	}
	var data []*Subscribe
	err := m.QueryNoCacheCtx(ctx, &data, func(conn *gorm.DB, v interface{}) error {
		err := conn.Model(&Subscribe{}).Where("subscribe_id = ? AND `status` IN ?", subscribeId, []int64{1, 0}).Find(v).Error

		if err != nil {
			return err
		}
		// update user subscribe status
		return conn.Model(&Subscribe{}).Where("subscribe_id = ? AND `status` = ?", subscribeId, 0).Update("status", 1).Error
	})
	return data, err
}

// QueryUserSubscribe returns a list of records that meet the conditions.
func (m *defaultUserModel) QueryUserSubscribe(ctx context.Context, userId int64, status ...int64) ([]*SubscribeDetails, error) {
	if m.subscriptionRepo.Available() {
		rows, err := m.subscriptionRepo.FindByTokenList(ctx, userId, intStatusSlice(status))
		if err != nil {
			return nil, err
		}
		list := make([]*SubscribeDetails, 0, len(rows))
		now := time.Now()
		sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
		for _, row := range rows {
			if row == nil {
				continue
			}
			if !(row.ExpireTime.After(now) || (row.FinishedAt != nil && !row.FinishedAt.Before(sevenDaysAgo)) || row.ExpireTime.Equal(time.UnixMilli(0))) {
				continue
			}
			list = append(list, m.subscriptionRecordToDetails(ctx, row))
		}
		return list, nil
	}
	var list []*SubscribeDetails
	key := fmt.Sprintf("%s%d", cacheUserSubscribeUserPrefix, userId)
	err := m.QueryCtx(ctx, &list, key, func(conn *gorm.DB, v interface{}) error {
		// 获取当前时间
		now := time.Now()
		// 获取当前时间向前推 7 天
		sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
		// 基础条件查询
		conn = conn.Model(&Subscribe{}).Where("`user_id` = ?", userId)
		if len(status) > 0 {
			conn = conn.Where("`status` IN ?", status)
		}
		// 订阅过期时间大于当前时间或者订阅结束时间大于当前时间
		return conn.Where("`expire_time` > ? OR `finished_at` >= ? OR `expire_time` = ?", now, sevenDaysAgo, time.UnixMilli(0)).
			Preload("Subscribe").
			Find(&list).Error
	})
	return list, err
}

// FindOneUserSubscribe  finds a subscribeDetails by id.
func (m *defaultUserModel) FindOneUserSubscribe(ctx context.Context, id int64) (subscribeDetails *SubscribeDetails, err error) {
	if m.subscriptionRepo.Available() {
		var record *modelsubscription.Record
		record, err = m.subscriptionRepo.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return m.subscriptionRecordToDetails(ctx, record), nil
	}
	//TODO cache
	//key := fmt.Sprintf("%s%d", cacheUserSubscribeUserPrefix, userId)
	err = m.QueryNoCacheCtx(ctx, subscribeDetails, func(conn *gorm.DB, v interface{}) error {
		if m.useIdentitySchema(conn) {
			var data SubscribeDetails
			if err := conn.Model(&Subscribe{}).Preload("Subscribe").Where("id = ?", id).First(&data).Error; err != nil {
				return err
			}
			var userRow *identity.User
			userRow, err = m.identityRepo.FindUserByID(ctx, data.UserId, conn)
			if err != nil {
				return err
			}
			data.User = m.identityUserToLegacy(userRow)
			subscribeDetails = &data
			return nil
		}
		return conn.Model(&Subscribe{}).Preload("Subscribe").Where("id = ?", id).First(&subscribeDetails).Error
	})
	return
}

func (m *defaultUserModel) ActivateAndFindUserSubscribeDetailsByIDs(ctx context.Context, ids []int64) ([]*SubscribeDetails, error) {
	if m.subscriptionRepo.Available() {
		rows, err := m.subscriptionRepo.ActivateAndFindByIDs(ctx, ids)
		if err != nil {
			return nil, err
		}
		result := make([]*SubscribeDetails, 0, len(rows))
		for _, row := range rows {
			result = append(result, m.subscriptionRecordToDetails(ctx, row))
		}
		return result, nil
	}
	ids = uniqueInt64(ids)
	if len(ids) == 0 {
		return nil, nil
	}

	var rows []*SubscribeDetails
	var pending []*Subscribe
	err := m.Transaction(ctx, func(tx *gorm.DB) error {
		query := tx.Model(&Subscribe{}).Preload("Subscribe").
			Where("id IN ? AND `status` IN ?", ids, []int64{0, 1})
		if err := query.Find(&rows).Error; err != nil {
			return err
		}
		for _, row := range rows {
			if row == nil || row.Status != 0 {
				continue
			}
			pending = append(pending, &Subscribe{
				Id:          row.Id,
				UserId:      row.UserId,
				OrderId:     row.OrderId,
				SubscribeId: row.SubscribeId,
				Token:       row.Token,
				UUID:        row.UUID,
				Status:      row.Status,
			})
		}
		if len(pending) == 0 {
			return nil
		}
		pendingIDs := make([]int64, 0, len(pending))
		for _, item := range pending {
			pendingIDs = append(pendingIDs, item.Id)
		}
		if err := tx.Model(&Subscribe{}).
			Where("id IN ? AND `status` = ?", pendingIDs, 0).
			Update("status", 1).Error; err != nil {
			return err
		}
		if err := m.clearSubscribeCacheWithTx(ctx, tx, pending...); err != nil {
			return err
		}
		for _, row := range rows {
			if row != nil && row.Status == 0 {
				row.Status = 1
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	byID := make(map[int64]*SubscribeDetails, len(rows))
	for _, row := range rows {
		if row != nil {
			byID[row.Id] = row
		}
	}
	ordered := make([]*SubscribeDetails, 0, len(rows))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			ordered = append(ordered, row)
		}
	}
	return ordered, nil
}

// FindOneSubscribeByToken  finds a record by token.
func (m *defaultUserModel) FindOneSubscribeByToken(ctx context.Context, token string) (*Subscribe, error) {
	if m.subscriptionRepo.Available() {
		record, err := m.subscriptionRepo.FindByToken(ctx, token)
		if err != nil {
			return nil, err
		}
		return subscriptionRecordToLegacy(record), nil
	}
	var data Subscribe
	key := fmt.Sprintf("%s%s", cacheUserSubscribeTokenPrefix, token)
	err := m.QueryCtx(ctx, &data, key, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Subscribe{}).Where("token = ?", token).First(&data).Error
	})
	return &data, err
}

// UpdateSubscribe updates a record.
func (m *defaultUserModel) UpdateSubscribe(ctx context.Context, data *Subscribe, tx ...*gorm.DB) error {
	old, err := m.FindOneSubscribe(ctx, data.Id)
	if err != nil {
		return err
	}

	err = m.withSubscribeWriteConn(ctx, tx, func(conn *gorm.DB) error {
		if m.subscriptionRepo.Available(conn) {
			if err := m.subscriptionRepo.Upsert(ctx, legacySubscribeToRecord(data), conn); err != nil {
				return err
			}
		} else {
			if err := conn.Model(&Subscribe{}).Where("id = ?", data.Id).Save(data).Error; err != nil {
				return err
			}
		}
		return m.syncNodeAssignments(ctx, conn, data)
	})
	if err != nil {
		return err
	}
	return m.clearSubscribeCacheWithTx(ctx, firstTx(tx), old, data)
}

// DeleteSubscribe deletes a record.
func (m *defaultUserModel) DeleteSubscribe(ctx context.Context, token string, tx ...*gorm.DB) error {
	data, err := m.FindOneSubscribeByToken(ctx, token)
	if err != nil {
		return err
	}

	err = m.withSubscribeWriteConn(ctx, tx, func(conn *gorm.DB) error {
		if m.subscriptionRepo.Available(conn) {
			if err := m.subscriptionRepo.Delete(ctx, data.Id, conn); err != nil {
				return err
			}
		} else {
			if err := conn.Where("token = ?", token).Delete(&Subscribe{}).Error; err != nil {
				return err
			}
		}
		return m.deleteNodeAssignments(ctx, conn, data.Id)
	})
	if err != nil {
		return err
	}
	return m.clearSubscribeCacheWithTx(ctx, firstTx(tx), data)
}

// InsertSubscribe insert Subscribe into the database.
func (m *defaultUserModel) InsertSubscribe(ctx context.Context, data *Subscribe, tx ...*gorm.DB) error {
	err := m.withSubscribeWriteConn(ctx, tx, func(conn *gorm.DB) error {
		if m.subscriptionRepo.Available(conn) {
			if err := m.subscriptionRepo.Upsert(ctx, legacySubscribeToRecord(data), conn); err != nil {
				return err
			}
		} else {
			if err := conn.Create(data).Error; err != nil {
				return err
			}
		}
		return m.syncNodeAssignments(ctx, conn, data)
	})
	if err != nil {
		return err
	}
	return m.clearSubscribeCacheWithTx(ctx, firstTx(tx), data)
}

func (m *defaultUserModel) DeleteSubscribeById(ctx context.Context, id int64, tx ...*gorm.DB) error {
	data, err := m.FindOneSubscribe(ctx, id)
	if err != nil {
		return err
	}

	err = m.withSubscribeWriteConn(ctx, tx, func(conn *gorm.DB) error {
		if m.subscriptionRepo.Available(conn) {
			if err := m.subscriptionRepo.Delete(ctx, id, conn); err != nil {
				return err
			}
		} else {
			if err := conn.Where("id = ?", id).Delete(&Subscribe{}).Error; err != nil {
				return err
			}
		}
		return m.deleteNodeAssignments(ctx, conn, id)
	})
	if err != nil {
		return err
	}
	return m.clearSubscribeCacheWithTx(ctx, firstTx(tx), data)
}

func (m *defaultUserModel) ClearSubscribeCache(ctx context.Context, data ...*Subscribe) error {
	return m.ClearSubscribeCacheByModels(ctx, data...)
}

func (m *defaultUserModel) syncNodeAssignments(ctx context.Context, conn *gorm.DB, data *Subscribe) error {
	if data == nil {
		return nil
	}
	if m.assignmentSyncer == nil {
		return nil
	}
	return m.assignmentSyncer.SyncUserSubscription(ctx, data.Id, data.SubscribeId, data.Status, conn)
}

func (m *defaultUserModel) deleteNodeAssignments(ctx context.Context, conn *gorm.DB, userSubscribeId int64) error {
	if m.assignmentSyncer == nil {
		return nil
	}
	return m.assignmentSyncer.DeleteUserSubscription(ctx, userSubscribeId, conn)
}

func (m *defaultUserModel) withSubscribeWriteConn(ctx context.Context, tx []*gorm.DB, fn func(conn *gorm.DB) error) error {
	if txConn := firstTx(tx); txConn != nil {
		return fn(txConn.WithContext(ctx))
	}
	return m.Transaction(ctx, fn)
}

func (m *defaultUserModel) clearSubscribeCacheWithTx(ctx context.Context, tx *gorm.DB, subscribes ...*Subscribe) error {
	if err := m.clearModelsCacheWithTx(ctx, tx, subscribesToCacheGenerators(subscribes)...); err != nil {
		logger.Errorf("failed to clear subscribe cache: %v", err)
		return err
	}
	return nil
}

func subscribesToCacheGenerators(subscribes []*Subscribe) []CacheKeyGenerator {
	models := make([]CacheKeyGenerator, 0, len(subscribes))
	for _, subscribe := range subscribes {
		if subscribe != nil {
			models = append(models, subscribe)
		}
	}
	return models
}

func firstTx(tx []*gorm.DB) *gorm.DB {
	if len(tx) == 0 {
		return nil
	}
	return tx[0]
}

func uniqueInt64(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(values))
	result := make([]int64, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func intStatusSlice(values []int64) []int64 {
	return values
}

func legacySubscribeToRecord(data *Subscribe) *modelsubscription.Record {
	if data == nil {
		return nil
	}
	return &modelsubscription.Record{
		ID:          data.Id,
		UserID:      data.UserId,
		OrderID:     data.OrderId,
		SubscribeID: data.SubscribeId,
		StartTime:   data.StartTime,
		ExpireTime:  data.ExpireTime,
		FinishedAt:  data.FinishedAt,
		Traffic:     data.Traffic,
		Download:    data.Download,
		Upload:      data.Upload,
		Token:       data.Token,
		UUID:        data.UUID,
		Status:      data.Status,
		Note:        data.Note,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func subscriptionRecordToLegacy(data *modelsubscription.Record) *Subscribe {
	if data == nil {
		return nil
	}
	return &Subscribe{
		Id:          data.ID,
		UserId:      data.UserID,
		OrderId:     data.OrderID,
		SubscribeId: data.SubscribeID,
		StartTime:   data.StartTime,
		ExpireTime:  data.ExpireTime,
		FinishedAt:  data.FinishedAt,
		Traffic:     data.Traffic,
		Download:    data.Download,
		Upload:      data.Upload,
		Token:       data.Token,
		UUID:        data.UUID,
		Status:      data.Status,
		Note:        data.Note,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

func (m *defaultUserModel) subscriptionRecordToDetails(ctx context.Context, data *modelsubscription.Record) *SubscribeDetails {
	if data == nil {
		return nil
	}
	details := &SubscribeDetails{
		Id:          data.ID,
		UserId:      data.UserID,
		OrderId:     data.OrderID,
		SubscribeId: data.SubscribeID,
		StartTime:   data.StartTime,
		ExpireTime:  data.ExpireTime,
		FinishedAt:  data.FinishedAt,
		Traffic:     data.Traffic,
		Download:    data.Download,
		Upload:      data.Upload,
		Token:       data.Token,
		UUID:        data.UUID,
		Status:      data.Status,
		Note:        data.Note,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
	if m.db != nil {
		var plan subscribe.Subscribe
		if err := m.db.WithContext(ctx).Model(&subscribe.Subscribe{}).Where("id = ?", data.SubscribeID).First(&plan).Error; err == nil {
			details.Subscribe = &plan
		}
	}
	if m.useIdentitySchema(nil) {
		if userRow, err := m.identityRepo.FindUserByID(ctx, data.UserID); err == nil {
			details.User = m.identityUserToLegacy(userRow)
		}
	} else if m.db != nil {
		var userRow User
		if err := m.db.WithContext(ctx).Model(&User{}).Where("id = ?", data.UserID).First(&userRow).Error; err == nil {
			details.User = &userRow
		}
	}
	return details
}
