package xerr

import (
	"errors"
	"net/http"
)

/**
General common fixed error
*/

type CodeError struct {
	errCode uint32
	errMsg  string
}

var ErrStatusNotModified = errors.New("304 Not Modified")

// GetErrCode returns the error code displayed to the front end
func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

// GetErrMsg returns the error message displayed to the front end
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

// Error implements the error interface
func (e *CodeError) Error() string {
	return e.errMsg
}

// GetStatus maps business error codes to HTTP status codes.
// Implements huma.StatusError interface so huma automatically
// uses the correct HTTP status code in responses.
func (e *CodeError) GetStatus() int {
	switch {
	// Auth errors → 401 Unauthorized
	case e.errCode == ErrorTokenEmpty,
		e.errCode == ErrorTokenInvalid,
		e.errCode == ErrorTokenExpire:
		return http.StatusUnauthorized
	// Access errors → 403 Forbidden
	case e.errCode == InvalidAccess:
		return http.StatusForbidden
	// Rate limiting → 429 Too Many Requests
	case e.errCode == TooManyRequests:
		return http.StatusTooManyRequests
	// Param validation → 422
	case e.errCode == InvalidParams:
		return http.StatusUnprocessableEntity
	// User errors (20xxx) → 400 Bad Request
	case e.errCode >= 20000 && e.errCode < 30000:
		return http.StatusBadRequest
	// Node errors (30xxx) → 502 Bad Gateway
	case e.errCode >= 30000 && e.errCode < 40000:
		return http.StatusBadGateway
	// Coupon errors (50xxx) → 400
	case e.errCode >= 50000 && e.errCode < 60000:
		return http.StatusBadRequest
	// Subscribe errors (60xxx) → 400
	case e.errCode >= 60000 && e.errCode < 61000:
		return http.StatusBadRequest
	// Order errors (61xxx) → 400
	case e.errCode >= 61000 && e.errCode < 70000:
		return http.StatusBadRequest
	// Auth verify errors (70xxx) → 400
	case e.errCode >= 70000 && e.errCode < 80000:
		return http.StatusBadRequest
	// DB errors (10xxx) → 500
	case e.errCode >= 10000 && e.errCode < 20000:
		return http.StatusInternalServerError
	// System errors (90xxx) → 400 (user-facing config/input errors)
	case e.errCode >= 90000 && e.errCode < 100000:
		return http.StatusBadRequest
	// Generic error
	case e.errCode == ERROR:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func NewErrCodeMsg(errCode uint32, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}
func NewErrCode(errCode uint32) *CodeError {
	return &CodeError{errCode: errCode, errMsg: MapErrMsg(errCode)}
}

func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: ERROR, errMsg: errMsg}
}
