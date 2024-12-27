package httperror

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Casagrande-Lucas/dnd/pkg/failure"
)

type ErrorMSG struct {
	ErrorCode string `json:"error_code"`
	ErrorMSG  string `json:"error_msg"`
}

type APIError struct {
	StatusCode int      `json:"status_code"`
	ErrorMSG   ErrorMSG `json:"error"`
}

func FormError(err error) APIError {
	var apiError APIError
	var svcError *failure.Error
	if errors.As(err, &svcError) {
		svcErr := svcError.SvcErr()
		switch {
		case errors.Is(svcErr, failure.ErrorInternalServer):
			apiError.StatusCode = http.StatusInternalServerError
		case errors.Is(svcErr, failure.ErrorBadRequest):
			apiError.StatusCode = http.StatusBadRequest
		case errors.Is(svcErr, failure.ErrorNotFound):
			apiError.StatusCode = http.StatusNotFound
		case errors.Is(svcErr, failure.ErrorUnauthorized):
			apiError.StatusCode = http.StatusUnauthorized
		case errors.Is(svcErr, failure.ErrorForbidden):
			apiError.StatusCode = http.StatusForbidden
		case errors.Is(svcErr, failure.ErrorDeadlineExceeded):
			apiError.StatusCode = http.StatusRequestTimeout
		case errors.Is(svcErr, failure.ErrorMethodNotAllowed):
			apiError.StatusCode = http.StatusMethodNotAllowed
		case errors.Is(svcErr, failure.ErrorNotAcceptable):
			apiError.StatusCode = http.StatusNotAcceptable
		case errors.Is(svcErr, failure.ErrorEmailAlreadyRegistered):
			apiError.StatusCode = http.StatusConflict
		default:
			apiError.StatusCode = http.StatusInternalServerError
		}

		msg := svcError.SvcErr().Error()
		apiError.ErrorMSG = ErrorMSG{ErrorCode: fmt.Sprintf("%d", apiError.StatusCode), ErrorMSG: msg}
	}
	return apiError
}
