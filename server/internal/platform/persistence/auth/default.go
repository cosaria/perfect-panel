package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/perfect-panel/server/internal/platform/cache"
	modelsystem "github.com/perfect-panel/server/internal/platform/persistence/system"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ Model = (*customAuthModel)(nil)
var (
	cacheAuthIdPrefix     = "cache:auth:id:"
	cacheAuthMethodPrefix = "cache:auth:method:"
)

type (
	Model interface {
		authModel
		customAuthLogicModel
	}
	authModel interface {
		Insert(ctx context.Context, data *Auth) error
		FindOne(ctx context.Context, id int64) (*Auth, error)
		Update(ctx context.Context, data *Auth) error
		Delete(ctx context.Context, id int64) error
		Transaction(ctx context.Context, fn func(db *gorm.DB) error) error
	}

	customAuthModel struct {
		*defaultAuthModel
	}
	defaultAuthModel struct {
		cache.CachedConn
		db         *gorm.DB
		table      string
		systemRepo *modelsystem.Repository
	}
)

func newAuthModel(db *gorm.DB, c *redis.Client) *defaultAuthModel {
	return &defaultAuthModel{
		CachedConn: cache.NewConn(db, c),
		db:         db,
		table:      "`auth_config`",
		systemRepo: modelsystem.NewRepository(db),
	}
}

//nolint:unused
func (m *defaultAuthModel) batchGetCacheKeys(Auths ...*Auth) []string {
	var keys []string
	for _, auth := range Auths {
		keys = append(keys, m.getCacheKeys(auth)...)
	}
	return keys

}
func (m *defaultAuthModel) getCacheKeys(data *Auth) []string {
	if data == nil {
		return []string{}
	}
	authIdKey := fmt.Sprintf("%s%v", cacheAuthIdPrefix, data.Id)
	platformKey := fmt.Sprintf("%s%s", cacheAuthMethodPrefix, data.Method)
	cacheKeys := []string{
		authIdKey,
		platformKey,
	}
	return cacheKeys
}

func (m *defaultAuthModel) Insert(ctx context.Context, data *Auth) error {
	err := m.ExecCtx(ctx, func(conn *gorm.DB) error {
		if m.useAuthProviderSchema(conn) {
			state, err := m.systemRepo.UpsertAuthProvider(ctx, data.Method, data.Config, data.Enabled, conn)
			if err != nil {
				return err
			}
			data.Id = state.Provider.ID
			return nil
		}
		return conn.Create(&data).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *defaultAuthModel) FindOne(ctx context.Context, id int64) (*Auth, error) {
	if m.useAuthProviderSchema(nil) {
		state, err := m.systemRepo.FindAuthProviderByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return m.authProviderStateToLegacy(state), nil
	}
	AuthIdKey := fmt.Sprintf("%s%v", cacheAuthIdPrefix, id)
	var resp Auth
	err := m.QueryCtx(ctx, &resp, AuthIdKey, func(conn *gorm.DB, v interface{}) error {
		return conn.Model(&Auth{}).Where("`id` = ?", id).First(&resp).Error
	})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *defaultAuthModel) Update(ctx context.Context, data *Auth) error {
	old, err := m.FindOne(ctx, data.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		db := conn
		if m.useAuthProviderSchema(db) {
			_, err := m.systemRepo.UpsertAuthProvider(ctx, data.Method, data.Config, data.Enabled, db)
			return err
		}
		return db.Save(data).Error
	}, m.getCacheKeys(old)...)
	return err
}

func (m *defaultAuthModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		db := conn
		if m.useAuthProviderSchema(db) {
			return m.systemRepo.DeleteAuthProvider(ctx, data.Method, db)
		}
		return db.Delete(&Auth{}, id).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *defaultAuthModel) Transaction(ctx context.Context, fn func(db *gorm.DB) error) error {
	return m.TransactCtx(ctx, fn)
}

func (m *defaultAuthModel) useAuthProviderSchema(conn *gorm.DB) bool {
	if m.systemRepo == nil {
		return false
	}
	return m.systemRepo.HasAuthProviderSchema(conn)
}

func (m *defaultAuthModel) authProviderStateToLegacy(state *modelsystem.AuthProviderState) *Auth {
	if state == nil || state.Provider == nil {
		return nil
	}
	config := ""
	if state.Config != nil {
		config = state.Config.Config
	}
	return &Auth{
		Id:      state.Provider.ID,
		Method:  state.Provider.Method,
		Config:  config,
		Enabled: state.Provider.Enabled,
	}
}
