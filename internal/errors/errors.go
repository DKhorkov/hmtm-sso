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

type InvalidPasswordError struct {
	message string
}

func (e *InvalidPasswordError) Error() string {
	if e.message != "" {
		return e.message
	}

	return "invalid password"
}
