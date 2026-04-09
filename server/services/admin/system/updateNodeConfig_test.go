package system

import (
	"context"
	"errors"
	"testing"

	serverconfig "github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUpdateNodeConfigUpdatesAllFieldsEvictsCacheAndReloads(t *testing.T) {
	updated := map[string]string{}
	var deletedKeys []string
	reloadCount := 0
	req := &types.NodeConfig{
		NodeSecret:             "node-secret",
		NodePullInterval:       30,
		NodePushInterval:       45,
		TrafficReportThreshold: 1024,
		IPStrategy:             "prefer_ipv4",
		DNS: []types.NodeDNS{
			{Proto: "udp", Address: "1.1.1.1", Domains: []string{"example.com"}},
		},
		Block: []string{"ads.example.com"},
		Outbound: []types.NodeOutbound{
			{Name: "proxy", Protocol: "socks5", Address: "127.0.0.1", Port: 1080, Password: "secret", Rules: []string{"geoip:private"}},
		},
	}
	deps := Deps{
		Config: &serverconfig.Config{},
		SystemModel: fakeSystemModel{
			transactionFn: func(_ context.Context, fn func(*gorm.DB) error) error {
				return fn(nil)
			},
		},
		UpdateConfigField: func(_ context.Context, _ *gorm.DB, category, fieldName, fieldValue string) error {
			require.Equal(t, "server", category)
			updated[fieldName] = fieldValue
			return nil
		},
		DeleteCacheKeys: func(_ context.Context, keys ...string) error {
			deletedKeys = append(deletedKeys, keys...)
			return nil
		},
		RunReloadNode: func() error {
			reloadCount++
			return nil
		},
	}
	logic := NewUpdateNodeConfigLogic(context.Background(), deps)

	err := logic.UpdateNodeConfig(req)

	require.NoError(t, err)
	require.Equal(t, req.NodeSecret, updated["NodeSecret"])
	require.Equal(t, "30", updated["NodePullInterval"])
	require.Equal(t, "45", updated["NodePushInterval"])
	require.Equal(t, "1024", updated["TrafficReportThreshold"])
	require.Equal(t, req.IPStrategy, updated["IPStrategy"])
	require.JSONEq(t, `[{"proto":"udp","address":"1.1.1.1","domains":["example.com"]}]`, updated["DNS"])
	require.JSONEq(t, `["ads.example.com"]`, updated["Block"])
	require.JSONEq(t, `[{"name":"proxy","protocol":"socks5","address":"127.0.0.1","port":1080,"password":"secret","rules":["geoip:private"]}]`, updated["Outbound"])
	require.ElementsMatch(t, []string{serverconfig.NodeConfigKey}, deletedKeys)
	require.Equal(t, 1, reloadCount)
}

func TestUpdateNodeConfigReturnsDatabaseUpdateErrorWhenTransactionFails(t *testing.T) {
	deleteCalled := false
	reloadCalled := false
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
		RunReloadNode: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateNodeConfigLogic(context.Background(), deps)

	err := logic.UpdateNodeConfig(&types.NodeConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, deleteCalled)
	require.False(t, reloadCalled)
}

func TestUpdateNodeConfigReturnsDatabaseUpdateErrorWhenCacheEvictionFails(t *testing.T) {
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
		RunReloadNode: func() error {
			reloadCalled = true
			return nil
		},
	}
	logic := NewUpdateNodeConfigLogic(context.Background(), deps)

	err := logic.UpdateNodeConfig(&types.NodeConfig{})

	requireSystemCodeError(t, err, xerr.DatabaseUpdateError)
	require.False(t, reloadCalled)
}

func TestUpdateNodeConfigAllowsMissingReloadHook(t *testing.T) {
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
	logic := NewUpdateNodeConfigLogic(context.Background(), deps)

	err := logic.UpdateNodeConfig(&types.NodeConfig{})

	require.NoError(t, err)
}

func TestUpdateNodeConfigReturnsErrorWhenReloadHookFails(t *testing.T) {
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
		RunReloadNode: func() error {
			return errors.New("reload failed")
		},
	}
	logic := NewUpdateNodeConfigLogic(context.Background(), deps)

	err := logic.UpdateNodeConfig(&types.NodeConfig{})

	requireSystemCodeError(t, err, xerr.ERROR)
}
