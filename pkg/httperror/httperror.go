package httperror

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Casagrande-Lucas/dnd/pkg/failure"
)

type ErrorResponse struct {
	Code string `json:"error_code"`
	MSG  string `json:"error_message"`
}

type APIError struct {
	StatusCode int           `json:"status_code"`
	ObjectErr  ErrorResponse `json:"error"`
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
		apiError.ObjectErr = ErrorResponse{Code: fmt.Sprintf("%d", apiError.StatusCode), MSG: msg}
	}
	return apiError
}
