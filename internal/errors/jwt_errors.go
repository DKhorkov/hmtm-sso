package errors

type InvalidJWTError struct {
	message string
}

func (e *InvalidJWTError) Error() string {
	if e.message != "" {
		return e.message
	}

	return "JWT token is invalid or has expired"
}

type JWTClaimsError struct {
	message string
}

func (e *JWTClaimsError) Error() string {
	if e.message != "" {
		return e.message
	}

	return "JWT claims error"
}
