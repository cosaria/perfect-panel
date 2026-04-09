package auth

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type customAuthLogicModel interface {
	GetAuthListByPage(ctx context.Context) ([]*Auth, error)
	FindOneByMethod(ctx context.Context, platform string) (*Auth, error)
	FindAll(ctx context.Context) ([]*Auth, error)
}

// NewModel returns a model for the database table.
func NewModel(conn *gorm.DB, c *redis.Client) Model {
	return &customAuthModel{
		defaultAuthModel: newAuthModel(conn, c),
	}
}

type Filter struct {
	Show   *bool
	Pinned *bool
	Popup  *bool
	Search string
}

// GetAuthListByPage  get auth list by page
func (m *customAuthModel) GetAuthListByPage(ctx context.Context) ([]*Auth, error) {
	if m.useAuthProviderSchema(nil) {
		return m.FindAll(ctx)
	}
	var list []*Auth
	err := m.QueryNoCacheCtx(ctx, &list, func(conn *gorm.DB, v interface{}) error {
		conn = conn.Model(&Auth{})
		return conn.Find(v).Error
	})
	return list, err
}

// FindOneByMethod  find one by method
func (m *customAuthModel) FindOneByMethod(ctx context.Context, method string) (*Auth, error) {
	if m.useAuthProviderSchema(nil) {
		state, err := m.systemRepo.FindAuthProviderByMethod(ctx, method)
		if err != nil {
			return nil, err
		}
		return m.authProviderStateToLegacy(state), nil
	}
	key := fmt.Sprintf("%s%s", cacheAuthMethodPrefix, method)
	var data Auth
	err := m.QueryCtx(ctx, &data, key, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Auth{}).Where("method = ?", method).First(v).Error
	})

	return &data, err
}

// FindAll find all
func (m *customAuthModel) FindAll(ctx context.Context) ([]*Auth, error) {
	if m.useAuthProviderSchema(nil) {
		states, err := m.systemRepo.ListAuthProviders(ctx)
		if err != nil {
			return nil, err
		}
		result := make([]*Auth, 0, len(states))
		for _, state := range states {
			result = append(result, m.authProviderStateToLegacy(state))
		}
		return result, nil
	}
	var list []*Auth
	err := m.QueryNoCacheCtx(ctx, &list, func(conn *gorm.DB, v interface{}) error {
		conn = conn.Model(&Auth{})
		return conn.Find(v).Error
	})
	return list, err
}
