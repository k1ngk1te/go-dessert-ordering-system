package app_errors

type NotFoundError struct {
	Message string
	Code    int
}

func (e *NotFoundError) Error() string {
	return e.Message
}
