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

func TestUpdateRegisterConfigUpdatesAllFieldsEvictsCacheAndReloads(t *testing.T) {
	updated := map[string]string{}
	var deletedKeys []string
	reloadCount := 0
	req := &types.RegisterConfig{
		StopRegister:            true,
		EnableTrial:             true,
		TrialSubscribe:          9,
		TrialTime:               30,
		TrialTimeUnit:           "day",
		EnableIpRegisterLimit:   true,
		IpRegisterLimit:         2,
		IpRegisterLimitDuration: 86400,
	}
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(_ context.Context, _ *gorm.DB, category, fieldName, fieldValue string) error {
			require.Equal(t, "register", category)
			updated[fieldName] = fieldValue
			return nil
		},
		DeleteCacheKeys: func(_ context.Context, keys ...string) error {
			deletedKeys = append(deletedKeys, keys...)
			return nil
		},
		RunReloadRegister: func() error {
			reloadCount++
			return nil
		},
	}
	logic := NewUpdateRegisterConfigLogic(context.Background(), deps)

	err := logic.UpdateRegisterConfig(req)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"StopRegister":            "true",
		"EnableTrial":             "true",
		"TrialSubscribe":          "9",
		"TrialTime":               "30",
		"TrialTimeUnit":           "day",
		"EnableIpRegisterLimit":   "true",
		"IpRegisterLimit":         "2",
		"IpRegisterLimitDuration": "86400",
	}, updated)
	require.ElementsMatch(t, []string{config.RegisterConfigKey, config.GlobalConfigKey}, deletedKeys)
	require.Equal(t, 1, reloadCount)
}

func TestUpdateRegisterConfigReturnsDatabaseUpdateErrorWhenTransactionFails(t *testing.T) {
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
		RunReloadRegister: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateRegisterConfigLogic(context.Background(), deps)

	err := logic.UpdateRegisterConfig(&types.RegisterConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, deleteCalled)
	require.False(t, reloadCalled)
}

func TestUpdateRegisterConfigReturnsDatabaseUpdateErrorWhenCacheEvictionFails(t *testing.T) {
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
		RunReloadRegister: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateRegisterConfigLogic(context.Background(), deps)

	err := logic.UpdateRegisterConfig(&types.RegisterConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, reloadCalled)
}

func TestUpdateRegisterConfigAllowsMissingReloadHook(t *testing.T) {
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
	logic := NewUpdateRegisterConfigLogic(context.Background(), deps)

	err := logic.UpdateRegisterConfig(&types.RegisterConfig{})

	require.NoError(t, err)
}

func TestUpdateRegisterConfigReturnsErrorWhenReloadHookFails(t *testing.T) {
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
		RunReloadRegister: func() error {
			return errors.New("reload failed")
		},
	}
	logic := NewUpdateRegisterConfigLogic(context.Background(), deps)

	err := logic.UpdateRegisterConfig(&types.RegisterConfig{})

	requireSystemCodeError(t, err, xerr.ERROR)
}
