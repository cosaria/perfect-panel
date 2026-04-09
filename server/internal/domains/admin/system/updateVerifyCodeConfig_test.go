package system

import (
	"context"
	"errors"
	"testing"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUpdateVerifyCodeConfigUpdatesAllFieldsAndEvictsBothCaches(t *testing.T) {
	updated := map[string]string{}
	var deletedKeys []string
	req := &types.VerifyCodeConfig{
		VerifyCodeExpireTime: 300,
		VerifyCodeLimit:      15,
		VerifyCodeInterval:   60,
	}
	deps := Deps{
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(_ context.Context, _ *gorm.DB, category, fieldName, fieldValue string) error {
			require.Equal(t, "verify_code", category)
			updated[fieldName] = fieldValue
			return nil
		},
		DeleteCacheKeys: func(_ context.Context, keys ...string) error {
			deletedKeys = append(deletedKeys, keys...)
			return nil
		},
	}
	logic := NewUpdateVerifyCodeConfigLogic(context.Background(), deps)

	err := logic.UpdateVerifyCodeConfig(req)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"VerifyCodeExpireTime": "300",
		"VerifyCodeLimit":      "15",
		"VerifyCodeInterval":   "60",
	}, updated)
	require.ElementsMatch(t, []string{config.VerifyCodeConfigKey, config.GlobalConfigKey}, deletedKeys)
}

func TestUpdateVerifyCodeConfigReturnsDatabaseUpdateErrorWhenTransactionFails(t *testing.T) {
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
	logic := NewUpdateVerifyCodeConfigLogic(context.Background(), deps)

	err := logic.UpdateVerifyCodeConfig(&types.VerifyCodeConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, deleteCalled)
}

func TestUpdateVerifyCodeConfigReturnsDatabaseUpdateErrorWhenCacheEvictionFails(t *testing.T) {
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
	logic := NewUpdateVerifyCodeConfigLogic(context.Background(), deps)

	err := logic.UpdateVerifyCodeConfig(&types.VerifyCodeConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
}
