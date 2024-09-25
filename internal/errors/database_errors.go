package errors

type NilDBConnectionError struct {
	message string
}

func (e NilDBConnectionError) Error() string {
	if e.message != "" {
		return e.message
	}

	return "DB connection error. Making operation on nil database connection."
}
