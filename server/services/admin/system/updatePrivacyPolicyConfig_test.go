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

func TestUpdatePrivacyPolicyConfigUpdatesFieldAndEvictsTosCache(t *testing.T) {
	updated := map[string]string{}
	var deletedKeys []string
	req := &types.PrivacyPolicyConfig{
		PrivacyPolicy: "<p>Privacy first</p>",
	}
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(_ context.Context, _ *gorm.DB, category, fieldName, fieldValue string) error {
			require.Equal(t, "tos", category)
			updated[fieldName] = fieldValue
			return nil
		},
		DeleteCacheKeys: func(_ context.Context, keys ...string) error {
			deletedKeys = append(deletedKeys, keys...)
			return nil
		},
	}
	logic := NewUpdatePrivacyPolicyConfigLogic(context.Background(), deps)

	err := logic.UpdatePrivacyPolicyConfig(req)

	require.NoError(t, err)
	require.Equal(t, map[string]string{"PrivacyPolicy": req.PrivacyPolicy}, updated)
	require.ElementsMatch(t, []string{config.TosConfigKey}, deletedKeys)
}

func TestUpdatePrivacyPolicyConfigReturnsDatabaseUpdateErrorWhenTransactionFails(t *testing.T) {
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
	}
	logic := NewUpdatePrivacyPolicyConfigLogic(context.Background(), deps)

	err := logic.UpdatePrivacyPolicyConfig(&types.PrivacyPolicyConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, deleteCalled)
}

func TestUpdatePrivacyPolicyConfigReturnsDatabaseUpdateErrorWhenCacheEvictionFails(t *testing.T) {
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
	}
	logic := NewUpdatePrivacyPolicyConfigLogic(context.Background(), deps)

	err := logic.UpdatePrivacyPolicyConfig(&types.PrivacyPolicyConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
}
