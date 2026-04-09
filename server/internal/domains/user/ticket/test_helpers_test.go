package ticket

import (
	"context"
	"testing"

	modelticket "github.com/perfect-panel/server/internal/platform/persistence/ticket"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/stretchr/testify/require"
)

type fakeTicketModel struct {
	modelticket.Model
	insertFn             func(context.Context, *modelticket.Ticket) error
	findOneFn            func(context.Context, int64) (*modelticket.Ticket, error)
	insertTicketFollowFn func(context.Context, *modelticket.Follow) error
	updateTicketStatusFn func(context.Context, int64, int64, uint8) error
}

func (f fakeTicketModel) Insert(ctx context.Context, data *modelticket.Ticket) error {
	if f.insertFn == nil {
		panic("unexpected Insert call")
	}
	return f.insertFn(ctx, data)
}

func (f fakeTicketModel) FindOne(ctx context.Context, id int64) (*modelticket.Ticket, error) {
	if f.findOneFn == nil {
		panic("unexpected FindOne call")
	}
	return f.findOneFn(ctx, id)
}

func (f fakeTicketModel) InsertTicketFollow(ctx context.Context, data *modelticket.Follow) error {
	if f.insertTicketFollowFn == nil {
		panic("unexpected InsertTicketFollow call")
	}
	return f.insertTicketFollowFn(ctx, data)
}

func (f fakeTicketModel) UpdateTicketStatus(ctx context.Context, id, userID int64, status uint8) error {
	if f.updateTicketStatusFn == nil {
		panic("unexpected UpdateTicketStatus call")
	}
	return f.updateTicketStatusFn(ctx, id, userID, status)
}

func requireTicketCodeError(t *testing.T, err error, want uint32) {
	t.Helper()

	require.Error(t, err)

	var codeErr *xerr.CodeError
	require.ErrorAs(t, err, &codeErr)
	require.Equal(t, want, codeErr.GetErrCode())
}
