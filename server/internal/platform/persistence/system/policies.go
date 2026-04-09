package system

import (
	"context"

	"gorm.io/gorm"
)

type VerificationPolicy struct {
	ID        int64  `gorm:"primaryKey"`
	Category  string `gorm:"type:varchar(100);not null;index:idx_verification_policies_category;comment:Category"`
	Key       string `gorm:"column:key;type:varchar(100);not null;uniqueIndex:idx_verification_policies_key;comment:Key Name"`
	Value     string `gorm:"type:text;not null;comment:Key Value"`
	Type      string `gorm:"type:varchar(50);default:'';not null;comment:Type"`
	Desc      string `gorm:"type:text;not null;comment:Description"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"`
}

func (VerificationPolicy) TableName() string {
	return "verification_policies"
}

type AuthProvider struct {
	ID        int64  `gorm:"primaryKey"`
	Method    string `gorm:"type:varchar(100);not null;uniqueIndex:idx_auth_providers_method;comment:Auth Method"`
	Enabled   *bool `gorm:"type:tinyint(1);not null;default:false;comment:Is Enabled"`
	CreatedAt int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func (AuthProvider) TableName() string {
	return "auth_providers"
}

type AuthProviderConfig struct {
	ID         int64 `gorm:"primaryKey"`
	ProviderID int64 `gorm:"not null;uniqueIndex:idx_auth_provider_configs_provider_id;comment:Auth Provider ID"`
	Config     string
	CreatedAt  int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt  int64 `gorm:"autoUpdateTime:milli"`
}

func (AuthProviderConfig) TableName() string {
	return "auth_provider_configs"
}

func (r *Repository) listPolicyValues(ctx context.Context, category string, tx ...*gorm.DB) ([]*System, error) {
	var rows []*VerificationPolicy
	if err := r.conn(ctx, tx...).Where("category = ?", category).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*System, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.toLegacy())
	}
	return result, nil
}

func (r *Repository) findPolicyByKey(ctx context.Context, key string, tx ...*gorm.DB) (*System, error) {
	var row VerificationPolicy
	if err := r.conn(ctx, tx...).Where("`key` = ?", key).First(&row).Error; err != nil {
		return nil, err
	}
	return row.toLegacy(), nil
}

func (p *VerificationPolicy) toLegacy() *System {
	if p == nil {
		return nil
	}
	return &System{
		Id:       p.ID,
		Category: p.Category,
		Key:      p.Key,
		Value:    p.Value,
		Type:     p.Type,
		Desc:     p.Desc,
	}
}
