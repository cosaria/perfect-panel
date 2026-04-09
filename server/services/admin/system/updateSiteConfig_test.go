package system

import (
	"context"
	"errors"
	"testing"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUpdateSiteConfigUpdatesAllFieldsEvictsCacheAndReloads(t *testing.T) {
	updated := map[string]string{}
	var deletedKeys []string
	reloadCount := 0
	req := &types.SiteConfig{
		Host:       "https://example.com",
		SiteName:   "Perfect Panel",
		SiteDesc:   "panel",
		SiteLogo:   "/logo.png",
		Keywords:   "proxy,panel",
		CustomHTML: "<script>console.log('ok')</script>",
		CustomData: "{\"theme\":\"light\"}",
	}
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(ctx context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateSiteField: func(_ context.Context, _ *gorm.DB, fieldName, fieldValue string) error {
			updated[fieldName] = fieldValue
			return nil
		},
		DeleteCacheKeys: func(_ context.Context, keys ...string) error {
			deletedKeys = append(deletedKeys, keys...)
			return nil
		},
		RunReloadSite: func() error {
			reloadCount++
			return nil
		},
	}
	logic := NewUpdateSiteConfigLogic(context.Background(), deps)

	err := logic.UpdateSiteConfig(req)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"Host":       req.Host,
		"SiteName":   req.SiteName,
		"SiteDesc":   req.SiteDesc,
		"SiteLogo":   req.SiteLogo,
		"Keywords":   req.Keywords,
		"CustomHTML": req.CustomHTML,
		"CustomData": req.CustomData,
	}, updated)
	require.ElementsMatch(t, []string{config.SiteConfigKey, config.GlobalConfigKey}, deletedKeys)
	require.Equal(t, 1, reloadCount)
}

func TestUpdateSiteConfigReturnsDatabaseUpdateErrorWhenTransactionFails(t *testing.T) {
	reloadCalled := false
	deleteCalled := false
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(context.Context, func(*gorm.DB) error) error {
				return errors.New("tx failed")
			},
		},
		DeleteCacheKeys: func(context.Context, ...string) error {
			deleteCalled = true
			return nil
		},
		RunReloadSite: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateSiteConfigLogic(context.Background(), deps)

	err := logic.UpdateSiteConfig(&types.SiteConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, deleteCalled)
	require.False(t, reloadCalled)
}

func TestUpdateSiteConfigReturnsDatabaseUpdateErrorWhenCacheEvictionFails(t *testing.T) {
	reloadCalled := false
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(ctx context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateSiteField: func(context.Context, *gorm.DB, string, string) error {
			return nil
		},
		DeleteCacheKeys: func(context.Context, ...string) error {
			return errors.New("redis delete failed")
		},
		RunReloadSite: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateSiteConfigLogic(context.Background(), deps)

	err := logic.UpdateSiteConfig(&types.SiteConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, reloadCalled)
}

func TestUpdateSiteConfigAllowsMissingReloadHook(t *testing.T) {
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(ctx context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateSiteField: func(context.Context, *gorm.DB, string, string) error {
			return nil
		},
		DeleteCacheKeys: func(context.Context, ...string) error {
			return nil
		},
	}
	logic := NewUpdateSiteConfigLogic(context.Background(), deps)

	err := logic.UpdateSiteConfig(&types.SiteConfig{})

	require.NoError(t, err)
}

func TestUpdateSiteConfigReturnsErrorWhenReloadHookFails(t *testing.T) {
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(ctx context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateSiteField: func(context.Context, *gorm.DB, string, string) error {
			return nil
		},
		DeleteCacheKeys: func(context.Context, ...string) error {
			return nil
		},
		RunReloadSite: func() error {
			return errors.New("reload failed")
		},
	}
	logic := NewUpdateSiteConfigLogic(context.Background(), deps)

	err := logic.UpdateSiteConfig(&types.SiteConfig{})

	requireSystemCodeError(t, err, xerr.ERROR)
}
