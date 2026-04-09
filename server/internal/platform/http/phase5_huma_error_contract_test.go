package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/internal/platform/http/response"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type phase5ValidationInput struct {
	Body struct {
		Email string `json:"email" validate:"required"`
	}
}

type phase5ValidationOutput struct {
	Body map[string]string
}

func decodeHumaProblemBody(t *testing.T, recorder *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &payload))
	return payload
}

func TestPhase5HumaHandlerCodeErrorsUseProblemDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.InstallHumaProblemFactory()

	router := gin.New()
	group := router.Group("/test")
	api := humagin.NewWithGroup(router, group, apiConfig("Phase5 Test API", "1.0.0"))

	registerOperation(api, huma.Operation{
		OperationID: "phase5TestHumaFailure",
		Method:      http.MethodGet,
		Path:        "/fail",
		Summary:     "Phase 5 Huma failure",
		Tags:        []string{"phase5"},
	}, func(context.Context, *struct{}) (*struct{}, error) {
		return nil, pkgerrors.Wrap(xerr.NewErrCode(xerr.UserNotExist), "database lookup failed")
	})

	req := httptest.NewRequest(http.MethodGet, "/test/fail", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Header().Get("Content-Type"), "application/problem+json")

	body := decodeHumaProblemBody(t, resp)
	require.EqualValues(t, http.StatusBadRequest, body["status"])
	require.Equal(t, "Bad Request", body["title"])
	require.Equal(t, "User does not exist", body["detail"])
	require.EqualValues(t, xerr.UserNotExist, body["code"])
	require.Equal(t, "urn:perfect-panel:error:20002", body["type"])
}

func TestPhase5HumaValidationErrorsUseProblemDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.InstallHumaProblemFactory()

	router := gin.New()
	group := router.Group("/test")
	api := humagin.NewWithGroup(router, group, apiConfig("Phase5 Test API", "1.0.0"))

	registerOperation(api, huma.Operation{
		OperationID: "phase5TestValidation",
		Method:      http.MethodPost,
		Path:        "/validate",
		Summary:     "Phase 5 validation",
		Tags:        []string{"phase5"},
	}, func(_ context.Context, input *phase5ValidationInput) (*phase5ValidationOutput, error) {
		return &phase5ValidationOutput{Body: map[string]string{"email": input.Body.Email}}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/test/validate", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	require.Contains(t, resp.Header().Get("Content-Type"), "application/problem+json")

	body := decodeHumaProblemBody(t, resp)
	require.EqualValues(t, http.StatusUnprocessableEntity, body["status"])
	require.Equal(t, "Unprocessable Entity", body["title"])
	require.Equal(t, "validation failed", body["detail"])

	errorsList, ok := body["errors"].([]interface{})
	require.True(t, ok, "expected validation details array")
	require.NotEmpty(t, errorsList)
}

func TestPhase5HumaCompatibilityModeReenvelopesRuntimeErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.InstallHumaProblemFactory()

	router := gin.New()
	group := router.Group("/test")
	api := humagin.NewWithGroup(router, group, apiConfig("Phase5 Compat API", "1.0.0"))
	configureHumaAPI(api, true)

	registerOperation(api, huma.Operation{
		OperationID: "phase5CompatFailure",
		Method:      http.MethodGet,
		Path:        "/compat-fail",
		Summary:     "Phase 5 compatibility failure",
		Tags:        []string{"phase5"},
	}, func(context.Context, *struct{}) (*struct{}, error) {
		return nil, pkgerrors.Wrap(xerr.NewErrCode(xerr.UserNotExist), "database lookup failed")
	})

	req := httptest.NewRequest(http.MethodGet, "/test/compat-fail", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Header().Get("Content-Type"), "application/json")

	body := decodeHumaProblemBody(t, resp)
	require.EqualValues(t, xerr.UserNotExist, body["code"])
	require.Equal(t, "User does not exist", body["msg"])
	require.NotContains(t, body, "status")
	require.NotContains(t, body, "detail")
}

func TestPhase5HumaCompatibilityModeIsScopedPerAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.InstallHumaProblemFactory()

	router := gin.New()

	compatGroup := router.Group("/compat")
	compatAPI := humagin.NewWithGroup(router, compatGroup, apiConfig("Compat API", "1.0.0"))
	configureHumaAPI(compatAPI, true)
	registerOperation(compatAPI, huma.Operation{
		OperationID: "compatOnlyFailure",
		Method:      http.MethodGet,
		Path:        "/fail",
		Summary:     "Compat scoped failure",
		Tags:        []string{"phase5"},
	}, func(context.Context, *struct{}) (*struct{}, error) {
		return nil, pkgerrors.Wrap(xerr.NewErrCode(xerr.UserNotExist), "database lookup failed")
	})

	strictGroup := router.Group("/strict")
	strictAPI := humagin.NewWithGroup(router, strictGroup, apiConfig("Strict API", "1.0.0"))
	configureHumaAPI(strictAPI, false)
	registerOperation(strictAPI, huma.Operation{
		OperationID: "strictOnlyFailure",
		Method:      http.MethodGet,
		Path:        "/fail",
		Summary:     "Strict scoped failure",
		Tags:        []string{"phase5"},
	}, func(context.Context, *struct{}) (*struct{}, error) {
		return nil, pkgerrors.Wrap(xerr.NewErrCode(xerr.UserNotExist), "database lookup failed")
	})

	compatReq := httptest.NewRequest(http.MethodGet, "/compat/fail", nil)
	compatResp := httptest.NewRecorder()
	router.ServeHTTP(compatResp, compatReq)
	require.Equal(t, http.StatusOK, compatResp.Code)
	compatBody := decodeHumaProblemBody(t, compatResp)
	require.EqualValues(t, xerr.UserNotExist, compatBody["code"])
	require.Equal(t, "User does not exist", compatBody["msg"])
	require.NotContains(t, compatBody, "status")

	strictReq := httptest.NewRequest(http.MethodGet, "/strict/fail", nil)
	strictResp := httptest.NewRecorder()
	router.ServeHTTP(strictResp, strictReq)
	require.Equal(t, http.StatusBadRequest, strictResp.Code)
	strictBody := decodeHumaProblemBody(t, strictResp)
	require.EqualValues(t, http.StatusBadRequest, strictBody["status"])
	require.Equal(t, "User does not exist", strictBody["detail"])
	require.NotContains(t, strictBody, "msg")
}

func TestPhase5HumaCompatibilityModeAlsoAppliesToFrameworkValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response.InstallHumaProblemFactory()

	router := gin.New()
	group := router.Group("/compat")
	api := humagin.NewWithGroup(router, group, apiConfig("Compat Validation API", "1.0.0"))
	configureHumaAPI(api, true)

	registerOperation(api, huma.Operation{
		OperationID: "compatValidationFailure",
		Method:      http.MethodPost,
		Path:        "/validate",
		Summary:     "Compat validation failure",
		Tags:        []string{"phase5"},
	}, func(_ context.Context, input *phase5ValidationInput) (*phase5ValidationOutput, error) {
		return &phase5ValidationOutput{Body: map[string]string{"email": input.Body.Email}}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/compat/validate", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	body := decodeHumaProblemBody(t, resp)
	require.EqualValues(t, xerr.InvalidParams, body["code"])
	require.Equal(t, "validation failed", body["msg"])
	require.NotContains(t, body, "status")
	require.NotContains(t, body, "detail")
}
