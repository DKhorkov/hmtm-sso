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
