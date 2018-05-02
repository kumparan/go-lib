package kumerror

import (
	"net/http"
)

const (
	_                        = iota
	ErrorBadRequest          // 400
	ErrorUnauthorized        // 401
	ErrorForbidden           // 403
	ErrorNotFound            // 404
	ErrorRequestTimeout      // 408
	ErrorConflict            // 409
	ErrorUnprocessableEntity // 422
	ErrorInternalServer      // 500
	ErrorNotImplemented      // 501
	ErrorBadGateway          // 502
	ErrorServiceUnavailable  // 503
	ErrorUndefined
)

var errorHttpMap = map[int64]int64{
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

func GetHttpStatus(e int64) int64 {
	if errorHttpMap[e] != 0 {
		return errorHttpMap[e]
	}
	return http.StatusUnprocessableEntity
}
