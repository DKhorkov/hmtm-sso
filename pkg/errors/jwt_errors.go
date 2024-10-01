package errors

type InvalidJWTError struct {
	Message string
}

func (e *InvalidJWTError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "JWT token is invalid or has expired"
}

type JWTClaimsError struct {
	Message string
}

func (e *JWTClaimsError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "JWT claims error"
}
