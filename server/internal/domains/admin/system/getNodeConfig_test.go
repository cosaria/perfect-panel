package system

import (
	"context"
	"testing"

	modelsystem "github.com/perfect-panel/server/models/system"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/stretchr/testify/require"
)

func TestGetNodeConfigReturnsErrorInsteadOfPanicWhenDNSJSONIsInvalid(t *testing.T) {
	deps := Deps{
		SystemModel: fakeSystemModel{
			getNodeConfigFn: func(context.Context) ([]*modelsystem.System, error) {
				return []*modelsystem.System{
					{Key: "DNS", Type: "string", Value: "{invalid-json"},
				}, nil
			},
		},
	}
	logic := NewGetNodeConfigLogic(context.Background(), deps)

	require.NotPanics(t, func() {
		_, err := logic.GetNodeConfig()
		requireSystemCodeError(t, err, xerr.ERROR)
		require.ErrorContains(t, err, "Unmarshal DNS config error")
	})
}

func TestGetNodeConfigReturnsErrorInsteadOfPanicWhenOutboundJSONIsInvalid(t *testing.T) {
	deps := Deps{
		SystemModel: fakeSystemModel{
			getNodeConfigFn: func(context.Context) ([]*modelsystem.System, error) {
				return []*modelsystem.System{
					{Key: "Outbound", Type: "string", Value: "{invalid-json"},
				}, nil
			},
		},
	}
	logic := NewGetNodeConfigLogic(context.Background(), deps)

	require.NotPanics(t, func() {
		_, err := logic.GetNodeConfig()
		requireSystemCodeError(t, err, xerr.ERROR)
		require.ErrorContains(t, err, "Unmarshal Outbound config error")
	})
}
