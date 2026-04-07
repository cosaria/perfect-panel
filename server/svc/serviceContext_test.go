package svc

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/perfect-panel/server/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestClearSendCountLimitKeysOnlyRemovesLimiterKeys(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	t.Cleanup(func() {
		_ = client.Close()
	})

	ctx := context.Background()
	require.NoError(t, client.Set(ctx, config.SendCountLimitKeyPrefix+"email:user@example.com", "1", 0).Err())
	require.NoError(t, client.Set(ctx, config.SendCountLimitKeyPrefix+"sms:13800000000", "1", 0).Err())
	require.NoError(t, client.Set(ctx, "session:user:1", "keep", 0).Err())

	require.NoError(t, clearSendCountLimitKeys(ctx, client))

	_, err := client.Get(ctx, config.SendCountLimitKeyPrefix+"email:user@example.com").Result()
	require.ErrorIs(t, err, redis.Nil)

	_, err = client.Get(ctx, config.SendCountLimitKeyPrefix+"sms:13800000000").Result()
	require.ErrorIs(t, err, redis.Nil)

	value, err := client.Get(ctx, "session:user:1").Result()
	require.NoError(t, err)
	require.Equal(t, "keep", value)
}
