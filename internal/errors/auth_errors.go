package errors

type InvalidPasswordError struct {
	message string
}

func (e *InvalidPasswordError) Error() string {
	if e.message != "" {
		return e.message
	}

	return "invalid password"
}

type UserAlreadyExistsError struct {
	message string
}

func (e *UserAlreadyExistsError) Error() string {
	if e.message != "" {
		return e.message
	}

	return "user with provided email already exists"
}
