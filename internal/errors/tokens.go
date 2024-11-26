package errors

type AccessTokenDoesNotBelongToRefreshTokenError struct {
	Message string
}

func (e AccessTokenDoesNotBelongToRefreshTokenError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "Access token does not belong to refresh token"
}
