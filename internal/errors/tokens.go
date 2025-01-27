package errors

import "fmt"

type AccessTokenDoesNotBelongToRefreshTokenError struct {
	Message string
	BaseErr error
}

func (e AccessTokenDoesNotBelongToRefreshTokenError) Error() string {
	template := "access token does not belong to refresh token"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

func (e AccessTokenDoesNotBelongToRefreshTokenError) Unwrap() error {
	return e.BaseErr
}
