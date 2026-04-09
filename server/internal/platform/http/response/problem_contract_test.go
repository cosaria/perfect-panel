package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func decodeProblemBody(t *testing.T, body *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(body.Body.Bytes(), &payload))
	return payload
}

func TestPhase5HttpResultUsesProblemDetailsForCodeErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	HttpResult(ctx, nil, pkgerrors.Wrap(xerr.NewErrCode(xerr.InvalidAccess), "wrapped detail that should stay internal"))

	require.Equal(t, http.StatusForbidden, recorder.Code)
	require.Contains(t, recorder.Header().Get("Content-Type"), "application/problem+json")

	body := decodeProblemBody(t, recorder)
	require.EqualValues(t, http.StatusForbidden, body["status"])
	require.Equal(t, "Forbidden", body["title"])
	require.Equal(t, "Invalid access", body["detail"])
	require.EqualValues(t, xerr.InvalidAccess, body["code"])
	require.Equal(t, "urn:perfect-panel:error:40005", body["type"])
}

func TestPhase5ParamErrorResultUsesProblemDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)

	ParamErrorResult(ctx, errors.New("email is required"))

	require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	require.Contains(t, recorder.Header().Get("Content-Type"), "application/problem+json")

	body := decodeProblemBody(t, recorder)
	require.EqualValues(t, http.StatusUnprocessableEntity, body["status"])
	require.Equal(t, "Unprocessable Entity", body["title"])
	require.Equal(t, "Param Error", body["detail"])
	require.EqualValues(t, xerr.InvalidParams, body["code"])
	require.Equal(t, "urn:perfect-panel:error:400", body["type"])

	errorsList, ok := body["errors"].([]interface{})
	require.True(t, ok, "expected validation details array")
	require.NotEmpty(t, errorsList)

	first, ok := errorsList[0].(map[string]interface{})
	require.True(t, ok, "expected first validation detail object")
	require.Equal(t, "email is required", first["message"])
}
