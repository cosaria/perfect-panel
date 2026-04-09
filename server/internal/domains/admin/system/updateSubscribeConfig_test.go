package system

import (
	"context"
	"errors"
	"testing"
	"time"

	serverconfig "github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUpdateSubscribeConfigUpdatesAllFieldsEvictsCacheAndReloadsWhenPathUnchanged(t *testing.T) {
	updated := map[string]string{}
	var deletedKeys []string
	reloadCount := 0
	restartCalled := false
	req := &types.SubscribeConfig{
		SingleModel:     true,
		SubscribePath:   "/sub",
		SubscribeDomain: "example.com",
		PanDomain:       true,
		UserAgentLimit:  true,
		UserAgentList:   "clash,sing-box",
	}
	deps := Deps{
		Config: &serverconfig.Config{
			Subscribe: serverconfig.SubscribeConfig{SubscribePath: "/sub"},
		},
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(_ context.Context, _ *gorm.DB, category, fieldName, fieldValue string) error {
			require.Equal(t, "subscribe", category)
			updated[fieldName] = fieldValue
			return nil
		},
		DeleteCacheKeys: func(_ context.Context, keys ...string) error {
			deletedKeys = append(deletedKeys, keys...)
			return nil
		},
		RunReloadSubscribe: func() error {
			reloadCount++
			return nil
		},
		Restart: func() error {
			restartCalled = true
			return nil
		},
	}
	logic := NewUpdateSubscribeConfigLogic(context.Background(), deps)

	err := logic.UpdateSubscribeConfig(req)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"SingleModel":     "true",
		"SubscribePath":   req.SubscribePath,
		"SubscribeDomain": req.SubscribeDomain,
		"PanDomain":       "true",
		"UserAgentLimit":  "true",
		"UserAgentList":   req.UserAgentList,
	}, updated)
	require.ElementsMatch(t, []string{serverconfig.SubscribeConfigKey, serverconfig.GlobalConfigKey}, deletedKeys)
	require.Equal(t, 1, reloadCount)
	require.False(t, restartCalled)
}

func TestUpdateSubscribeConfigSchedulesRestartWhenPathChanges(t *testing.T) {
	restarted := make(chan struct{}, 1)
	reloadCount := 0
	deps := Deps{
		Config: &serverconfig.Config{
			Subscribe: serverconfig.SubscribeConfig{SubscribePath: "/old"},
		},
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
		RunReloadSubscribe: func() error {
			reloadCount++
			return nil
		},
		Restart: func() error {
			restarted <- struct{}{}
			return nil
		},
	}
	logic := NewUpdateSubscribeConfigLogic(context.Background(), deps)

	err := logic.UpdateSubscribeConfig(&types.SubscribeConfig{SubscribePath: "/new"})

	require.NoError(t, err)
	select {
	case <-restarted:
	case <-time.After(time.Second):
		t.Fatal("expected restart to be scheduled")
	}
	require.Zero(t, reloadCount)
}

func TestUpdateSubscribeConfigReturnsDatabaseUpdateErrorWhenTransactionFails(t *testing.T) {
	reloadCalled := false
	deleteCalled := false
	restartCalled := false
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
		RunReloadSubscribe: func() error {
			reloadCalled = true
			return nil
		},
		Restart: func() error {
			restartCalled = true
			return nil
		},
	}
	logic := NewUpdateSubscribeConfigLogic(context.Background(), deps)

	err := logic.UpdateSubscribeConfig(&types.SubscribeConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, deleteCalled)
	require.False(t, reloadCalled)
	require.False(t, restartCalled)
}

func TestUpdateSubscribeConfigReturnsDatabaseUpdateErrorWhenCacheEvictionFails(t *testing.T) {
	reloadCalled := false
	deps := Deps{
		Config: &serverconfig.Config{
			Subscribe: serverconfig.SubscribeConfig{SubscribePath: "/sub"},
		},
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
		RunReloadSubscribe: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateSubscribeConfigLogic(context.Background(), deps)

	err := logic.UpdateSubscribeConfig(&types.SubscribeConfig{SubscribePath: "/sub"})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, reloadCalled)
}

func TestUpdateSubscribeConfigAllowsMissingReloadHook(t *testing.T) {
	deps := Deps{
		Config: &serverconfig.Config{
			Subscribe: serverconfig.SubscribeConfig{SubscribePath: "/sub"},
		},
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
	logic := NewUpdateSubscribeConfigLogic(context.Background(), deps)

	err := logic.UpdateSubscribeConfig(&types.SubscribeConfig{SubscribePath: "/sub"})

	require.NoError(t, err)
}

func TestUpdateSubscribeConfigReturnsErrorWhenReloadHookFails(t *testing.T) {
	deps := Deps{
		Config: &serverconfig.Config{
			Subscribe: serverconfig.SubscribeConfig{SubscribePath: "/sub"},
		},
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
		RunReloadSubscribe: func() error {
			return errors.New("reload failed")
		},
	}
	logic := NewUpdateSubscribeConfigLogic(context.Background(), deps)

	err := logic.UpdateSubscribeConfig(&types.SubscribeConfig{SubscribePath: "/sub"})

	requireSystemCodeError(t, err, xerr.ERROR)
}
