package system

import (
	"context"

	"gorm.io/gorm"
)

const (
	identitySystemRevisionName = "0002_identity_system"
	schemaRegistryTable        = "schema_registry"
	revisionStateApplied       = "applied"
)

type Setting struct {
	ID        int64  `gorm:"primaryKey"`
	Category  string `gorm:"type:varchar(100);not null;index:idx_system_settings_category;comment:Category"`
	Key       string `gorm:"type:varchar(100);not null;uniqueIndex:idx_system_settings_key;comment:Key Name"`
	Value     string `gorm:"type:text;not null;comment:Key Value"`
	Type      string `gorm:"type:varchar(50);default:'';not null;comment:Type"`
	Desc      string `gorm:"type:text;not null;comment:Description"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"`
}

func (Setting) TableName() string {
	return "system_settings"
}

type Repository struct {
	db *gorm.DB
}

type AuthProviderState struct {
	Provider *AuthProvider
	Config   *AuthProviderConfig
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) HasSettingsSchema(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil {
		return false
	}
	if !r.revisionApplied(db) {
		return false
	}
	return db.Migrator().HasTable(&Setting{}) &&
		db.Migrator().HasTable(&VerificationPolicy{})
}

func (r *Repository) HasAuthProviderSchema(conn ...*gorm.DB) bool {
	db := r.conn(nil, conn...)
	if db == nil {
		return false
	}
	if !r.revisionApplied(db) {
		return false
	}
	return db.Migrator().HasTable(&AuthProvider{}) &&
		db.Migrator().HasTable(&AuthProviderConfig{})
}

func (r *Repository) ListCategoryValues(ctx context.Context, category string, tx ...*gorm.DB) ([]*System, error) {
	switch category {
	case "verify", "verify_code":
		return r.listPolicyValues(ctx, category, tx...)
	default:
		return r.listSettingValues(ctx, category, tx...)
	}
}

func (r *Repository) FindValueByKey(ctx context.Context, key string, tx ...*gorm.DB) (*System, error) {
	data, err := r.findSettingByKey(ctx, key, tx...)
	if err == nil {
		return data, nil
	}
	return r.findPolicyByKey(ctx, key, tx...)
}

func (r *Repository) FindValueByID(ctx context.Context, id int64, tx ...*gorm.DB) (*System, error) {
	var setting Setting
	if err := r.conn(ctx, tx...).Where("id = ?", id).First(&setting).Error; err == nil {
		return setting.toLegacy(), nil
	}

	var policy VerificationPolicy
	if err := r.conn(ctx, tx...).Where("id = ?", id).First(&policy).Error; err != nil {
		return nil, err
	}
	return policy.toLegacy(), nil
}

func (r *Repository) UpdateValue(ctx context.Context, category, key, value string, tx ...*gorm.DB) error {
	switch category {
	case "verify", "verify_code":
		return r.conn(ctx, tx...).Model(&VerificationPolicy{}).
			Where("category = ? AND `key` = ?", category, key).
			Update("value", value).Error
	default:
		return r.conn(ctx, tx...).Model(&Setting{}).
			Where("category = ? AND `key` = ?", category, key).
			Update("value", value).Error
	}
}

func (r *Repository) SaveValue(ctx context.Context, data *System, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	switch data.Category {
	case "verify", "verify_code":
		row := VerificationPolicy{
			ID:       data.Id,
			Category: data.Category,
			Key:      data.Key,
			Value:    data.Value,
			Type:     data.Type,
			Desc:     data.Desc,
		}
		return r.conn(ctx, tx...).Save(&row).Error
	default:
		row := Setting{
			ID:       data.Id,
			Category: data.Category,
			Key:      data.Key,
			Value:    data.Value,
			Type:     data.Type,
			Desc:     data.Desc,
		}
		return r.conn(ctx, tx...).Save(&row).Error
	}
}

func (r *Repository) DeleteValue(ctx context.Context, data *System, tx ...*gorm.DB) error {
	if data == nil {
		return nil
	}
	switch data.Category {
	case "verify", "verify_code":
		return r.conn(ctx, tx...).Delete(&VerificationPolicy{}, data.Id).Error
	default:
		return r.conn(ctx, tx...).Delete(&Setting{}, data.Id).Error
	}
}

func (r *Repository) FindAuthProviderByMethod(ctx context.Context, method string, tx ...*gorm.DB) (*AuthProviderState, error) {
	var provider AuthProvider
	if err := r.conn(ctx, tx...).Where("method = ?", method).First(&provider).Error; err != nil {
		return nil, err
	}
	var config AuthProviderConfig
	if err := r.conn(ctx, tx...).Where("provider_id = ?", provider.ID).First(&config).Error; err != nil {
		return nil, err
	}
	return &AuthProviderState{Provider: &provider, Config: &config}, nil
}

func (r *Repository) FindAuthProviderByID(ctx context.Context, id int64, tx ...*gorm.DB) (*AuthProviderState, error) {
	var provider AuthProvider
	if err := r.conn(ctx, tx...).Where("id = ?", id).First(&provider).Error; err != nil {
		return nil, err
	}
	var config AuthProviderConfig
	if err := r.conn(ctx, tx...).Where("provider_id = ?", provider.ID).First(&config).Error; err != nil {
		return nil, err
	}
	return &AuthProviderState{Provider: &provider, Config: &config}, nil
}

func (r *Repository) ListAuthProviders(ctx context.Context, tx ...*gorm.DB) ([]*AuthProviderState, error) {
	var providers []*AuthProvider
	if err := r.conn(ctx, tx...).Find(&providers).Error; err != nil {
		return nil, err
	}

	result := make([]*AuthProviderState, 0, len(providers))
	for _, provider := range providers {
		var config AuthProviderConfig
		if err := r.conn(ctx, tx...).Where("provider_id = ?", provider.ID).First(&config).Error; err != nil {
			return nil, err
		}
		result = append(result, &AuthProviderState{
			Provider: provider,
			Config:   &config,
		})
	}
	return result, nil
}

func (r *Repository) UpsertAuthProvider(ctx context.Context, method, config string, enabled *bool, tx ...*gorm.DB) (*AuthProviderState, error) {
	conn := r.conn(ctx, tx...)
	var result *AuthProviderState
	err := conn.Transaction(func(tx *gorm.DB) error {
		var provider AuthProvider
		err := tx.Where("method = ?", method).First(&provider).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			provider = AuthProvider{Method: method, Enabled: enabled}
			if err := tx.Create(&provider).Error; err != nil {
				return err
			}
		} else {
			provider.Enabled = enabled
			if err := tx.Save(&provider).Error; err != nil {
				return err
			}
		}

		var providerConfig AuthProviderConfig
		err = tx.Where("provider_id = ?", provider.ID).First(&providerConfig).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			providerConfig = AuthProviderConfig{
				ProviderID: provider.ID,
				Config:     config,
			}
			if err := tx.Create(&providerConfig).Error; err != nil {
				return err
			}
		} else {
			providerConfig.Config = config
			if err := tx.Save(&providerConfig).Error; err != nil {
				return err
			}
		}

		result = &AuthProviderState{
			Provider: &provider,
			Config:   &providerConfig,
		}
		return nil
	})
	return result, err
}

func (r *Repository) DeleteAuthProvider(ctx context.Context, method string, tx ...*gorm.DB) error {
	conn := r.conn(ctx, tx...)
	return conn.Transaction(func(tx *gorm.DB) error {
		var provider AuthProvider
		if err := tx.Where("method = ?", method).First(&provider).Error; err != nil {
			return err
		}
		if err := tx.Where("provider_id = ?", provider.ID).Delete(&AuthProviderConfig{}).Error; err != nil {
			return err
		}
		return tx.Delete(&provider).Error
	})
}

func (r *Repository) listSettingValues(ctx context.Context, category string, tx ...*gorm.DB) ([]*System, error) {
	var rows []*Setting
	if err := r.conn(ctx, tx...).Where("category = ?", category).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*System, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.toLegacy())
	}
	return result, nil
}

func (r *Repository) findSettingByKey(ctx context.Context, key string, tx ...*gorm.DB) (*System, error) {
	var row Setting
	if err := r.conn(ctx, tx...).Where("`key` = ?", key).First(&row).Error; err != nil {
		return nil, err
	}
	return row.toLegacy(), nil
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
	if db == nil || !db.Migrator().HasTable(schemaRegistryTable) {
		return false
	}

	var count int64
	if err := db.Table(schemaRegistryTable).
		Where("id = ? AND state = ?", identitySystemRevisionName, revisionStateApplied).
		Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (s *Setting) toLegacy() *System {
	if s == nil {
		return nil
	}
	return &System{
		Id:       s.ID,
		Category: s.Category,
		Key:      s.Key,
		Value:    s.Value,
		Type:     s.Type,
		Desc:     s.Desc,
	}
}
