package kumerror

import (
	"net/http"
)

const (
	ErrorBadRequest          = "Bad Request"           //400
	ErrorUnauthorized        = "Unauthorized"          // 401
	ErrorForbidden           = "Forbidden"             // 403
	ErrorNotFound            = "Not Found"             // 404
	ErrorRequestTimeout      = "Request Timeout"       // 408
	ErrorConflict            = "Conflict"              // 409
	ErrorUnprocessableEntity = "Unprocessable Entity"  // 422
	ErrorInternalServer      = "Internal Server Error" // 500
	ErrorNotImplemented      = "Not Implemented"       // 501
	ErrorBadGateway          = "Bad Gateway"           // 502
	ErrorServiceUnavailable  = "Service Unavailable"   // 503
	ErrorUndefined           = "Undefined Error"
)

var errorHttpMap = map[string]int{
	ErrorBadRequest:          http.StatusBadRequest,
	ErrorUnauthorized:        http.StatusUnauthorized,
	ErrorForbidden:           http.StatusForbidden,
	ErrorNotFound:            http.StatusNotFound,
	ErrorRequestTimeout:      http.StatusRequestTimeout,
	ErrorConflict:            http.StatusConflict,
	ErrorUnprocessableEntity: http.StatusUnprocessableEntity,
	ErrorInternalServer:      http.StatusInternalServerError,
	ErrorNotImplemented:      http.StatusNotImplemented,
	ErrorBadGateway:          http.StatusBadGateway,
	ErrorServiceUnavailable:  http.StatusServiceUnavailable,
	ErrorUndefined:           http.StatusUnprocessableEntity,
}

func GetHttpStatus(e string) int {
	if errorHttpMap[e] != 0 {
		return errorHttpMap[e]
	}
	return http.StatusUnprocessableEntity
}
