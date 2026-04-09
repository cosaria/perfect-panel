package ticket

import (
	"context"
	"errors"
	"testing"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	modelticket "github.com/perfect-panel/server/internal/platform/persistence/ticket"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/stretchr/testify/require"
)

func TestCreateUserTicketReturnsInvalidAccessWhenUserMissing(t *testing.T) {
	logic := NewCreateUserTicketLogic(context.Background(), Deps{})

	err := logic.CreateUserTicket(&types.CreateUserTicketRequest{
		Title:       "Need help",
		Description: "Please assist",
	})

	requireTicketCodeError(t, err, xerr.InvalidAccess)
}

func TestCreateUserTicketRejectsBlankTitleOrDescription(t *testing.T) {
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 12})
	logic := NewCreateUserTicketLogic(ctx, Deps{})

	err := logic.CreateUserTicket(&types.CreateUserTicketRequest{
		Title:       "   ",
		Description: "\n\t",
	})

	requireTicketCodeError(t, err, xerr.InvalidParams)
}

func TestCreateUserTicketInsertsPendingTicketForCurrentUser(t *testing.T) {
	var inserted *modelticket.Ticket
	deps := Deps{
		TicketModel: fakeTicketModel{
			insertFn: func(_ context.Context, data *modelticket.Ticket) error {
				copied := *data
				inserted = &copied
				return nil
			},
		},
	}
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 19})
	logic := NewCreateUserTicketLogic(ctx, deps)

	err := logic.CreateUserTicket(&types.CreateUserTicketRequest{
		Title:       "  API failure  ",
		Description: "  Need logs  ",
	})

	require.NoError(t, err)
	require.NotNil(t, inserted)
	require.Equal(t, int64(19), inserted.UserId)
	require.Equal(t, uint8(modelticket.Pending), inserted.Status)
	require.Equal(t, "API failure", inserted.Title)
	require.Equal(t, "Need logs", inserted.Description)
}

func TestCreateUserTicketReturnsDatabaseInsertErrorWhenInsertFails(t *testing.T) {
	deps := Deps{
		TicketModel: fakeTicketModel{
			insertFn: func(context.Context, *modelticket.Ticket) error {
				return errors.New("insert failed")
			},
		},
	}
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 21})
	logic := NewCreateUserTicketLogic(ctx, deps)

	err := logic.CreateUserTicket(&types.CreateUserTicketRequest{
		Title:       "Need help",
		Description: "Please assist",
	})

	requireTicketCodeError(t, err, xerr.DatabaseInsertError)
}
