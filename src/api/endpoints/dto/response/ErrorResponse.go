package response

type Error struct {
	Message string `json:"message"`
}

func NewError(errorMessage string) *Error {
	return &Error{Message: errorMessage}
}
