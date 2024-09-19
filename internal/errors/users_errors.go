package errors

type UserNotFoundError struct {
	message string
}

func (e *UserNotFoundError) Error() string {
	if e.message != "" {
		return e.message
	}

	return "user not found"
}
