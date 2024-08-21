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
		Message: "Requisição mal formulada, verifique as informações enviadas e tente novamente.",
	}
}

func NewUnauthorizedError() *HttpError {
	return &HttpError{
		Code:    http.StatusUnauthorized,
		Message: "Acesso não autorizado.",
	}
}

func NewForbiddenError() *HttpError {
	return &HttpError{
		Code:    http.StatusForbidden,
		Message: "Acesso negado.",
	}
}

func NewInternalServerError() *HttpError {
	return &HttpError{
		Code:    http.StatusInternalServerError,
		Message: "Ocorreu um erro inesperado no servidor. Por favor, tente novamente.",
	}
}

func NewServiceUnavailableError() *HttpError {
	return &HttpError{
		Code:    http.StatusServiceUnavailable,
		Message: "Este serviço está temporariamente indisponível. Por favor, tente novamente em alguns minutos.",
	}
}
