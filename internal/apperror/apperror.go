package apperror

type AppError struct {
	Code    int
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, status int, message string) *AppError {
	return &AppError{
		Code:    code,
		Status:  status,
		Message: message,
	}
}
