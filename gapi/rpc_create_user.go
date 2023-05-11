package gapi

import (
	"context"

	db "github.com/julysNICK/simplebank/db/sqlc"
	"github.com/julysNICK/simplebank/pb"
	"github.com/julysNICK/simplebank/utils"
	"github.com/julysNICK/simplebank/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validateCreateUserRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}


	hash, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %v", err)

	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hash,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.InvalidArgument, "username or email already exist")
			}
		}

		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)

	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err.Error()))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err.Error()))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err.Error()))
	}

	if err := val.ValidEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err.Error()))
	}

	return violations

}
