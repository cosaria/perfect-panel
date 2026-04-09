package subscription

import (
	"context"
	"sort"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	subscriptionRevisionName  = "0004_billing_subscription"
	subscriptionRegistryTable = "schema_registry"
	subscriptionRevisionState = "applied"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
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
	return db.Migrator().HasTable(&Subscription{}) &&
		db.Migrator().HasTable(&SubscriptionToken{}) &&
		db.Migrator().HasTable(&SubscriptionPeriod{}) &&
		db.Migrator().HasTable(&SubscriptionUsageSnapshot{})
}

func (r *Repository) FindByID(ctx context.Context, id int64, tx ...*gorm.DB) (*Record, error) {
	return r.findOne(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("subscriptions.id = ?", id)
	}, tx...)
}

func (r *Repository) FindByToken(ctx context.Context, token string, tx ...*gorm.DB) (*Record, error) {
	return r.findOne(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("subscription_tokens.token = ?", token)
	}, tx...)
}

func (r *Repository) FindByOrderID(ctx context.Context, orderID int64, tx ...*gorm.DB) (*Record, error) {
	return r.findOne(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("subscriptions.order_id = ?", orderID)
	}, tx...)
}

func (r *Repository) FindByTokenList(ctx context.Context, userID int64, statuses []int64, tx ...*gorm.DB) ([]*Record, error) {
	query := r.baseQuery(ctx, tx...).Where("subscriptions.user_id = ?", userID)
	if len(statuses) > 0 {
		query = query.Where("subscriptions.status IN ?", statuses)
	}
	var rows []joinedRow
	if err := query.Order("subscriptions.id desc").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rowsToRecords(rows), nil
}

func (r *Repository) FindBySubscribeID(ctx context.Context, subscribeID int64, statuses []int64, tx ...*gorm.DB) ([]*Record, error) {
	query := r.baseQuery(ctx, tx...).Where("subscriptions.subscribe_id = ?", subscribeID)
	if len(statuses) > 0 {
		query = query.Where("subscriptions.status IN ?", statuses)
	}
	var rows []joinedRow
	if err := query.Order("subscriptions.id asc").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rowsToRecords(rows), nil
}

func (r *Repository) ActivateAndFindByIDs(ctx context.Context, ids []int64, tx ...*gorm.DB) ([]*Record, error) {
	ids = uniqueInt64(ids)
	if len(ids) == 0 {
		return nil, nil
	}
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil, nil
	}
	var records []*Record
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Subscription{}).
			Where("id IN ? AND status = ?", ids, 0).
			Update("status", 1).Error; err != nil {
			return err
		}
		if tx.Migrator().HasTable("user_subscribe") {
			if err := tx.Table("user_subscribe").
				Where("id IN ? AND status = ?", ids, 0).
				Update("status", 1).Error; err != nil {
				return err
			}
		}
		rows, err := r.findMany(ctx, ids, tx)
		if err != nil {
			return err
		}
		records = rows
		return nil
	})
	return records, err
}

func (r *Repository) Upsert(ctx context.Context, data *Record, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil
	}
	if data.ExpireTime.IsZero() {
		data.ExpireTime = time.UnixMilli(0)
	}
	subscription := Subscription{
		ID:          data.ID,
		UserID:      data.UserID,
		OrderID:     data.OrderID,
		SubscribeID: data.SubscribeID,
		StartTime:   data.StartTime,
		ExpireTime:  data.ExpireTime,
		FinishedAt:  data.FinishedAt,
		Traffic:     data.Traffic,
		Download:    data.Download,
		Upload:      data.Upload,
		Status:      data.Status,
		Note:        data.Note,
	}
	period := SubscriptionPeriod{
		ID:             data.ID,
		SubscriptionID: data.ID,
		StartTime:      data.StartTime,
		ExpireTime:     data.ExpireTime,
		FinishedAt:     data.FinishedAt,
		Status:         data.Status,
	}
	token := SubscriptionToken{
		ID:             data.ID,
		SubscriptionID: data.ID,
		Token:          data.Token,
		UUID:           data.UUID,
		IsPrimary:      true,
	}
	usage := SubscriptionUsageSnapshot{
		ID:             data.ID,
		SubscriptionID: data.ID,
		Traffic:        data.Traffic,
		Download:       data.Download,
		Upload:         data.Upload,
		CapturedAt:     time.Now().UTC(),
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&subscription).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&period).Error; err != nil {
			return err
		}
		if err := tx.Where("subscription_id = ?", data.ID).Delete(&SubscriptionToken{}).Error; err != nil {
			return err
		}
		if err := tx.Create(&token).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&usage).Error; err != nil {
			return err
		}
		if tx.Migrator().HasTable("user_subscribe") {
			if err := tx.Table("user_subscribe").Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(subscriptionLegacyMap(data)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) Delete(ctx context.Context, id int64, tx ...*gorm.DB) error {
	db := r.conn(ctx, tx...)
	if db == nil {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("subscription_id = ?", id).Delete(&SubscriptionToken{}).Error; err != nil {
			return err
		}
		if err := tx.Where("subscription_id = ?", id).Delete(&SubscriptionPeriod{}).Error; err != nil {
			return err
		}
		if err := tx.Where("subscription_id = ?", id).Delete(&SubscriptionUsageSnapshot{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&Subscription{}, id).Error; err != nil {
			return err
		}
		if tx.Migrator().HasTable("user_subscribe") {
			if err := tx.Exec("DELETE FROM user_subscribe WHERE id = ?", id).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) CountActiveBySubscribeID(ctx context.Context, subscribeIDs []int64, tx ...*gorm.DB) (map[int64]int64, error) {
	result := make(map[int64]int64)
	if len(subscribeIDs) == 0 {
		return result, nil
	}
	type row struct {
		SubscribeID int64
		Total       int64
	}
	var rows []row
	err := r.conn(ctx, tx...).Model(&Subscription{}).
		Where("subscribe_id IN ? AND status IN ?", subscribeIDs, []int64{0, 1}).
		Select("subscribe_id, COUNT(id) AS total").
		Group("subscribe_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, item := range rows {
		result[item.SubscribeID] = item.Total
	}
	return result, nil
}

type joinedRow struct {
	ID          int64
	UserID      int64
	OrderID     int64
	SubscribeID int64
	StartTime   time.Time
	ExpireTime  time.Time
	FinishedAt  *time.Time
	Traffic     int64
	Download    int64
	Upload      int64
	Status      uint8
	Note        string
	Token       string
	UUID        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *Repository) baseQuery(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	return r.conn(ctx, tx...).Model(&Subscription{}).
		Joins("LEFT JOIN subscription_tokens ON subscription_tokens.subscription_id = subscriptions.id AND subscription_tokens.is_primary = ?", true).
		Select("subscriptions.*, subscription_tokens.token, subscription_tokens.uuid")
}

func (r *Repository) findOne(ctx context.Context, scope func(*gorm.DB) *gorm.DB, tx ...*gorm.DB) (*Record, error) {
	var row joinedRow
	query := r.baseQuery(ctx, tx...)
	if scope != nil {
		query = scope(query)
	}
	err := query.Take(&row).Error
	if err != nil {
		return nil, err
	}
	return rowToRecord(row), nil
}

func (r *Repository) findMany(ctx context.Context, ids []int64, tx ...*gorm.DB) ([]*Record, error) {
	var rows []joinedRow
	err := r.baseQuery(ctx, tx...).
		Where("subscriptions.id IN ? AND subscriptions.status IN ?", ids, []int64{0, 1}).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	byID := make(map[int64]*Record, len(rows))
	for _, row := range rows {
		byID[row.ID] = rowToRecord(row)
	}
	result := make([]*Record, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			result = append(result, row)
		}
	}
	return result, nil
}

func rowToRecord(row joinedRow) *Record {
	return &Record{
		ID:          row.ID,
		UserID:      row.UserID,
		OrderID:     row.OrderID,
		SubscribeID: row.SubscribeID,
		StartTime:   row.StartTime,
		ExpireTime:  row.ExpireTime,
		FinishedAt:  row.FinishedAt,
		Traffic:     row.Traffic,
		Download:    row.Download,
		Upload:      row.Upload,
		Token:       row.Token,
		UUID:        row.UUID,
		Status:      row.Status,
		Note:        row.Note,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func rowsToRecords(rows []joinedRow) []*Record {
	result := make([]*Record, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToRecord(row))
	}
	return result
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

func subscriptionLegacyMap(data *Record) map[string]any {
	return map[string]any{
		"id":           data.ID,
		"user_id":      data.UserID,
		"order_id":     data.OrderID,
		"subscribe_id": data.SubscribeID,
		"start_time":   data.StartTime,
		"expire_time":  data.ExpireTime,
		"finished_at":  data.FinishedAt,
		"traffic":      data.Traffic,
		"download":     data.Download,
		"upload":       data.Upload,
		"token":        data.Token,
		"uuid":         data.UUID,
		"status":       data.Status,
		"note":         data.Note,
		"created_at":   data.CreatedAt,
		"updated_at":   data.UpdatedAt,
	}
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
	if db == nil || !db.Migrator().HasTable(subscriptionRegistryTable) {
		return false
	}
	var count int64
	if err := db.Table(subscriptionRegistryTable).
		Where("id = ? AND state = ?", subscriptionRevisionName, subscriptionRevisionState).
		Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
