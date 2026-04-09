package ticket

import (
	"context"
	"errors"
	"testing"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	modelticket "github.com/perfect-panel/server/models/ticket"
	modeluser "github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/stretchr/testify/require"
)

func TestCreateUserTicketFollowReturnsInvalidAccessWhenUserMissing(t *testing.T) {
	logic := NewCreateUserTicketFollowLogic(context.Background(), Deps{})

	err := logic.CreateUserTicketFollow(&types.CreateUserTicketFollowRequest{
		TicketId: 1,
		From:     "User",
		Type:     1,
		Content:  "hello",
	})

	requireTicketCodeError(t, err, xerr.InvalidAccess)
}

func TestCreateUserTicketFollowRejectsInvalidParams(t *testing.T) {
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 41})
	logic := NewCreateUserTicketFollowLogic(ctx, Deps{})

	testCases := []struct {
		name string
		req  types.CreateUserTicketFollowRequest
	}{
		{
			name: "invalid ticket id",
			req: types.CreateUserTicketFollowRequest{
				TicketId: 0,
				From:     "User",
				Type:     1,
				Content:  "hello",
			},
		},
		{
			name: "blank text content",
			req: types.CreateUserTicketFollowRequest{
				TicketId: 3,
				From:     "User",
				Type:     1,
				Content:  " \n\t ",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := logic.CreateUserTicketFollow(&tc.req)

			requireTicketCodeError(t, err, xerr.InvalidParams)
		})
	}
}

func TestCreateUserTicketFollowRejectsAccessToAnotherUsersTicket(t *testing.T) {
	deps := Deps{
		TicketModel: fakeTicketModel{
			findOneFn: func(context.Context, int64) (*modelticket.Ticket, error) {
				return &modelticket.Ticket{Id: 8, UserId: 999}, nil
			},
		},
	}
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 42})
	logic := NewCreateUserTicketFollowLogic(ctx, deps)

	err := logic.CreateUserTicketFollow(&types.CreateUserTicketFollowRequest{
		TicketId: 8,
		From:     "User",
		Type:     1,
		Content:  "hello",
	})

	requireTicketCodeError(t, err, xerr.InvalidAccess)
}

func TestCreateUserTicketFollowForcesUserOriginAndResetsStatusToPending(t *testing.T) {
	var inserted *modelticket.Follow
	var updatedID int64
	var updatedUserID int64
	var updatedStatus uint8
	deps := Deps{
		TicketModel: fakeTicketModel{
			findOneFn: func(context.Context, int64) (*modelticket.Ticket, error) {
				return &modelticket.Ticket{Id: 9, UserId: 43}, nil
			},
			insertTicketFollowFn: func(_ context.Context, data *modelticket.Follow) error {
				copied := *data
				inserted = &copied
				return nil
			},
			updateTicketStatusFn: func(_ context.Context, id, userID int64, status uint8) error {
				updatedID = id
				updatedUserID = userID
				updatedStatus = status
				return nil
			},
		},
	}
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 43})
	logic := NewCreateUserTicketFollowLogic(ctx, deps)

	err := logic.CreateUserTicketFollow(&types.CreateUserTicketFollowRequest{
		TicketId: 9,
		From:     "System",
		Type:     1,
		Content:  "  hello support  ",
	})

	require.NoError(t, err)
	require.NotNil(t, inserted)
	require.Equal(t, int64(9), inserted.TicketId)
	require.Equal(t, "User", inserted.From)
	require.Equal(t, uint8(1), inserted.Type)
	require.Equal(t, "hello support", inserted.Content)
	require.Equal(t, int64(9), updatedID)
	require.Equal(t, int64(43), updatedUserID)
	require.Equal(t, uint8(modelticket.Pending), updatedStatus)
}

func TestCreateUserTicketFollowReturnsDatabaseErrors(t *testing.T) {
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 44})

	t.Run("find ticket", func(t *testing.T) {
		logic := NewCreateUserTicketFollowLogic(ctx, Deps{
			TicketModel: fakeTicketModel{
				findOneFn: func(context.Context, int64) (*modelticket.Ticket, error) {
					return nil, errors.New("lookup failed")
				},
			},
		})

		err := logic.CreateUserTicketFollow(&types.CreateUserTicketFollowRequest{
			TicketId: 1,
			From:     "User",
			Type:     1,
			Content:  "hello",
		})

		requireTicketCodeError(t, err, xerr.DatabaseQueryError)
	})

	t.Run("insert follow", func(t *testing.T) {
		logic := NewCreateUserTicketFollowLogic(ctx, Deps{
			TicketModel: fakeTicketModel{
				findOneFn: func(context.Context, int64) (*modelticket.Ticket, error) {
					return &modelticket.Ticket{Id: 1, UserId: 44}, nil
				},
				insertTicketFollowFn: func(context.Context, *modelticket.Follow) error {
					return errors.New("insert failed")
				},
			},
		})

		err := logic.CreateUserTicketFollow(&types.CreateUserTicketFollowRequest{
			TicketId: 1,
			From:     "User",
			Type:     1,
			Content:  "hello",
		})

		requireTicketCodeError(t, err, xerr.DatabaseInsertError)
	})

	t.Run("update ticket status", func(t *testing.T) {
		logic := NewCreateUserTicketFollowLogic(ctx, Deps{
			TicketModel: fakeTicketModel{
				findOneFn: func(context.Context, int64) (*modelticket.Ticket, error) {
					return &modelticket.Ticket{Id: 1, UserId: 44}, nil
				},
				insertTicketFollowFn: func(context.Context, *modelticket.Follow) error {
					return nil
				},
				updateTicketStatusFn: func(context.Context, int64, int64, uint8) error {
					return errors.New("update failed")
				},
			},
		})

		err := logic.CreateUserTicketFollow(&types.CreateUserTicketFollowRequest{
			TicketId: 1,
			From:     "User",
			Type:     1,
			Content:  "hello",
		})

		requireTicketCodeError(t, err, xerr.DatabaseUpdateError)
	})
}
