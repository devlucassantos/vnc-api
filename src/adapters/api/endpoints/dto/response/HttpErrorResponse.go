package response

import "net/http"

type HttpError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func NewHttpError(code int, errorMessage string) *HttpError {
	return &HttpError{
		Code:    code,
		Message: errorMessage,
	}
}

func NewBadRequestError() *HttpError {
	return &HttpError{
		Code:    http.StatusBadRequest,
		Message: "Badly formatted request, check the data sent and try again",
	}
}

func NewUnauthorizedError() *HttpError {
	return &HttpError{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized access",
	}
}

func NewForbiddenError() *HttpError {
	return &HttpError{
		Code:    http.StatusForbidden,
		Message: "Access denied",
	}
}

func NewInternalServerError() *HttpError {
	return &HttpError{
		Code:    http.StatusInternalServerError,
		Message: "An unexpected server error has occurred. Please try again",
	}
}

func NewServiceUnavailableError() *HttpError {
	return &HttpError{
		Code:    http.StatusServiceUnavailable,
		Message: "This service is temporarily unavailable. Please try again in a few minutes",
	}
}
