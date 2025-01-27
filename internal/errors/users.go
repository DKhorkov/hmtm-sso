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

func (e UserNotFoundError) Unwrap() error {
	return e.BaseErr
}
