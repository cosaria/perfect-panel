package system

import (
	"context"
	"errors"
	"testing"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUpdateCurrencyConfigUpdatesAllFieldsEvictsCacheAndReloads(t *testing.T) {
	updated := map[string]string{}
	var deletedKeys []string
	reloadCount := 0
	req := &types.CurrencyConfig{
		AccessKey:      "currency-key",
		CurrencyUnit:   "USD",
		CurrencySymbol: "$",
	}
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(_ context.Context, _ *gorm.DB, category, fieldName, fieldValue string) error {
			require.Equal(t, "currency", category)
			updated[fieldName] = fieldValue
			return nil
		},
		DeleteCacheKeys: func(_ context.Context, keys ...string) error {
			deletedKeys = append(deletedKeys, keys...)
			return nil
		},
		RunReloadCurrency: func() error {
			reloadCount++
			return nil
		},
	}
	logic := NewUpdateCurrencyConfigLogic(context.Background(), deps)

	err := logic.UpdateCurrencyConfig(req)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"AccessKey":      req.AccessKey,
		"CurrencyUnit":   req.CurrencyUnit,
		"CurrencySymbol": req.CurrencySymbol,
	}, updated)
	require.ElementsMatch(t, []string{config.CurrencyConfigKey, config.GlobalConfigKey}, deletedKeys)
	require.Equal(t, 1, reloadCount)
}

func TestUpdateCurrencyConfigReturnsDatabaseUpdateErrorWhenTransactionFails(t *testing.T) {
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
		RunReloadCurrency: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateCurrencyConfigLogic(context.Background(), deps)

	err := logic.UpdateCurrencyConfig(&types.CurrencyConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, deleteCalled)
	require.False(t, reloadCalled)
}

func TestUpdateCurrencyConfigReturnsDatabaseUpdateErrorWhenCacheEvictionFails(t *testing.T) {
	reloadCalled := false
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(context.Context, *gorm.DB, string, string, string) error {
			return nil
		},
		DeleteCacheKeys: func(context.Context, ...string) error {
			return errors.New("redis delete failed")
		},
		RunReloadCurrency: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateCurrencyConfigLogic(context.Background(), deps)

	err := logic.UpdateCurrencyConfig(&types.CurrencyConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, reloadCalled)
}

func TestUpdateCurrencyConfigAllowsMissingReloadHook(t *testing.T) {
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(context.Context, *gorm.DB, string, string, string) error {
			return nil
		},
		DeleteCacheKeys: func(context.Context, ...string) error {
			return nil
		},
	}
	logic := NewUpdateCurrencyConfigLogic(context.Background(), deps)

	err := logic.UpdateCurrencyConfig(&types.CurrencyConfig{})

	require.NoError(t, err)
}

func TestUpdateCurrencyConfigReturnsErrorWhenReloadHookFails(t *testing.T) {
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(context.Context, *gorm.DB, string, string, string) error {
			return nil
		},
		DeleteCacheKeys: func(context.Context, ...string) error {
			return nil
		},
		RunReloadCurrency: func() error {
			return errors.New("reload failed")
		},
	}
	logic := NewUpdateCurrencyConfigLogic(context.Background(), deps)

	err := logic.UpdateCurrencyConfig(&types.CurrencyConfig{})

	requireSystemCodeError(t, err, xerr.ERROR)
}
