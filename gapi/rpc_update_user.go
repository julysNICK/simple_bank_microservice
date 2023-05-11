package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/julysNICK/simplebank/db/sqlc"
	"github.com/julysNICK/simplebank/pb"
	"github.com/julysNICK/simplebank/utils"
	"github.com/julysNICK/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authpayload, err := server.authorizeUser(ctx)

	if err != nil {
		return nil, unauthorizedError(err)
	}

	if authpayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's profile")
	}

	violations := validateUpdateUserRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},

		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hash, err := utils.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot hash password: %v", err)

		}
		arg.HashedPassword = sql.NullString{
			String: hash,
			Valid:  true,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

	}

	user, err := server.store.UpdateUser(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "cannot find user: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot Update user: %v", err)

	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err.Error()))
	}

	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err.Error()))
		}
	}

	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err.Error()))
		}
	}

	if req.Email != nil {
		if err := val.ValidEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err.Error()))
		}
	}

	return violations

}
