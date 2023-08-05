package errors

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsernameAlreadyExistErr struct {
	Name string
}

func (e UsernameAlreadyExistErr) Error() string {
	return fmt.Sprintf("the entity with username:%s already exist", e.Name)
}

func (e UsernameAlreadyExistErr) GRPCStatus() *status.Status {
	return status.Newf(codes.AlreadyExists, "a username %s already exists", e.Name)
}

var AutenticationFailErr = fmt.Errorf("authentication failed")

type ValidationErr struct {
	err   error
	filed string
}

func (e ValidationErr) Error() string {
	return fmt.Sprintf("validation error: %+v", e.err)
}
func (e ValidationErr) GRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.err.Error())
}

func NewValidationErr(err error) ValidationErr {
	return ValidationErr{err: err}
}
