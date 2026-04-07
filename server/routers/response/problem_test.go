package response

import (
	"errors"
	"net/http"
	"testing"

	"github.com/perfect-panel/server/modules/infra/xerr"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewProblemFromErrorUsesXerrStatusAndBusinessCode(t *testing.T) {
	problem := NewProblemFromError(pkgerrors.Wrap(xerr.NewErrCode(xerr.InvalidAccess), "internal detail"))

	require.Equal(t, http.StatusForbidden, problem.Status)
	require.Equal(t, "Forbidden", problem.Title)
	require.Equal(t, "Invalid access", problem.Detail)
	require.Equal(t, uint32(xerr.InvalidAccess), problem.Code)
	require.Equal(t, "urn:perfect-panel:error:40005", problem.Type)
	require.Nil(t, problem.Errors)
}

func TestNewProblemFromErrorRedactsGenericInternalErrors(t *testing.T) {
	problem := NewProblemFromError(errors.New("dial tcp 10.0.0.7:6379: connect: connection refused"))

	require.Equal(t, http.StatusInternalServerError, problem.Status)
	require.Equal(t, "Internal Server Error", problem.Title)
	require.Equal(t, "Internal Server Error", problem.Detail)
	require.Zero(t, problem.Code)
	require.Equal(t, "about:blank", problem.Type)
	require.Nil(t, problem.Errors)
}

func TestNewValidationProblemPreservesValidationDetails(t *testing.T) {
	problem := NewValidationProblem(errors.New("email is required"))

	require.Equal(t, http.StatusUnprocessableEntity, problem.Status)
	require.Equal(t, "Param Error", problem.Detail)
	require.Equal(t, uint32(xerr.InvalidParams), problem.Code)
	require.Len(t, problem.Errors, 1)
	require.Equal(t, "email is required", problem.Errors[0].Message)
}
