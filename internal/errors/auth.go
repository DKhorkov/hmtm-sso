package errors

import "fmt"

type UserNotFoundError struct {
	Message string
	BaseErr error
}

func (e UserNotFoundError) Error() string {
	template := "user not found"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

type WrongPasswordError struct {
	Message string
	BaseErr error
}

func (e WrongPasswordError) Error() string {
	template := "wrong password"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

type UserAlreadyExistsError struct {
	Message string
	BaseErr error
}

func (e UserAlreadyExistsError) Error() string {
	template := "user with provided email already exists"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

type RefreshTokenAlreadyExistsError struct {
	Message string
	BaseErr error
}

func (e RefreshTokenAlreadyExistsError) Error() string {
	template := "refresh token already exists"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

type RefreshTokenNotFoundError struct {
	Message string
	BaseErr error
}

func (e RefreshTokenNotFoundError) Error() string {
	template := "refresh token not found"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

type InvalidEmailError struct {
	Message string
	BaseErr error
}

func (e InvalidEmailError) Error() string {
	template := "email does not meet the requirements"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

type InvalidPasswordError struct {
	Message string
	BaseErr error
}

func (e InvalidPasswordError) Error() string {
	template := "password does not meet the requirements"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}
