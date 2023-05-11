package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldViolation(field string, description string) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: description,
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
			badRequest := &errdetails.BadRequest{
			FieldViolations: violations,
		}
		statusInvalid := status.New(codes.InvalidArgument, "invalid argument")

		statusInvalid, err := statusInvalid.WithDetails(badRequest)

		if err != nil {
			return  err
		}

		return  statusInvalid.Err()
}

func unauthorizedError(err error) error {
	return status.Errorf(codes.Unauthenticated,"unauthorized: %s", err)
}