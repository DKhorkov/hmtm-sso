package errors

type UserNotFoundError struct {
	Message string
}

func (e UserNotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "user not found"
}

type InvalidPasswordError struct {
	Message string
}

func (e InvalidPasswordError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "wrong password"
}

type UserAlreadyExistsError struct {
	Message string
}

func (e UserAlreadyExistsError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "user with provided email already exists"
}
