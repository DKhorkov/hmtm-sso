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

type RefreshTokenAlreadyExistsError struct {
	Message string
}

func (e RefreshTokenAlreadyExistsError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "Refresh token already exists"
}

type RefreshTokenNotFoundError struct {
	Message string
}

func (e RefreshTokenNotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "Refresh token not found"
}
