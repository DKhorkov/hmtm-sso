package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCError struct {
	Message string     `json:"message"`
	Status  codes.Code `json:"-"`
}

func (e GRPCError) Error() string {
	return e.Message
}

// GRPCStatus is a member function, which is used by gRPC when converting an error into a status.
func (e GRPCError) GRPCStatus() *status.Status {
	return status.New(e.Status, e.Error())
}
