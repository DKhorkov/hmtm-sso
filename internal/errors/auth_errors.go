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

type UserAlreadyExistsError struct {
	message string
}

func (e *UserAlreadyExistsError) Error() string {
	if e.message != "" {
		return e.message
	}

	return "user with this email already exists"
}

type InvalidPasswordError struct {
	message string
}

func (e *InvalidPasswordError) Error() string {
	if e.message != "" {
		return e.message
	}

	return "wrong password"
}
