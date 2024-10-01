package errors

type NilDBConnectionError struct {
	Message string
}

func (e NilDBConnectionError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "DB connection error. Making operation on nil database connection."
}
