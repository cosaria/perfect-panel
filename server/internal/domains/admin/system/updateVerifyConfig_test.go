package system

import (
	"context"
	"errors"
	"testing"

	serverconfig "github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUpdateVerifyConfigUpdatesAllFieldsEvictsCacheReloadsAndSyncsConfig(t *testing.T) {
	updated := map[string]string{}
	var deletedKeys []string
	reloadCount := 0
	cfg := &serverconfig.Config{}
	req := &types.VerifyConfig{
		TurnstileSiteKey:          "site-key",
		TurnstileSecret:           "secret-key",
		EnableLoginVerify:         true,
		EnableRegisterVerify:      true,
		EnableResetPasswordVerify: true,
	}
	deps := Deps{
		Config: cfg,
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(_ context.Context, _ *gorm.DB, category, fieldName, fieldValue string) error {
			require.Equal(t, "verify", category)
			updated[fieldName] = fieldValue
			return nil
		},
		DeleteCacheKeys: func(_ context.Context, keys ...string) error {
			deletedKeys = append(deletedKeys, keys...)
			return nil
		},
		RunReloadVerify: func() error {
			reloadCount++
			return nil
		},
	}
	logic := NewUpdateVerifyConfigLogic(context.Background(), deps)

	err := logic.UpdateVerifyConfig(req)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"TurnstileSiteKey":          req.TurnstileSiteKey,
		"TurnstileSecret":           req.TurnstileSecret,
		"EnableLoginVerify":         "true",
		"EnableRegisterVerify":      "true",
		"EnableResetPasswordVerify": "true",
	}, updated)
	require.ElementsMatch(t, []string{serverconfig.VerifyConfigKey, serverconfig.GlobalConfigKey}, deletedKeys)
	require.Equal(t, 1, reloadCount)
	require.Equal(t, req.TurnstileSiteKey, cfg.Verify.TurnstileSiteKey)
	require.Equal(t, req.TurnstileSecret, cfg.Verify.TurnstileSecret)
	require.Equal(t, req.EnableLoginVerify, cfg.Verify.LoginVerify)
	require.Equal(t, req.EnableRegisterVerify, cfg.Verify.RegisterVerify)
	require.Equal(t, req.EnableResetPasswordVerify, cfg.Verify.ResetPasswordVerify)
}

func TestUpdateVerifyConfigReturnsDatabaseUpdateErrorWhenTransactionFails(t *testing.T) {
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
		RunReloadVerify: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateVerifyConfigLogic(context.Background(), deps)

	err := logic.UpdateVerifyConfig(&types.VerifyConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, deleteCalled)
	require.False(t, reloadCalled)
}

func TestUpdateVerifyConfigReturnsDatabaseUpdateErrorWhenCacheEvictionFails(t *testing.T) {
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
		RunReloadVerify: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateVerifyConfigLogic(context.Background(), deps)

	err := logic.UpdateVerifyConfig(&types.VerifyConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, reloadCalled)
}

func TestUpdateVerifyConfigAllowsMissingReloadHook(t *testing.T) {
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
	logic := NewUpdateVerifyConfigLogic(context.Background(), deps)

	err := logic.UpdateVerifyConfig(&types.VerifyConfig{})

	require.NoError(t, err)
}

func TestUpdateVerifyConfigReturnsErrorWhenReloadHookFails(t *testing.T) {
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
		RunReloadVerify: func() error {
			return errors.New("reload failed")
		},
	}
	logic := NewUpdateVerifyConfigLogic(context.Background(), deps)

	err := logic.UpdateVerifyConfig(&types.VerifyConfig{})

	requireSystemCodeError(t, err, xerr.ERROR)
}
