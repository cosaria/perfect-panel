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

func TestUpdateInviteConfigUpdatesAllFieldsEvictsCacheAndReloads(t *testing.T) {
	updated := map[string]string{}
	var deletedKeys []string
	reloadCount := 0
	req := &types.InviteConfig{
		ForcedInvite:       true,
		ReferralPercentage: 35,
		OnlyFirstPurchase:  true,
	}
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(_ context.Context, _ *gorm.DB, category, fieldName, fieldValue string) error {
			require.Equal(t, "invite", category)
			updated[fieldName] = fieldValue
			return nil
		},
		DeleteCacheKeys: func(_ context.Context, keys ...string) error {
			deletedKeys = append(deletedKeys, keys...)
			return nil
		},
		RunReloadInvite: func() error {
			reloadCount++
			return nil
		},
	}
	logic := NewUpdateInviteConfigLogic(context.Background(), deps)

	err := logic.UpdateInviteConfig(req)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"ForcedInvite":       "true",
		"ReferralPercentage": "35",
		"OnlyFirstPurchase":  "true",
	}, updated)
	require.ElementsMatch(t, []string{config.InviteConfigKey, config.GlobalConfigKey}, deletedKeys)
	require.Equal(t, 1, reloadCount)
}

func TestUpdateInviteConfigReturnsDatabaseUpdateErrorWhenTransactionFails(t *testing.T) {
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
		RunReloadInvite: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateInviteConfigLogic(context.Background(), deps)

	err := logic.UpdateInviteConfig(&types.InviteConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, deleteCalled)
	require.False(t, reloadCalled)
}

func TestUpdateInviteConfigReturnsDatabaseUpdateErrorWhenCacheEvictionFails(t *testing.T) {
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
		RunReloadInvite: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateInviteConfigLogic(context.Background(), deps)

	err := logic.UpdateInviteConfig(&types.InviteConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, reloadCalled)
}

func TestUpdateInviteConfigAllowsMissingReloadHook(t *testing.T) {
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
	logic := NewUpdateInviteConfigLogic(context.Background(), deps)

	err := logic.UpdateInviteConfig(&types.InviteConfig{})

	require.NoError(t, err)
}

func TestUpdateInviteConfigReturnsErrorWhenReloadHookFails(t *testing.T) {
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
		RunReloadInvite: func() error {
			return errors.New("reload failed")
		},
	}
	logic := NewUpdateInviteConfigLogic(context.Background(), deps)

	err := logic.UpdateInviteConfig(&types.InviteConfig{})

	requireSystemCodeError(t, err, xerr.ERROR)
}
