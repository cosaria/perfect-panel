package response

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	pkgerrors "github.com/pkg/errors"

	"github.com/perfect-panel/server/internal/platform/support/xerr"
)

const problemContentType = "application/problem+json; charset=utf-8"

var humaProblemFactoryOnce sync.Once

type humaCompatibilityModeKey struct{}

const (
	ProblemTypeUnauthorized    = "urn:perfect-panel:error:unauthorized"
	ProblemTypeForbidden       = "urn:perfect-panel:error:forbidden"
	ProblemTypeInvalidRequest  = "urn:perfect-panel:error:invalid-request"
	ProblemTypeNodeUnavailable = "urn:perfect-panel:error:node-unavailable"
)

// Problem is the canonical RFC 9457 error contract for first-party JSON APIs.
type Problem struct {
	Type     string              `json:"type,omitempty" format:"uri"`
	Title    string              `json:"title,omitempty"`
	Status   int                 `json:"status,omitempty"`
	Detail   string              `json:"detail,omitempty"`
	Instance string              `json:"instance,omitempty" format:"uri"`
	Errors   []*huma.ErrorDetail `json:"errors,omitempty"`
	Code     uint32              `json:"code,omitempty"`
}

type LegacyHumaError struct {
	Status int    `json:"-"`
	Code   uint32 `json:"code"`
	Msg    string `json:"msg"`
}

func (e *LegacyHumaError) Error() string {
	return e.Msg
}

func (e *LegacyHumaError) GetStatus() int {
	return e.Status
}

func (p *Problem) Error() string {
	return p.Detail
}

func (p *Problem) GetStatus() int {
	return p.Status
}

func (p *Problem) ContentType(ct string) string {
	if ct == "application/json" {
		return "application/problem+json"
	}
	if ct == "application/cbor" {
		return "application/problem+cbor"
	}
	return ct
}

// InstallHumaProblemFactory forces huma's framework-generated errors through
// the shared problem contract used by Gin response helpers.
func InstallHumaProblemFactory() {
	humaProblemFactoryOnce.Do(func() {
		huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
			return humaRuntimeError(context.Background(), NewProblemForStatus(status, msg, errs...))
		}
		huma.NewErrorWithContext = func(ctx huma.Context, status int, msg string, errs ...error) huma.StatusError {
			return humaRuntimeError(ctx.Context(), NewProblemForStatus(status, msg, errs...))
		}
	})
}

func WithHumaCompatibilityMode(ctx context.Context, enabled bool) context.Context {
	if ctx == nil || !enabled {
		return ctx
	}
	return context.WithValue(ctx, humaCompatibilityModeKey{}, true)
}

// NewProblemFromError converts an application error into the shared problem
// contract, using xerr when present and falling back to a generic 500.
func NewProblemFromError(err error) *Problem {
	return newProblemFromError(err, 0, "", nil)
}

// NewValidationProblem builds a 422 invalid-params problem with detail
// extensions for request validation errors.
func NewValidationProblem(errs ...error) *Problem {
	return newProblemFromError(xerr.NewErrCode(xerr.InvalidParams), 0, "", errs)
}

// NewProblemForStatus builds a shared problem from Huma/framework errors where
// the HTTP status and top-level message are known in advance.
func NewProblemForStatus(status int, msg string, errs ...error) *Problem {
	return newProblemFromError(nil, status, msg, errs)
}

// NewPublicProblem builds a coarse public problem without leaking business
// codes or internal error taxonomy.
func NewPublicProblem(status int, typeURI string, detail string, errs ...error) *Problem {
	problem := &Problem{
		Type:   typeURI,
		Title:  http.StatusText(status),
		Status: status,
		Detail: strings.TrimSpace(detail),
		Errors: sanitizeErrorDetails(status, errs),
	}

	if problem.Type == "" {
		problem.Type = "about:blank"
	}
	if problem.Detail == "" {
		problem.Detail = http.StatusText(status)
	}
	if len(problem.Errors) == 0 {
		problem.Errors = nil
	}
	return problem
}

// AsHumaStatusError converts any error into a shared problem suitable for
// returning from a huma handler. Existing response problems are preserved.
func AsHumaStatusError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	var problem *Problem
	if stderrors.As(err, &problem) {
		return humaRuntimeError(ctx, problem)
	}

	converted := NewProblemFromError(err)

	var headersErr huma.HeadersError
	if stderrors.As(err, &headersErr) {
		return huma.ErrorWithHeaders(humaRuntimeError(ctx, converted), headersErr.GetHeaders())
	}
	return humaRuntimeError(ctx, converted)
}

// WriteProblem writes the canonical problem payload for Gin handlers.
func WriteProblem(ctx *gin.Context, problem *Problem) {
	if problem == nil {
		problem = NewProblemForStatus(http.StatusInternalServerError, "", nil)
	}

	body, err := json.Marshal(problem)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, problemContentType, []byte(`{"type":"about:blank","title":"Internal Server Error","status":500,"detail":"Internal Server Error"}`))
		return
	}

	ctx.Data(problem.Status, problemContentType, body)
}

func newProblemFromError(err error, status int, msg string, extraErrs []error) *Problem {
	codeErr := codeErrorFrom(err)
	if codeErr != nil {
		status = codeErr.GetStatus()
		msg = codeErr.GetErrMsg()
	}

	if status == 0 {
		status = http.StatusInternalServerError
	}

	problem := &Problem{
		Type:   "about:blank",
		Title:  http.StatusText(status),
		Status: status,
		Detail: sanitizeDetail(status, msg, codeErr),
		Errors: sanitizeErrorDetails(status, extraErrs),
	}

	if codeErr != nil {
		problem.Code = codeErr.GetErrCode()
		problem.Type = problemTypeURI(codeErr.GetErrCode())
	}

	if len(problem.Errors) == 0 {
		problem.Errors = nil
	}

	return problem
}

func codeErrorFrom(err error) *xerr.CodeError {
	if err == nil {
		return nil
	}

	var codeErr *xerr.CodeError
	if stderrors.As(err, &codeErr) {
		return codeErr
	}

	if cause := pkgerrors.Cause(err); cause != nil && cause != err {
		if stderrors.As(cause, &codeErr) {
			return codeErr
		}
	}
	return nil
}

func sanitizeDetail(status int, msg string, codeErr *xerr.CodeError) string {
	if codeErr != nil {
		return codeErr.GetErrMsg()
	}

	trimmed := strings.TrimSpace(msg)
	switch {
	case trimmed == "":
		return http.StatusText(status)
	case status >= 500:
		return http.StatusText(status)
	default:
		return trimmed
	}
}

func sanitizeErrorDetails(status int, errs []error) []*huma.ErrorDetail {
	if len(errs) == 0 || status >= 500 {
		return nil
	}

	details := make([]*huma.ErrorDetail, 0, len(errs))
	for _, err := range errs {
		if err == nil {
			continue
		}

		if converted, ok := err.(huma.ErrorDetailer); ok {
			detail := converted.ErrorDetail()
			if detail == nil {
				continue
			}
			copied := *detail
			copied.Message = strings.TrimSpace(copied.Message)
			if copied.Message == "" {
				continue
			}
			details = append(details, &copied)
			continue
		}

		message := strings.TrimSpace(err.Error())
		if message == "" {
			continue
		}
		details = append(details, &huma.ErrorDetail{Message: message})
	}

	return details
}

func problemTypeURI(code uint32) string {
	return "urn:perfect-panel:error:" + strconv.FormatUint(uint64(code), 10)
}

func humaRuntimeError(ctx context.Context, problem *Problem) huma.StatusError {
	if humaCompatibilityModeEnabled(ctx) {
		return compatibilityEnvelope(problem)
	}
	return problem
}

func humaCompatibilityModeEnabled(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	enabled, _ := ctx.Value(humaCompatibilityModeKey{}).(bool)
	return enabled
}

func compatibilityEnvelope(problem *Problem) *LegacyHumaError {
	code := problem.Code
	switch {
	case code != 0:
	case problem.Status == http.StatusUnprocessableEntity:
		code = xerr.InvalidParams
	case problem.Status == http.StatusTooManyRequests:
		code = xerr.TooManyRequests
	default:
		code = xerr.ERROR
	}

	msg := strings.TrimSpace(problem.Detail)
	if msg == "" {
		msg = http.StatusText(problem.Status)
	}

	return &LegacyHumaError{
		Status: http.StatusOK,
		Code:   code,
		Msg:    msg,
	}
}
