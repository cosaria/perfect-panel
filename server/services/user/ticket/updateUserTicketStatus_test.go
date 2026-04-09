package ticket

import (
	"context"
	"errors"
	"testing"

	"github.com/perfect-panel/server/config"
	modelticket "github.com/perfect-panel/server/models/ticket"
	modeluser "github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserTicketStatusReturnsInvalidAccessWhenUserMissing(t *testing.T) {
	status := uint8(modelticket.Closed)
	logic := NewUpdateUserTicketStatusLogic(context.Background(), Deps{})

	err := logic.UpdateUserTicketStatus(&types.UpdateUserTicketStatusRequest{
		Id:     1,
		Status: &status,
	})

	requireTicketCodeError(t, err, xerr.InvalidAccess)
}

func TestUpdateUserTicketStatusRejectsInvalidParams(t *testing.T) {
	validStatus := uint8(modelticket.Closed)
	invalidStatus := uint8(modelticket.Waiting)
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 31})
	logic := NewUpdateUserTicketStatusLogic(ctx, Deps{})

	testCases := []struct {
		name string
		req  types.UpdateUserTicketStatusRequest
	}{
		{
			name: "missing status",
			req: types.UpdateUserTicketStatusRequest{
				Id: 3,
			},
		},
		{
			name: "invalid ticket id",
			req: types.UpdateUserTicketStatusRequest{
				Id:     0,
				Status: &validStatus,
			},
		},
		{
			name: "user cannot set waiting",
			req: types.UpdateUserTicketStatusRequest{
				Id:     3,
				Status: &invalidStatus,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := logic.UpdateUserTicketStatus(&tc.req)

			requireTicketCodeError(t, err, xerr.InvalidParams)
		})
	}
}

func TestUpdateUserTicketStatusDelegatesClosedStatusToTicketModel(t *testing.T) {
	status := uint8(modelticket.Closed)
	var gotID int64
	var gotUserID int64
	var gotStatus uint8
	deps := Deps{
		TicketModel: fakeTicketModel{
			updateTicketStatusFn: func(_ context.Context, id, userID int64, newStatus uint8) error {
				gotID = id
				gotUserID = userID
				gotStatus = newStatus
				return nil
			},
		},
	}
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 32})
	logic := NewUpdateUserTicketStatusLogic(ctx, deps)

	err := logic.UpdateUserTicketStatus(&types.UpdateUserTicketStatusRequest{
		Id:     88,
		Status: &status,
	})

	require.NoError(t, err)
	require.Equal(t, int64(88), gotID)
	require.Equal(t, int64(32), gotUserID)
	require.Equal(t, uint8(modelticket.Closed), gotStatus)
}

func TestUpdateUserTicketStatusReturnsDatabaseUpdateErrorWhenModelFails(t *testing.T) {
	status := uint8(modelticket.Closed)
	deps := Deps{
		TicketModel: fakeTicketModel{
			updateTicketStatusFn: func(context.Context, int64, int64, uint8) error {
				return errors.New("update failed")
			},
		},
	}
	ctx := context.WithValue(context.Background(), config.CtxKeyUser, &modeluser.User{Id: 33})
	logic := NewUpdateUserTicketStatusLogic(ctx, deps)

	err := logic.UpdateUserTicketStatus(&types.UpdateUserTicketStatusRequest{
		Id:     99,
		Status: &status,
	})

	requireTicketCodeError(t, err, xerr.DatabaseUpdateError)
}
